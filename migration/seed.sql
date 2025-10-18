INSERT INTO quizzes (title, description, author_id)
VALUES ('Math Basics', 'Basic arithmetic quiz', 3);

INSERT INTO questions (quiz_id, question_text, question_type)
VALUES (3, 'What is 2 + 2?', 'multiple_choice');

INSERT INTO question_options (question_id, option_text, is_correct)
VALUES
  (3, '3', false),
  (3, '4', true),
  (3, '5', false);

INSERT INTO questions (quiz_id, question_text, question_type)
VALUES (3, 'What is 2 + 2?', 'open_ended');

SELECT * 
FROM questions q
LEFT JOIN question_options qo ON q.id = qo.question_id
WHERE q.quiz_id = 3;
