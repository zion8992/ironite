package main

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"time"
)

/** USERS **/

func (a *App) CreateUser(u User) (int64, error) {
	hash, err := a.HashPassword(u.Password)
	if err != nil {
		return 0, err
	}

	res, err := a.DB.Exec(
		"INSERT INTO users(username,email,password,date_created) VALUES (?,?,?,?)",
		u.Username, u.Email, hash, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (a *App) GetUser(id uint64) (*User, error) {
	var u User

	err := a.DB.QueryRow(
		`SELECT id,username,email,password,date_created,session_token,session_token_expires 
		 FROM users WHERE id=?`,
		id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.DateCreated,
		&u.SessionToken,
		&u.SessionTokenExpires,
	)

	return &u, err
}

func (a *App) UpdateUser(u User) error {
	_, err := a.DB.Exec(
		`UPDATE users SET username=?, email=? WHERE id=?`,
		u.Username, u.Email, u.ID,
	)
	return err
}

func (a *App) DeleteUser(id uint64) error {
	_, err := a.DB.Exec("DELETE FROM users WHERE id=?", id)
	return err
}

func (a *App) GetServersByUser(userID uint64) ([]Server, error) {
	rows, err := a.DB.Query(
		"SELECT id,name,description,xml_feed_link,playercount,owner_id FROM servers WHERE owner_id=?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	for rows.Next() {
		var s Server
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.XMLFeedLink,
			&s.PlayerCount,
			&s.OwnerID,
		); err != nil {
			return nil, err
		}
		servers = append(servers, s)
	}

	return servers, nil
}

/** PASSWORDS **/

func (a *App) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func (a *App) CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

/** SESSION TOKENS **/

func (a *App) GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (a *App) SetSessionToken(userID uint64, token string, expires time.Time) error {
	_, err := a.DB.Exec(
		"UPDATE users SET session_token=?, session_token_expires=? WHERE id=?",
		token, expires, userID,
	)
	return err
}

func (a *App) CheckSessionToken(userID uint64, token string) (bool, error) {
	var stored string
	var expires time.Time

	err := a.DB.QueryRow(
		"SELECT session_token, session_token_expires FROM users WHERE id=?",
		userID,
	).Scan(&stored, &expires)

	if err != nil {
		return false, err
	}

	if stored != token {
		return false, nil
	}

	return time.Now().Before(expires), nil
}
