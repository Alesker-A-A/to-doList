package friends

import (
	"database/sql"
	"errors"

	"todo-calendar/models"
)

var (
	ErrCannotFriendSelf = errors.New("Нельзя добавить в друзья самого себя")
	ErrAlreadyFriends   = errors.New("Вы уже друзья")
	ErrRequestExists    = errors.New("Заявка уже отправлена")
	ErrRequestNotFound  = errors.New("Заявка не найдена")
)

func GetFriendship(db *sql.DB, userA, userB int) (*models.Friendship, error) {
	var f models.Friendship
	err := db.QueryRow(
		`SELECT id, requester_id, addressee_id, status, created_at, updated_at
			FROM friendships
			WHERE (requester_id = ? AND addressee_id = ?)
			OR (requester_id = ? AND addressee_id = ?)`,
		userA, userB, userB, userA,
	).Scan(&f.ID, &f.RequesterID, &f.AddresseeID, &f.Status, &f.CreatedAt, &f.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func SendRequest(db *sql.DB, requesterID, addresseeID int) (*models.Friendship, error) {
	if requesterID == addresseeID {
		return nil, ErrCannotFriendSelf
	}

	existing, err := GetFriendship(db, requesterID, addresseeID)
	if err != nil {
		return nil, err
	}

	switch {
	case existing == nil:
		if _, err := db.Exec(
			`INSERT INTO friendships (requester_id, addressee_id, status)
			 VALUES (?, ?, 'pending')`,
			requesterID, addresseeID,
		); err != nil {
			return nil, err
		}

	case existing.Status == models.FriendshipAccepted:
		return nil, ErrAlreadyFriends

	case existing.Status == models.FriendshipPending:
		if existing.RequesterID == addresseeID {
			if _, err := db.Exec(
				`UPDATE friendships
				    SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
				  WHERE id = ?`,
				existing.ID,
			); err != nil {
				return nil, err
			}
		} else {
			return nil, ErrRequestExists
		}

	case existing.Status == models.FriendshipDeclined:
		if _, err := db.Exec(
			`UPDATE friendships
			    SET requester_id = ?, addressee_id = ?, status = 'pending',
			        updated_at = CURRENT_TIMESTAMP
			  WHERE id = ?`,
			requesterID, addresseeID, existing.ID,
		); err != nil {
			return nil, err
		}
	}

	return GetFriendship(db, requesterID, addresseeID)
}

func AcceptRequest(db *sql.DB, requesterID, addresseeID int) error {
	res, err := db.Exec(
		`UPDATE friendships
		SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
		WHERE requester_id = ? AND addressee_id = ? AND status = 'pending'`,
		requesterID, addresseeID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrRequestNotFound
	}
	return nil
}

func DeclineRequest(db *sql.DB, requesterID, addresseeID int) error {
	res, err := db.Exec(
		`UPDATE friendships
		    SET status = 'declined', updated_at = CURRENT_TIMESTAMP
		  WHERE requester_id = ? AND addressee_id = ? AND status = 'pending'`,
		requesterID, addresseeID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrRequestNotFound
	}
	return nil
}

func ListFriends(db *sql.DB, userID int) ([]models.User, error) {
	rows, err := db.Query(
		`SELECT u.id, u.username, u.created_at
		   FROM friendships f
		   JOIN users u
		     ON (u.id = f.requester_id AND f.addressee_id = ?)
		     OR (u.id = f.addressee_id AND f.requester_id = ?)
		  WHERE f.status = 'accepted'`,
		userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

func ListIncomingRequests(db *sql.DB, userID int) ([]models.User, error) {
	rows, err := db.Query(
		`SELECT u.id, u.username, u.created_at
		   FROM friendships f
		   JOIN users u ON u.id = f.requester_id
		  WHERE f.addressee_id = ? AND f.status = 'pending'`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return requests, nil
}
