-- name: CreateBook :exec
INSERT INTO books (
    id, name, description, author, isbn, genre, cover_image_file_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: GetBookByID :one
SELECT * FROM books 
WHERE id = $1 LIMIT 1;

-- name: GetAllBooks :many
SELECT * FROM books 
ORDER BY created_at DESC;

-- name: GetBooksWithPagination :many
SELECT * FROM books 
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2;