package database

import (
	"github.com/surrealdb/surrealdb.go/pkg/models"
	"time"
)

type EpicAccount struct {
	ID                    *models.RecordID `json:"id,omitempty"`
	AccountId             string           `json:"account_id"`
	CreatedClientId       string           `json:"created_client_id"`
	Flags                 int              `json:"flags"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
}

type User struct {
	ID        *models.RecordID `json:"id,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	DiscordId string           `json:"discord_id"`
	//EpicAccounts          []models.RecordID `json:"epic_accounts"`
	GlobalFlags           int       `json:"global_flags"`
	SelectedEpicAccountId string    `json:"selected_epic_account_id"`
	UpdatedAt             time.Time `json:"updated_at"`
}
