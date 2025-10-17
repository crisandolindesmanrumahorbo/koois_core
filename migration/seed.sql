INSERT INTO quizzes (title, description, author_id)
VALUES ('Math Basics', 'Basic arithmetic quiz', 1);

INSERT INTO questions (quiz_id, question_text, question_type)
VALUES (1, 'What is 2 + 2?', 'multiple_choice');

INSERT INTO question_options (question_id, option_text, is_correct)
VALUES
  (1, '3', false),
  (1, '4', true),
  (1, '5', false);
