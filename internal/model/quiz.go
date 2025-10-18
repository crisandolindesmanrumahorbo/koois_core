package model

import (
	"database/sql"
	"time"
)

type Quiz struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	AuthorId    int            `json:"author_id"`
	State       string         `json:"state"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Question struct {
	Id           int            `json:"id"`
	QuizId       int            `json:"quiz_id"`
	QuestionText string         `json:"question_text"`
	QuestionType string         `json:"question_type"`
	ImageUrl     sql.NullString `json:"image_url"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type QuestionOptions struct {
	Id          int            `json:"id"`
	QuestionId  int            `json:"question_id"`
	OptionText  sql.NullString `json:"option_text"`
	ImageUrl    sql.NullString `json:"image_url"`
	IsCorrect   bool           `json:"is_correct"`
	OptionOrder int            `json:"option_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type QuizAttempt struct {
	Id         int          `json:"id"`
	QuizId     int          `json:"quiz_id"`
	UserId     int          `json:"user_id"`
	StartedAt  time.Time    `json:"started_at"`
	FinishedAt sql.NullTime `json:"finished_at"`
	IsCorrect  bool         `json:"is_correct"`
	Score      int          `json:"score"`
	Status     string       `json:"status"`
}

// TODO add user_id to easier query GET
// and others GET query check field addition for more efficient query
type UserAnswer struct {
	Id               int            `json:"id"`
	AttemptId        int            `json:"attempt_id"`
	QuestionId       int            `json:"question_id"`
	SelectedOptionId int            `json:"selected_option_id"`
	AnswerText       sql.NullString `json:"answer_text"`
	IsCorrect        sql.NullBool   `json:"is_correct"`
	AnsweredAt       time.Time      `json:"answered_at"`
}

type GetQuestion struct {
	Id       int             `json:"id"`
	Text     string          `json:"text"`
	Type     string          `json:"type"`
	ImageUrl *sql.NullString `json:"image_url"`
	Options  []GetOption     `json:"options"`
}

type GetOption struct {
	Id        int             `json:"id"`
	Text      string          `json:"text"`
	ImageUrl  *sql.NullString `json:"image_url"`
	IsCorrect bool            `json:"is_correct"`
	Order     *sql.NullInt32  `json:"order"`
}

type CreateQuizReq struct {
	Title       string              `json:"title"`
	Description *string             `json:"description"`
	Questions   []CreateQuestionReq `json:"questions"`
}

type CreateQuestionReq struct {
	QuestionText string                 `json:"question_text"`
	QuestionType string                 `json:"question_type"`
	ImageUrl     *string                `json:"image_url"`
	Options      []CreateQuestionOption `json:"question_options"`
}

type CreateQuestionOption struct {
	OptionText string  `json:"option_text"`
	IsCorrect  bool    `json:"is_correct"`
	ImageUrl   *string `json:"image_url"`
}
