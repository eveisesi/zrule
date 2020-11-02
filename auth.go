package zrule

type AuthStatus struct {
	Status Status `json:"status"`
	State  string `json:"state,omitempty"`
	Token  string `json:"token,omitempty"`
}

type Status string

const (
	StatusCreated   = "created"
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusInvalid   = "invalid"
)
