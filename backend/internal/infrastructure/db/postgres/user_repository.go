package postgres

import (
	"context"
	"database/sql"

	"github.com/project-capillary/backend/internal/domain/entities"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.Role,
		user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
		&user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE username = $1
	`
	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
		&user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
		&user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = $2, password_hash = $3, first_name = $4, last_name = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.IsActive, user.UpdatedAt)
	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, first_name, last_name, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role,
			&user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
