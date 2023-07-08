package blaze

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/safe/governance/config"
	"github.com/MixinNetwork/safe/governance/models"
)

type mixinBlazeHandler func(ctx context.Context, msg bot.MessageView, clientID string) error

func (f mixinBlazeHandler) OnMessage(ctx context.Context, msg bot.MessageView, clientID string) error {
	return f(ctx, msg, clientID)
}

func (f mixinBlazeHandler) OnAckReceipt(ctx context.Context, msg bot.MessageView, clientID string) error {
	return nil
}

func (f mixinBlazeHandler) SyncAck() bool {
	return true
}

func Boot(ctx context.Context) error {
	log.Println("Mixin Safe Governance start blaze service")
	mixin := config.AppConfig.Mixin
	for {
		client := bot.NewBlazeClient(mixin.ClientID, mixin.SessionID, mixin.PrivateKey)
		h := func(ctx context.Context, botMsg bot.MessageView, clientID string) error {
			err := handleMessage(ctx, client, botMsg)
			if err != nil {
				log.Printf("blaze.handleMessage() => %v", err)
				return err
			}
			return nil
		}
		if err := client.Loop(ctx, mixinBlazeHandler(h)); err != nil {
			log.Printf("client.Loop() => %#v", err)
		}
		time.Sleep(time.Second)
	}
}

func handleMessage(ctx context.Context, bc *bot.BlazeClient, bm bot.MessageView) error {
	dataRaw, err := base64.StdEncoding.DecodeString(bm.Data)
	if err != nil {
		return err
	}
	mixin := config.AppConfig.Mixin
	if bm.Category == "SYSTEM_ACCOUNT_SNAPSHOT" && bm.UserId != mixin.ClientID {
		var transfer bot.TransferView
		err = json.Unmarshal(dataRaw, &transfer)
		if err != nil {
			return err
		}
		governance := config.AppConfig.Governance
		if transfer.AssetId != governance.FeeAssetID {
			return nil
		}
		if transfer.Amount != governance.Fee {
			return nil
		}
		if transfer.Memo == "" {
			return nil // TODO deposit to the bot will error
		}
		_, err = models.PaymentNode(ctx, transfer.Memo)
		return err
	}
	return nil
}
