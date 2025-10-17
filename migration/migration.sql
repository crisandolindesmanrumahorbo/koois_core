-- QUIZZES
CREATE TABLE quizzes (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  author_id INT REFERENCES users(user_id) ON DELETE SET NULL,
  state VARCHAR(50) DEFAULT 'draft' CHECK (state IN ('draft', 'published')),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- QUESTIONS
CREATE TABLE questions (
  id SERIAL PRIMARY KEY,
  quiz_id INT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
  question_text TEXT NOT NULL,
  question_type VARCHAR(50) NOT NULL CHECK (question_type IN (
    'multiple_choice',
    'true_false',
    'fill_blank',
    'matching',
    'open_ended'
  )),
  image_url TEXT, 
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- QUESTION OPTIONS
CREATE TABLE question_options (
  id SERIAL PRIMARY KEY,
  question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  option_text TEXT,
  image_url TEXT, 
  is_correct BOOLEAN DEFAULT FALSE,
  option_order INT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- QUIZ ATTEMPTS
CREATE TABLE quiz_attempts (
  id SERIAL PRIMARY KEY,
  quiz_id INT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
  user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  started_at TIMESTAMP DEFAULT NOW(),
  finished_at TIMESTAMP,
  score NUMERIC(5,2), 
  status VARCHAR(50) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'submitted', 'graded')),
  UNIQUE (quiz_id, user_id, started_at) -- allow multiple attempts if needed, but each has unique start
);

-- USER ANSWERS
CREATE TABLE user_answers (
  id SERIAL PRIMARY KEY,
  attempt_id INT NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
  question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  selected_option_id INT REFERENCES question_options(id) ON DELETE SET NULL,
  answer_text TEXT, 
  is_correct BOOLEAN,
  answered_at TIMESTAMP DEFAULT NOW(),
  UNIQUE (attempt_id, question_id)
);

-- INDEXES
CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);
CREATE INDEX idx_question_options_question_id ON question_options(question_id);
CREATE INDEX idx_quiz_attempts_user_id ON quiz_attempts(user_id);
CREATE INDEX idx_user_answers_attempt_id ON user_answers(attempt_id);
