package models

import (
	"context"
	"os"

	"github.com/MixinNetwork/safe/governance/config"
	"github.com/MixinNetwork/safe/governance/session"
	"github.com/MixinNetwork/safe/governance/store"
)

func teardownTestContext(ctx context.Context) {
	db := config.AppConfig.Database
	err := os.Remove(db.Path)
	if err != nil {
		panic(err)
	}
}

func setupTestContext() context.Context {
	config.InitConfiguration("test")

	db, err := store.OpenDatabase()
	if err != nil {
		panic(err)
	}

	return session.WithDatabase(context.Background(), db)
}
