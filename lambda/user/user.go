package user

import (
	"time"
)

// User is a struct representing a user used for marshalling via sqlx
type User struct {
	ID                int       `db:"id"`
	Username          string    `db:"username"`
	Password          string    `db:"password"`
	CreationDate      time.Time `db:"creation_date"`
	APIKey            string    `db:"api_key"`
	EncryptionEnabled bool      `db:"encryption_enabled"` // Unused
	ThemeName         string    `db:"theme_name"`         // Unused
}
