package repository

import (
	"github.com/MeiChihChang/owt_contacts_api/internal/models"
)

type DatabaseRepo interface {
	AllContacts() ([]*models.Contact, error)
	GetContactByEmail(email string) (*models.Contact, error)
	GetContactByID(id int) (*models.Contact, error)
	InsertContact(movie models.Contact) (int, error)
	UpdateContact(movie models.Contact) error
	DeleteContact(id int) error

	AllSkills() ([]*models.Skill, error)
	GetSkillByID(id int) (*models.Skill, error)
	InsertSkill(skill models.Skill) (int, error)
	UpdateSkill(skill models.Skill) error
	DeleteSkill(id int) error

	UpdateContactSkills(id int, skillIDs []int) error
}
