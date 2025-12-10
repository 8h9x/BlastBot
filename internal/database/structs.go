package database

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

//type EpicAccount struct {
//	ID                    *models.RecordID `json:"id,omitempty"`
//	AccountId             string           `json:"account_id"`
//	CreatedClientId       string           `json:"created_client_id"`
//	Flags                 int              `json:"flags"`
//	RefreshToken          string           `json:"refresh_token"`
//	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
//}
//
//type User struct {
//	ID        *models.RecordID `json:"id,omitempty"`
//	CreatedAt time.Time        `json:"created_at"`
//	DiscordId string           `json:"discord_id"`
//	//EpicAccounts          []models.RecordID `json:"epic_accounts"`
//	GlobalFlags           int       `json:"global_flags"`
//	SelectedEpicAccountId string    `json:"selected_epic_account_id"`
//	UpdatedAt             time.Time `json:"updated_at"`
//}
