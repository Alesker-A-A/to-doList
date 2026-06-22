package models

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
	FriendshipDeclined FriendshipStatus = "declined"
)

type Friendship struct {
	ID          int              `json:"id"`
	RequesterID int              `json:"requester_id"`
	AddresseeID int              `json:"addressee_id"`
	Status      FriendshipStatus `json:"status"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}
