-- name: CreateBookPage :exec
INSERT INTO book_pages (
    id, book_id, page_number, content, is_bookmarked, highlight, attached_image_file_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: GetBookPageByID :one
SELECT * FROM book_pages 
WHERE id = $1 LIMIT 1;

-- name: GetBookPagesByBookID :many
SELECT * FROM book_pages 
WHERE book_id = $1 
ORDER BY page_number ASC;

-- name: GetBookpagesByBookIDWithPagination :many
SELECT * FROM book_pages 
WHERE book_id = $1 
ORDER BY page_number ASC 
LIMIT $2 OFFSET $3;

-- name: GetBookpagesByPageRange :many
SELECT * FROM book_pages 
WHERE book_id = $1 AND page_number >= $2 AND page_number <= $3
ORDER BY page_number ASC;
