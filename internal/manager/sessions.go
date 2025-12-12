package manager

import (
	"github.com/8h9x/vinderman"
	"github.com/disgoorg/snowflake/v2"
)

// TODO: At some point it would probably make sense to move some of the auth invalidation logic and coupled event handlers upstream to vinderman

type SessionManager struct {
	sessions map[snowflake.ID]vinderman.Client
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// func (sm *SessionManager)

// func (am *SessionManager) StartRefreshLoop() {

// }
