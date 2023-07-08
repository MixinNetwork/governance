package config

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"log"
	"testing"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	InitConfiguration("test")
	assert.Equal("/tmp/governance_test.sqlite3", AppConfig.Database.Path)
	assert.Equal("a9f616af-635a-4047-a32c-526340d3c241", AppConfig.Mixin.ClientID)

	assert.Equal("100", AppConfig.Governance.Fee)

	apps, err := FetchApps()
	assert.Nil(err)
	assert.Len(apps, 50)
	assert.Equal("f857e241-9f04-4c55-b3ca-48bfda6675df", apps[0].AppID)

	key, _ := mixin.KeyFromString("6196d87a2ee934da04f51e21c4542674c7ac9a57bf1eb6a39c7d65ac7318680b")

	tipBody := bot.TipBodyForOwnershipTransfer("f857e241-9f04-4c55-b3ca-48bfda6675df")
	sig := key.Sign(tipBody)
	log.Println(hex.EncodeToString(sig[:]), key.Public())

	keyBuf, _ := hex.DecodeString("6196d87a2ee934da04f51e21c4542674c7ac9a57bf1eb6a39c7d65ac7318680bf0475fbc7284330aba8f42ade1b321097e64a1c7a81a8a6792013783450b5610")
	privateKey := make([]byte, ed25519.PrivateKeySize)
	copy(privateKey, keyBuf)
	log.Println(hex.EncodeToString(privateKey))

	sigBuf := ed25519.Sign(ed25519.PrivateKey(privateKey), tipBody)
	log.Println(hex.EncodeToString(sigBuf[:]))
	pub, private, _ := ed25519.GenerateKey(rand.Reader)
	log.Println(hex.EncodeToString(pub), hex.EncodeToString(private))
	seedBuf, _ := hex.DecodeString("628ac11ce4f075f4933fcfbd269bf1bebabb9a2c4b0d0ba6eae573d9a7d980e8")
	private = ed25519.NewKeyFromSeed(seedBuf)
	log.Println(hex.EncodeToString(pub), hex.EncodeToString(private))
}
