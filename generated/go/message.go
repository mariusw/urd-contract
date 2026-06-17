// Code generated from message.schema.json. Do not edit.

package contract

import "time"

type Message struct {
	ID              string    `json:"id"`
	SessionID       string    `json:"session_id"`
	ParentSessionID string    `json:"parent_session_id,omitempty"`
	Harness         string    `json:"harness"`
	CWD             string    `json:"cwd"`
	GitBranch       string    `json:"git_branch"`
	Project         string    `json:"project,omitempty"`
	Role            string    `json:"role"`
	Content         string    `json:"content"`
	Timestamp       time.Time `json:"timestamp"`
	Sequence        int       `json:"sequence"`
}
