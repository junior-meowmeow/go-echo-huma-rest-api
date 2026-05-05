//nolint:dupl // Adaptors are intended to follow a similar pattern.
package database

import (
	"context"

	"github.com/stretchr/testify/require"
)

func (p *PostgresAdapter) CleanBooks(t require.TestingT) {
	_, err := p.pool.Exec(context.Background(), "TRUNCATE TABLE books CASCADE;")
	require.NoError(t, err)
}

func (p *PostgresAdapter) CountBooks(t require.TestingT) int64 {
	var count int64
	err := p.pool.QueryRow(context.Background(), "SELECT count(*) FROM books").Scan(&count)
	require.NoError(t, err)
	return count
}

func (p *PostgresAdapter) GetBookByID(t require.TestingT, id string) TestBookRecord {
	var b TestBookRecord
	b.ID = id
	err := p.pool.QueryRow(context.Background(),
		"SELECT name, description, author, isbn, genre, cover_image_file_id, created_at, updated_at FROM books WHERE id = $1", id,
	).Scan(&b.Name, &b.Description, &b.Author, &b.ISBN, &b.Genre, &b.CoverImageFileID, &b.CreatedAt, &b.UpdatedAt)
	require.NoError(t, err)
	return b
}

func (p *PostgresAdapter) SeedBooks(t require.TestingT, books []TestBookRecord) {
	query := `INSERT INTO books (id, name, description, author, isbn, genre, cover_image_file_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	for _, b := range books {
		_, err := p.pool.Exec(
			context.Background(),
			query,
			b.ID,
			b.Name,
			b.Description,
			b.Author,
			b.ISBN,
			b.Genre,
			b.CoverImageFileID,
			b.CreatedAt,
			b.UpdatedAt,
		)
		require.NoError(t, err)
	}
}
