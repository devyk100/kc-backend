-- name: GetUserFromEmail :one
SELECT * FROM "User"
WHERE email = $1 LIMIT 1;

-- name: GetUserFromUsername :one
SELECT * FROM "User"
WHERE username = $1 LIMIT 1;

-- name: GetQuestionFromId :one
SELECT * FROM "Question"
WHERE id = $1 LIMIT 1;

-- name: GetAllTestcases :many
SELECT 
    q.id AS question_id, 
    q.body AS question_body, 
    q.driver_code, 
    q.email AS author_email, 
    t.id AS testcase_id, 
    t.input AS testcase_input, 
    t.output AS testcase_output, 
    t.order AS testcase_order
FROM "Question" q
LEFT JOIN "Testcases" t ON q.id = t.qid
WHERE q.id = $1
ORDER BY t.order ASC;

-- name: InsertSubmission :one
INSERT INTO "Submission" (
    "code",
    "message",
    "correct",
    "question_id",
    "language",
    "duration"
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;
