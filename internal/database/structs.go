package database

import (
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type AccountFlags struct {
	AutoDailyClaim bool `bson:"autoDailyClaim"` // TODO: remove, daily rewards are gone
}

type EpicAccount struct {
	AccountID        string    `bson:"accountId"`
	CreatedClientID  string    `bson:"createdClientId"`
	RefreshToken     string    `bson:"refreshToken"`
	RefreshExpiresAt time.Time `bson:"refreshExpiresAt"`
	Flags            UserFlag  `bson:"flags"`
}

type User struct {
	ID                    string        `bson:"_id"`
	DiscordID             snowflake.ID  `bson:"discordId"`
	Accounts              []EpicAccount `bson:"accounts"`
	SelectedEpicAccountId string        `bson:"selectedAccount"`
	BulkFlags             UserFlag      `bson:"bulkFlags"`
	CreatedAt             time.Time     `bson:"createdAt"`
	UpdatedAt             time.Time     `bson:"updatedAt"`
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
