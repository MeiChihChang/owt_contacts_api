package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var LEVELS = []string{"beginner", "intermediate", "advanced", "expert"}

const MAX_NAME_LENGHT = 20

type Contact struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	FullName      string    `json:"full_name"`
	Address       string    `json:"address"`
	Email         string    `json:"email"`
	Mobile        string    `json:"mobile"`
	Password      string    `json:"password"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Skills        []*Skill  `json:"skills,omitempty"`
	SkillIDs      []int     `json:"skillids,omitempty"`
	SkillIDstring string    `json:"skillidsstring,omitempty"`
}

type Skill struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *Contact) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
