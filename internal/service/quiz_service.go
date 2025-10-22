package service

import (
	"context"
	"encoding/json"
	"errors"
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
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	var quizID int
	err = tx.QueryRow(ctx,
		`INSERT INTO quizzes (title, description, author_id) VALUES ($1, $2, $3) RETURNING id`,
		quizReq.Title, quizReq.Description, authorId,
	).Scan(&quizID)
	if err != nil {
		return 0, fmt.Errorf("insert quiz: %w", err)
	}
	type BareQuestion struct {
		QuestionText string  `json:"question_text"`
		QuestionType string  `json:"question_type"`
		ImageUrl     *string `json:"image_url"`
	}
	var bare []BareQuestion
	for _, q := range quizReq.Questions {
		bare = append(bare, BareQuestion{QuestionText: q.QuestionText, QuestionType: q.QuestionType, ImageUrl: q.ImageUrl})
	}
	qJSON, _ := json.Marshal(bare)

	rows, err := tx.Query(ctx, `
		INSERT INTO questions (quiz_id, question_text, question_type, image_url)
		SELECT $1, q.question_text, q.question_type, q.image_url
		FROM jsonb_to_recordset($2::jsonb) AS q(question_text TEXT, question_type TEXT, image_url TEXT)
		RETURNING id;
	`, quizID, string(qJSON))
	if err != nil {
		return 0, fmt.Errorf("insert questions: %w", err)
	}
	defer rows.Close()

	var questionIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		questionIDs = append(questionIDs, id)
	}
	if len(questionIDs) != len(quizReq.Questions) {
		return 0, errors.New("mismatch between inserted question IDs and input count")
	}

	type FlatOption struct {
		QuestionID int     `json:"question_id"`
		OptionText *string `json:"option_text"`
		IsCorrect  bool    `json:"is_correct"`
		ImageUrl   *string `json:"image_url"`
	}
	var flat []FlatOption
	for i, q := range quizReq.Questions {
		for _, opt := range q.Options {
			flat = append(flat, FlatOption{
				QuestionID: questionIDs[i],
				OptionText: opt.OptionText,
				IsCorrect:  opt.IsCorrect,
				ImageUrl:   opt.ImageUrl,
			})
		}
	}

	if len(flat) > 0 {
		optJSON, _ := json.Marshal(flat)
		_, err = tx.Exec(ctx, `
			INSERT INTO question_options (question_id, option_text, is_correct, image_url)
			SELECT o.question_id, o.option_text, o.is_correct, o.image_url
			FROM jsonb_to_recordset($1::jsonb)
			AS o(question_id INT, option_text TEXT, is_correct BOOL, image_url TEXT);
		`, string(optJSON))
		if err != nil {
			return 0, fmt.Errorf("insert options: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return quizID, nil
}
