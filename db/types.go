package db

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type AccountFlags struct {
	AutoDailyClaim bool `bson:"autoDailyClaim"`
}

type EpicAccountEntry struct {
	AccountID        string    `bson:"accountId"`
	RefreshToken     string    `bson:"refreshToken"`
	RefreshExpiresAt time.Time `bson:"refreshExpiresAt"`
	ClientId         string    `bson:"clientId"`
	Flags            UserFlag
}

type UserEntry struct {
	ID              string             `bson:"_id"`
	DiscordID       snowflake.ID       `bson:"discordId"`
	Accounts        []EpicAccountEntry `bson:"accounts"`
	SelectedAccount int                `bson:"selectedAccount"`
	BulkFlags       UserFlag
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}

type UserFlag uint32

func (f UserFlag) HasFlag(flag UserFlag) bool { return f&flag != 0 }
func (f *UserFlag) AddFlag(flag UserFlag)     { *f |= flag }
func (f *UserFlag) ClearFlag(flag UserFlag)   { *f &= ^flag }
func (f *UserFlag) ToggleFlag(flag UserFlag)  { *f ^= flag }

const (
	USER UserFlag = 1 << iota
	DEVELOPER
	BETA
	VIP
	AUTODAILY
)
