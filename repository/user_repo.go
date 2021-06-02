package repository

import (
	"context"
	"database/sql"

	"github.com/prastamaha/auth-basic/models"
	"github.com/prastamaha/auth-basic/utils"
)

// User Repository interface
type UserRepository interface {
	Insert(ctx context.Context, user *models.User) (string, error)
	DeleteByID(ctx context.Context, id uint64) (string, error)
	QueryByEmail(ctx context.Context, email string) (models.User, error)
}

// User Repository implementation struct
type UserRepositoryImpl struct {
	DB *sql.DB
}

// Create New User Repository Implementation
func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

// Insert User Implementation
func (repo *UserRepositoryImpl) Insert(ctx context.Context, user *models.User) (string, error) {
	query := "INSERT INTO users(name, email, password) values (?,?,?)"
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		return "failed to prepare context", err
	}
	defer stmt.Close()

	password, err := utils.HashPassword(user.Password)
	if err != nil {
		return err.Error(), err
	}

	result, err := stmt.ExecContext(ctx, &user.Name, &user.Email, &password)
	if err != nil {
		return "failed to register user", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "failed to get user id", err
	}

	user.Id = uint(id)

	return "user create success", nil
}

// Delete User by id Implementation

func (repo *UserRepositoryImpl) DeleteByID(ctx context.Context, id uint64) (string, error) {
	query := "DELETE FROM users WHERE id = ?"
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		return "failed to prepare context", err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, &id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "user not found", err
		}
		return "failed to delete user", err
	}

	return "user delete success", nil
}

func (repo *UserRepositoryImpl) QueryByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := "SELECT id,name,email,password FROM users WHERE email = ?"
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, &email)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		return user, nil
	} else {
		return user, err
	}

}
