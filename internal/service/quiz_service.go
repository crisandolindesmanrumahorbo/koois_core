package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"koois_core/internal/model"
)

type QuizService struct {
	db *pgxpool.Pool
}

func NewQuizService(db *pgxpool.Pool) *QuizService {
	return &QuizService{db: db}
}

func (s *QuizService) GetQuizQuestions(ctx context.Context, quizID int) ([]*model.GetQuestion, error) {
	rows, err := s.db.Query(ctx, `
		SELECT 
			q.id AS question_id,
			q.question_text AS question_text,
			q.question_type,
			q.image_url,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT(
						'id', o.id,
						'text', o.option_text,
						'is_correct', o.is_correct
					)
				) FILTER (WHERE o.id IS NOT NULL),
				'[]'
			) AS options
		FROM questions q
		LEFT JOIN question_options o ON o.question_id = q.id
		WHERE q.quiz_id = $1
		GROUP BY q.id
	`, quizID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var result []*model.GetQuestion

	for rows.Next() {
		var q model.GetQuestion
		var optionsJSON []byte

		err := rows.Scan(&q.Id, &q.Text, &q.Type, &q.ImageUrl, &optionsJSON)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}

		result = append(result, &q)
	}

	return result, nil
}

func (s *QuizService) Create(ctx context.Context, quizReq model.CreateQuizReq, authorId int) (int, error) {
	var quizID int
	err := s.db.QueryRow(ctx,
		`INSERT INTO quizzes (title, description, author_id) VALUES ($1, $2, $3) RETURNING id`,
		quizReq.Title, quizReq.Description,
	).Scan(&quizID, &authorId)
	if err != nil {
		return 0, fmt.Errorf("insert quiz: %w", err)
	}
}
