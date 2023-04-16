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
	Flags            UserFlag
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

type UserFlag uint32

func (f UserFlag) HasFlag(flag UserFlag) bool { return f&flag != 0 }
func (f *UserFlag) AddFlag(flag UserFlag)     { *f |= flag }
func (f *UserFlag) ClearFlag(flag UserFlag)   { *f &= ^flag }
func (f *UserFlag) ToggleFlag(flag UserFlag)  { *f ^= flag }

const (
	DEVELOPER UserFlag = 1 << iota
	BETA
	VIP
	AUTODAILY
)
