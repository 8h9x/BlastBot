package db

import "time"

type EpicAccountEntry struct {
	ID               string    `bson:"accountId"`
	RefreshToken     string    `bson:"refreshToken"`
	RefreshExpiresAt time.Time `bson:"refreshExpiresAt"`
	ClientId         string    `bson:"clientId"`
}

type UserEntry struct {
	ID        string             `bson:"_id"`
	DiscordID string             `bson:"discordId"`
	Accounts  []EpicAccountEntry `bson:"accounts"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
