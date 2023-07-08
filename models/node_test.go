package models

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/MixinNetwork/safe/governance/config"
	"github.com/stretchr/testify/assert"
)

func TestNodeCRUD(t *testing.T) {
	assert := assert.New(t)

	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	node, err := CreateNode(ctx, "custodian", "payee", "kernel", "app", "hash")
	assert.Nil(err)
	assert.NotNil(node)
	node, err = ReadNode(ctx, node.Custodian)
	assert.Nil(err)
	assert.NotNil(node)
	assert.Equal("custodian", node.Custodian)
	nodes, err := ReadNodes(ctx)
	assert.Nil(err)
	assert.Len(nodes, 1)

	addressStr := "XINLijHgKMB25puXb16JGrmD8DY3rmZno4iR2Bemtc8Ct465H5AJ32Sr4YCgTCYZVz79ZsY4c9ut7ZeUc6wQCZboNUqWegCe"

	address, err := common.NewAddressFromString(addressStr)
	assert.Nil(err)
	data, err := address.MarshalJSON()
	assert.Nil(err)
	// Address byte len 98
	log.Println(len(data))
	addressNew := common.Address{}
	err = addressNew.UnmarshalJSON(data)
	assert.Nil(err)
	assert.Equal(addressStr, addressNew.String())

	spentStr := "a229abeb186aad03a6968a9c3c6c5e0f90fe3fb498eec1f9acacf70aad6dce0f"
	spent, _ := crypto.KeyFromString(spentStr)
	sig := spent.Sign([]byte(addressStr))
	view := spent.Public()
	assert.True(view.Verify([]byte(addressStr), sig))

	// custodian
	// address:	XINJYiri2BU4dLGdsj33C5pvDuhzxK7DmWB9PvABa7u53tCoabApajFRsNTbsLjm2tjPfRQJEN2Awpe8SP3V35CMGRm2A5N1
	// view key:	b38f4c858cdc2f861c9c04b4714f24c8e950ce516541373390055d24c567a40a
	// spend key:	bdfe0792f1d613d7842587e6bce8a05e549876b5a840c47a0577b0540864ba0e

	// signer
	// address:	XIN4qtYcAuAsJFnHp61waUheVsiK1byouLqbhrA8VpSQwxHs4z8LPjpFRrx3zdmiXZuFSwJ8CAMCwLkxap1LbRWHk2iVsLyx
	// view key:	af458bbe67afd1230e2b5b128840be40c10398472d45d88c278510fd9ecdf10a
	// spend key:	ed4c90d8a0a34e4a3e564ea1ee5399a14a920a8cc2fdc56be3e5fba88c44350e

	// payee
	// address:	XINYvDWLAqoa1PxNxAaJcecrrehHVaaqqT4owg7ST1Yt2Gs5VUX62ArnVW7rx3vBMxfRdA5Y6kEg1Y5jSdQDFF3msunpmED4
	// view key:	926abf607d33577dc25947764c95e4fe5699a8ec45ec072fd9d89d2f0b700a06
	// spend key:	f38222cdd1c17bbf748afa4b74c829785b4af24ea3b5b2172db04f413adc260c

	// node
	// id: 394e7b2131b7d0a996bb094e30d05ac7d51f5a09156e5f7349cac55d2179a144

	extra := []byte{1}
	custodian, err := common.NewAddressFromString("XINJYiri2BU4dLGdsj33C5pvDuhzxK7DmWB9PvABa7u53tCoabApajFRsNTbsLjm2tjPfRQJEN2Awpe8SP3V35CMGRm2A5N1")
	assert.Nil(err)
	extra = append(extra, custodian.PublicViewKey[:]...)
	extra = append(extra, custodian.PublicSpendKey[:]...)

	// signer, err := common.NewAddressFromString("XIN4qtYcAuAsJFnHp61waUheVsiK1byouLqbhrA8VpSQwxHs4z8LPjpFRrx3zdmiXZuFSwJ8CAMCwLkxap1LbRWHk2iVsLyx")
	// assert.Nil(err)

	payee, err := common.NewAddressFromString("XINYvDWLAqoa1PxNxAaJcecrrehHVaaqqT4owg7ST1Yt2Gs5VUX62ArnVW7rx3vBMxfRdA5Y6kEg1Y5jSdQDFF3msunpmED4")
	assert.Nil(err)
	extra = append(extra, payee.PublicViewKey[:]...)
	extra = append(extra, payee.PublicSpendKey[:]...)

	kernel, err := crypto.HashFromString("394e7b2131b7d0a996bb094e30d05ac7d51f5a09156e5f7349cac55d2179a144")
	assert.Nil(err)
	extra = append(extra, kernel[:]...)

	keySigner, _ := crypto.KeyFromString("ed4c90d8a0a34e4a3e564ea1ee5399a14a920a8cc2fdc56be3e5fba88c44350e")
	sigSigner := keySigner.Sign(extra)
	extra = append(extra, sigSigner[:]...)

	keyPayee, _ := crypto.KeyFromString("f38222cdd1c17bbf748afa4b74c829785b4af24ea3b5b2172db04f413adc260c")
	sigPayee := keyPayee.Sign(extra[:161])
	extra = append(extra, sigPayee[:]...)

	keyCustodian, _ := crypto.KeyFromString("bdfe0792f1d613d7842587e6bce8a05e549876b5a840c47a0577b0540864ba0e")
	sigCustodian := keyCustodian.Sign(extra[:161])
	extra = append(extra, sigCustodian[:]...)

	node, err = CreateNodeByExtra(ctx, base64.RawURLEncoding.EncodeToString(extra))
	assert.Nil(err)
	assert.NotNil(node)
	node, err = ReadNode(ctx, node.Custodian)
	assert.Nil(err)
	assert.NotNil(node)
	assert.Equal("", node.AppID.String)
	return

	node, err = PaymentNode(ctx, "5e7f37fd76bea1647d46c396e21c6496f3033f03ea50121500c6e6c2df5294b7")
	assert.Nil(err)
	assert.NotNil(node)
	node, err = ReadNode(ctx, node.Custodian)
	assert.Nil(err)
	assert.NotNil(node)
	assert.NotEqual("", node.AppID.String)
	assert.NotEqual("", node.Keystore)

	apps, _ := config.FetchApps()
	assert.Equal("a9f616af-635a-4047-a32c-526340d3c241", apps[0].AppID)
	mixin := config.AppConfig.Mixin
	privateBuf, _ := base64.RawURLEncoding.DecodeString(mixin.PrivateKey)

	privateCustodian, _ := crypto.KeyFromString("bdfe0792f1d613d7842587e6bce8a05e549876b5a840c47a0577b0540864ba0e")
	privateBotKey := crypto.NewKeyFromSeed(privateBuf)

	publicBotKey := privateBotKey.Public()
	key1 := hex.EncodeToString(crypto.KeyMultPubPriv(&publicBotKey, &privateCustodian).Bytes())
	key2 := hex.EncodeToString(crypto.KeyMultPubPriv(&custodian.PublicSpendKey, &privateBotKey).Bytes())
	assert.Equal(key1, key2)

	keystore, _ := base64.RawURLEncoding.DecodeString("pzJxZtFBaAHQ6JBQgj_b8-u-hpizaBK9HZlkZ2xuLXpYebbcZSAnGuvoG0tu-jffAxczNWMI6FUMgnR2LCb1Wp5ksiFeq7VuwbyZ51djxmXjW4nk3qEGvhDF2rNqXNQ0XgLrtFMWzk4VKOrUsVW6ySTsmo1i4o8lLFczzWYja0_ayLqrP1Hoc_YiwKMYz6fpRTf0k8oRbng5-5kqWiZaRLz_7GjqQyQLNVCWRYH0BTUOVnDJ1jS1twIxTshewEEhFIth4lHXcNS8KOaRYA9zA_yNQxItG-QkjmJDMMWq8wk")
	publicBotKey, _ = crypto.KeyFromString("b80200fafa1856c0dcb9fdb422e3b941cc86c9c1c4e352058c845d45db8b99e3")

	keystoreBuf, _ := AesDecryptCBC(crypto.KeyMultPubPriv(&publicBotKey, &privateCustodian).Bytes(), keystore)
	log.Println(string(keystoreBuf))
}
