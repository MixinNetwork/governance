package models

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/go-number"
	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/MixinNetwork/safe/governance/config"
	"github.com/MixinNetwork/safe/governance/externals"
	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/store"
	"github.com/gofrs/uuid"
)

type Node struct {
	Custodian string
	Payee     string
	KernelID  string
	AppID     sql.NullString
	MixinHash sql.NullString
	Keystore  string
	PublicKey string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var nodesColumns = []string{"custodian", "payee", "kernel_id", "app_id", "mixin_hash", "keystore", "public_key", "created_at", "updated_at"}

func (n *Node) values() []any {
	return []any{n.Custodian, n.Payee, n.KernelID, n.AppID, n.MixinHash, n.Keystore, n.PublicKey, n.CreatedAt, n.UpdatedAt}
}

func nodeFromRow(row store.Row) (*Node, error) {
	var n Node
	err := row.Scan(&n.Custodian, &n.Payee, &n.KernelID, &n.AppID, &n.MixinHash, &n.Keystore, &n.PublicKey, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &n, err
}

func CreateNode(ctx context.Context, custodian, payee, kernelID, appID, hash string) (*Node, error) {
	t := time.Now()
	node := &Node{
		Custodian: custodian,
		Payee:     payee,
		KernelID:  kernelID,
		AppID:     sql.NullString{String: appID, Valid: true},
		MixinHash: sql.NullString{String: hash, Valid: true},
		Keystore:  "",
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := store.BuildInsertionSQL("nodes", nodesColumns)
		tx.ExecContext(ctx, query, node.values()...)
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return node, nil
}

// extra: custodian (common.Address) || payee (common.Address) || node id (crypto.Hash)
func CreateNodeByExtra(ctx context.Context, extra string) (*Node, error) {
	custodian, payee, kernel, err := validateExtra(ctx, extra)
	if err != nil {
		return nil, err
	}
	t := time.Now()
	node := &Node{
		Custodian: custodian.String(),
		Payee:     payee.String(),
		KernelID:  kernel.String(),
		CreatedAt: t,
		UpdatedAt: t,
	}

	traceID := bot.UniqueObjectId(node.Custodian, node.Payee, node.KernelID)
	data := bot.EncodeMixinExtra(uuid.Nil.String(), extra)
	amount := number.FromString(fmt.Sprint(len(data)/1024 + 2)).Mul(number.FromString("0.001"))
	in := &bot.ObjectInput{
		Amount:  amount,
		TraceId: traceID,
		Memo:    extra,
	}
	mixin := config.AppConfig.Mixin
	snapshot, err := bot.CreateObject(context.Background(), in, mixin.ClientID, mixin.SessionID, mixin.PrivateKey, mixin.Pin, mixin.PinToken)
	if err != nil {
		return nil, err
	} else if snapshot == nil || snapshot.TransactionHash == "" {
		return nil, session.ServerError(ctx, fmt.Errorf("invalid snapshot %s", traceID))
	}
	node.MixinHash = sql.NullString{String: snapshot.TransactionHash, Valid: true}

	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		old, err := findNode(ctx, tx, node.Custodian, node.Payee, node.KernelID, node.AppID.String, "")
		if err != nil {
			return err
		} else if old != nil {
			node = old
			return nil
		}
		query := store.BuildInsertionSQL("nodes", nodesColumns)
		_, err = tx.ExecContext(ctx, query, node.values()...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return node, nil
}

func PaymentNode(ctx context.Context, hash string) (*Node, error) {
	set, err := ReadNodeSet(ctx)
	if err != nil {
		return nil, err
	}
	apps, err := config.FetchApps()
	if err != nil {
		return nil, err
	}
	transaction, err := externals.ReadTransaction(hash)
	if err != nil {
		return nil, err
	} else if transaction == nil {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "hash", "invalid", hash)
	}
	extraBuf, err := hex.DecodeString(transaction.Extra)
	if err != nil {
		return nil, err
	}
	pack := bot.DecodeMixinExtra(extraBuf)
	_, _, _, err = validateExtra(ctx, pack.M)
	if err != nil {
		return nil, err
	}

	var node *Node
	err = session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		old, err := findNode(ctx, tx, "", "", "", "", hash)
		if err != nil || old == nil {
			return err
		}
		node = old
		if node.AppID.String != "" {
			return nil
		}
		var app *config.App
		for _, a := range apps {
			if set[a.AppID] != nil {
				continue
			}
			app = a
			node.AppID = sql.NullString{String: app.AppID, Valid: true}
		}
		if node.AppID.String == "" {
			return session.BadDataErrorWithFieldAndData(ctx, "app id", "invalid", "")
		}

		appBuf, err := json.Marshal(app)
		if err != nil {
			return err
		}

		custodian, err := common.NewAddressFromString(node.Custodian)
		if err != nil {
			return err
		}
		mixin := config.AppConfig.Mixin
		privateBuf, _ := base64.RawURLEncoding.DecodeString(mixin.PrivateKey)
		privateBot := crypto.NewKeyFromSeed(privateBuf)
		key := crypto.KeyMultPubPriv(&custodian.PublicSpendKey, &privateBot)
		encryptedBuf := AesEncryptCBC(key.Bytes(), appBuf)
		node.Keystore = base64.RawURLEncoding.EncodeToString(encryptedBuf)
		node.PublicKey = privateBot.Public().String()

		_, err = tx.ExecContext(ctx, "UPDATE nodes SET app_id=?,keystore=?,public_key=? WHERE custodian=?", node.AppID, base64.RawURLEncoding.EncodeToString(encryptedBuf), node.PublicKey, node.Custodian)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return node, nil
}

func ReadNodes(ctx context.Context) ([]*Node, error) {
	var nodes []*Node
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM nodes WHERE app_id IS NOT NULL LIMIT 100", strings.Join(nodesColumns, ","))
		rows, err := tx.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		for rows.Next() {
			node, err := nodeFromRow(rows)
			if err != nil {
				return err
			}
			nodes = append(nodes, node)
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return nodes, nil
}

func ReadNodeSet(ctx context.Context) (map[string]*Node, error) {
	nodes, err := ReadNodes(ctx)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*Node, 0)
	for _, n := range nodes {
		set[n.AppID.String] = n
	}
	return set, nil
}

func ReadNode(ctx context.Context, custodian string) (*Node, error) {
	var node *Node
	err := session.Database(ctx).RunInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		old, err := findNode(ctx, tx, custodian, "", "", "", "")
		node = old
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return node, nil
}

func findNode(ctx context.Context, tx *sql.Tx, custodian, payee, kernel, app, hash string) (*Node, error) {
	query := fmt.Sprintf("SELECT %s FROM nodes WHERE custodian=? OR payee=? OR kernel_id=? OR app_id=? OR mixin_hash=?", strings.Join(nodesColumns, ","))
	row := tx.QueryRowContext(ctx, query, custodian, payee, kernel, app, hash)
	node, err := nodeFromRow(row)
	return node, err
}

func validateExtra(ctx context.Context, extra string) (*common.Address, *common.Address, *crypto.Hash, error) {
	raw, err := base64.RawURLEncoding.DecodeString(extra)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(raw) != 353 {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "extra", "invalid", extra)
	}
	if raw[0] != byte(1) {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "extra head", "invalid", extra)
	}
	custodianBytes := raw[1:65]
	payeeBytes := raw[65:129]
	nodeBytes := raw[129:161]
	sigSignerBytes := raw[161:225]
	sigPayeeBytes := raw[225:289]
	sigCustodianBytes := raw[289:]

	custodian := common.Address{}
	copy(custodian.PublicSpendKey[:], custodianBytes[:32])
	copy(custodian.PublicViewKey[:], custodianBytes[32:])

	payee := common.Address{}
	copy(payee.PublicSpendKey[:], payeeBytes[:32])
	copy(payee.PublicViewKey[:], payeeBytes[32:])

	var kernel crypto.Hash
	copy(kernel[:], nodeBytes)

	var sigSigner crypto.Signature
	copy(sigSigner[:], sigSignerBytes)

	var sigPayee crypto.Signature
	copy(sigPayee[:], sigPayeeBytes)

	var signCustodian crypto.Signature
	copy(signCustodian[:], sigCustodianBytes)

	nodes, err := externals.ListAllNodes()
	if err != nil {
		return nil, nil, nil, err
	}
	var signerStr string
	for _, n := range nodes {
		if n.Id == kernel.String() {
			if n.State == "ACCEPTED" {
				if n.Payee == payee.String() {
					signerStr = n.Signer
				}
				break
			}
			if n.State == "REMOVED" {
				t := time.Unix(0, n.Timestamp)
				if t.Add(7 * 24 * time.Hour).After(time.Now()) {
					if n.Payee == payee.String() {
						signerStr = n.Signer
					}
					break
				}
			}
		}
	}
	if signerStr == "" {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "signer", "not existing", extra)
	}

	signer, err := common.NewAddressFromString(signerStr)
	if err != nil {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "signer", "invalid", extra)
	}

	if !signer.PublicSpendKey.Verify(raw[:161], sigSigner) {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "signer signature verify", "invalid", extra)
	}

	if !payee.PublicSpendKey.Verify(raw[:161], sigPayee) {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "payee signature verify", "invalid", extra)
	}

	if !custodian.PublicSpendKey.Verify(raw[:161], signCustodian) {
		return nil, nil, nil, session.BadDataErrorWithFieldAndData(ctx, "custodian signature verify", "invalid", extra)
	}
	return &custodian, &payee, &kernel, nil
}

func AesEncryptCBC(key, msg []byte) []byte {
	padding := aes.BlockSize - len(msg)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	msg = append(msg, padtext...)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(hex.EncodeToString(key))
	}
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	n, err := io.ReadFull(rand.Reader, iv)
	if n != aes.BlockSize || err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], msg)
	return ciphertext
}

func AesDecryptCBC(key, ciphertext []byte) ([]byte, error) {
	if cl := len(ciphertext); cl < aes.BlockSize || cl%aes.BlockSize != 0 {
		return nil, fmt.Errorf("AES cipher text invalid length %d", cl)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]
	source := ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(source, source)

	length := len(source)
	unpadding := int(source[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("AES CBC padding invalid %d %d", unpadding, length)
	}
	return source[:length-unpadding], nil
}
