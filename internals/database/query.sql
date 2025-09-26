-- name: CreateAccount :one
INSERT INTO accounts (
  username, email, password, account_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAccounts :many
SELECT * FROM accounts ORDER BY username;

-- name: GetAccount :one
SELECT * FROM accounts WHERE account_id = $1 LIMIT 1;

-- name: VerifyEmail :one
UPDATE accounts SET email_verified = $2
WHERE account_id = $1 RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE account_id = $1;



-- name: CreateView :one
INSERT INTO views (
  title, paragraph, user_id, public, view_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetView :one
SELECT * FROM  views WHERE view_id = $1 LIMIT 1;

-- name: GetAllViews :many
SELECT * FROM views ORDER BY user_id;

-- name: GetUserViews :many
SELECT * FROM views WHERE user_id = $1;

-- name: GetAllPublicViews :many
SELECT * FROM views WHERE public = true ORDER BY created_at;

-- name: DeleteView :exec
DELETE FROM views WHERE view_id = $1;