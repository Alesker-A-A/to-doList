package friends

import (
	"database/sql"
	"errors"

	"todo-calendar/models"
)

var (
	ErrNotFriends   = errors.New("Доступ можно выдать только другу")
	ErrInvalidLevel = errors.New("Недопустимый уровень доступа")
)

func Grant(db *sql.DB, ownerID, granteeID int, level models.AccessLevel) error {
	if level != models.AccessFreeBusy && level != models.AccessView {
		return ErrInvalidLevel
	}

	f, err := GetFriendship(db, ownerID, granteeID)
	if err != nil {
		return err
	}
	if f == nil || f.Status != models.FriendshipAccepted {
		return ErrNotFriends
	}

	_, err = db.Exec(
		`INSERT INTO calendar_access (owner_id, grantee_id, level)
		 VALUES (?, ?, ?)
		 ON CONFLICT(owner_id, grantee_id) DO UPDATE SET level = excluded.level`,
		ownerID, granteeID, level,
	)
	return err
}

func Revoke(db *sql.DB, ownerID, granteeID int) error {
	_, err := db.Exec(
		`DELETE FROM calendar_access WHERE owner_id = ? AND grantee_id = ?`,
		ownerID, granteeID,
	)
	return err
}

func GetLevel(db *sql.DB, ownerID, granteeID int) (models.AccessLevel, error) {
	var level models.AccessLevel
	err := db.QueryRow(
		`SELECT level FROM calendar_access WHERE owner_id = ? AND grantee_id = ?`,
		ownerID, granteeID,
	).Scan(&level)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return level, nil
}

func ListGrantedToMe(db *sql.DB, granteeID int) ([]models.SharedCalendar, error) {
	rows, err := db.Query(
		`SELECT u.id, u.username, u.created_at, ca.level
		   FROM calendar_access ca
		   JOIN users u ON u.id = ca.owner_id
		  WHERE ca.grantee_id = ?`,
		granteeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.SharedCalendar{}
	for rows.Next() {
		var sc models.SharedCalendar
		if err := rows.Scan(&sc.User.ID, &sc.User.Username, &sc.User.CreatedAt, &sc.Level); err != nil {
			return nil, err
		}
		result = append(result, sc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func ListMyGrants(db *sql.DB, ownerID int) ([]models.SharedCalendar, error) {
	rows, err := db.Query(
		`SELECT u.id, u.username, u.created_at, ca.level
		   FROM calendar_access ca
		   JOIN users u ON u.id = ca.grantee_id
		  WHERE ca.owner_id = ?`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.SharedCalendar{}
	for rows.Next() {
		var sc models.SharedCalendar
		if err := rows.Scan(&sc.User.ID, &sc.User.Username, &sc.User.CreatedAt, &sc.Level); err != nil {
			return nil, err
		}
		result = append(result, sc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
