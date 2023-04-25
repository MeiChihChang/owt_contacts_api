package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/MeiChihChang/owt_contacts_api/internal/models"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) AllContacts() ([]*models.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var contacts []*models.Contact

	query := `select id, first_name, last_name, full_name, email, address, mobile, created_at, updated_at from contacts`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		if rows == nil {
			return contacts, nil
		}
		return contacts, err
	}

	defer rows.Close()

	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID,
			&contact.FirstName,
			&contact.LastName,
			&contact.FullName,
			&contact.Email,
			&contact.Address,
			&contact.Mobile,
			&contact.UpdatedAt,
			&contact.CreatedAt,
		)
		if err != nil {
			return contacts, err
		}

		contacts = append(contacts, &contact)
	}

	if err = rows.Err(); err != nil {
		return contacts, err
	}

	return contacts, nil

}

func (m *PostgresDBRepo) AllSkills() ([]*models.Skill, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, level, created_at, updated_at from skills order by name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*models.Skill

	for rows.Next() {
		var s models.Skill
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Level,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		skills = append(skills, &s)
	}

	return skills, nil
}

func (m *PostgresDBRepo) GetContactByEmail(email string) (*models.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password,
			created_at, updated_at from contacts where email = $1`

	var user models.Contact
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) InsertContact(user models.Contact) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into contacts (first_name, last_name, full_name, address, email, mobile, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.FullName,
		user.Address,
		user.Email,
		user.Mobile,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateContact(user models.Contact) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update contacts set first_name = $1, last_name = $2, full_name = $3,
	address = $4, email = $5, mobile = $6, updated_at = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.FullName,
		user.Address,
		user.Email,
		user.Mobile,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteContact(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from contacts where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) UpdateContactSkills(id int, skillIDs []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from contacts_skills where contact_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	for _, n := range skillIDs {
		stmt := `insert into contacts_skills (contact_id, skill_id) values ($1, $2)`
		_, err := m.DB.ExecContext(ctx, stmt, id, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostgresDBRepo) GetContactByID(id int) (*models.Contact, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, first_name, last_name, full_name, address, email, mobile,
			created_at, updated_at from contacts where id = $1`

	var contact models.Contact
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&contact.ID,
		&contact.FirstName,
		&contact.LastName,
		&contact.FullName,
		&contact.Address,
		&contact.Email,
		&contact.Mobile,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// get skills, if any
	query = `select s.id, s.name, s.level from contacts_skills mg
		left join skills s on (mg.skill_id = s.id)
		where mg.contact_id = $1
		order by s.name`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var skills []*models.Skill

	for rows.Next() {
		var s models.Skill
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Level,
		)
		if err != nil {
			return nil, err
		}

		skills = append(skills, &s)
	}

	contact.Skills = skills

	return &contact, err
}

func (m *PostgresDBRepo) GetSkillByID(id int) (*models.Skill, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, level, created_at, updated_at from skills where id = $1`

	var skill models.Skill
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&skill.ID,
		&skill.Name,
		&skill.Level,
		&skill.CreatedAt,
		&skill.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &skill, err
}

func (m *PostgresDBRepo) InsertSkill(skill models.Skill) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into skills (name, level, created_at, updated_at) values ($1, $2, $3, $4) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		skill.Name,
		skill.Level,
		skill.CreatedAt,
		skill.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateSkill(skill models.Skill) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update skills set name = $1, level = $2, updated_at = $3 where id = $4`

	_, err := m.DB.ExecContext(ctx, stmt,
		skill.Name,
		skill.Level,
		skill.UpdatedAt,
		skill.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) DeleteSkill(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from skills where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
