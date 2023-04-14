package db

import "time"

type AccountFlags struct {
	AutoDailyClaim bool `bson:"autoDailyClaim"`
}

type EpicAccountEntry struct {
	AccountID        string    `bson:"accountId"`
	RefreshToken     string    `bson:"refreshToken"`
	RefreshExpiresAt time.Time `bson:"refreshExpiresAt"`
	ClientId         string    `bson:"clientId"`
	Flags            AccountFlags
}

type UserEntry struct {
	ID              string             `bson:"_id"`
	DiscordID       string             `bson:"discordId"`
	Accounts        []EpicAccountEntry `bson:"accounts"`
	SelectedAccount int                `bson:"selectedAccount"`
	BulkFlags       AccountFlags
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}
