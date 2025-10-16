package service

import (
	"context"
	// "errors"
	// "time"
	// "golang.org/x/crypto/bcrypt"
	//
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"koois_core/internal/model"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	rows, err := s.db.Query(ctx, "SELECT user_id, username, email, password, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserId, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
// 	var user model.User
// 	err := s.db.QueryRow(ctx, "SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1", id).
// 		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
//
// 	if err == pgx.ErrNoRows {
// 		return nil, errors.New("user not found")
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &user, nil
// }
//
// func (s *UserService) Create(ctx context.Context, req model.CreateUserReq) (*model.User, error) {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	now := time.Now()
// 	var id int
//
// 	err = s.db.QueryRow(ctx,
// 		"INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
// 		req.Name, req.Email, string(hashedPassword), now, now).
// 		Scan(&id)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &model.User{
// 		ID:        id,
// 		Username:  req.Name,
// 		Email:     req.Email,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}, nil
// }
//
// func (s *UserService) Update(ctx context.Context, id int, req model.UpdateUserReq) (*model.User, error) {
// 	now := time.Now()
//
// 	var user model.User
// 	err := s.db.QueryRow(ctx,
// 		"UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4 RETURNING id, name, email, password, created_at, updated_at",
// 		req.Name, req.Email, now, id).
// 		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
//
// 	if err == pgx.ErrNoRows {
// 		return nil, errors.New("user not found")
// 	}
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &user, nil
// }
//
// func (s *UserService) Delete(ctx context.Context, id int) error {
// 	result, err := s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
// 	if err != nil {
// 		return err
// 	}
//
// 	if result.RowsAffected() == 0 {
// 		return errors.New("user not found")
// 	}
//
// 	return nil
// }
