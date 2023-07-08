package cmd

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/MixinNetwork/safe/governance/config"
	"github.com/MixinNetwork/safe/governance/models"
	"github.com/gofrs/uuid"
	"github.com/urfave/cli/v2"
)

func MigrateCMD(c *cli.Context) error {
	keystore := c.String("keystore")
	private := c.String("private")
	public := c.String("public")
	userID := c.String("user")
	encrypted := c.String("encrypted")

	if encrypted == "true" {
		if len(private) != 64 {
			return fmt.Errorf("Invalid private: %s", private)
		}
		if len(public) != 64 {
			return fmt.Errorf("Invalid public: %s", public)
		}
	}
	if uid, _ := uuid.FromString(userID); uid.String() != userID {
		return fmt.Errorf("Invalid user: %s", userID)
	}

	keystoreBuf, err := base64.RawURLEncoding.DecodeString(keystore)
	if err != nil {
		return err
	}
	keystoreRaw := keystoreBuf

	if encrypted == "true" {
		custodian, err := crypto.KeyFromString(private)
		if err != nil {
			return err
		}
		publicKey, err := crypto.KeyFromString(public)
		if err != nil {
			return err
		}
		key := crypto.KeyMultPubPriv(&publicKey, &custodian)
		keystoreRaw, err = models.AesDecryptCBC(key.Bytes(), keystoreBuf)
		if err != nil {
			return err
		}
	}
	var app config.App
	err = json.Unmarshal(keystoreRaw, &app)
	if err != nil {
		return err
	}
	log.Println("app", app)
	if len(app.Pin) == 6 {
		log.Printf("1. Keystore before migrate: %s", string(keystoreRaw))

		tipPub, tipPriv, _ := ed25519.GenerateKey(rand.Reader)
		log.Printf("2. Your tip private key: %s", hex.EncodeToString(tipPriv))

		err = bot.UpdateTipPin(context.Background(), app.Pin, hex.EncodeToString(tipPub), app.PinToken, app.AppID, app.SessionID, app.PrivateKey)
		if err != nil {
			return fmt.Errorf("bot.UpdateTipPin() => %v", err)
		}

		app.Pin = hex.EncodeToString(tipPriv)
		keystoreRaw, _ = json.Marshal(app)
		log.Printf("3. Keystore after migrate: %s", string(keystoreRaw))
	}

	appa, err := bot.Migrate(context.Background(), userID, app.AppID, app.SessionID, app.PrivateKey, app.Pin, app.PinToken)
	if err != nil {
		return err
	}
	log.Printf("4. App owner: %s", appa.CreatorId)
	return err
}
