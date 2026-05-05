//nolint:dupl // Adaptors are intended to follow a similar pattern.
package database

import (
	"context"

	"github.com/stretchr/testify/require"
)

func (p *PostgresAdapter) CleanBookPages(t require.TestingT) {
	_, err := p.pool.Exec(context.Background(), "TRUNCATE TABLE book_pages CASCADE;")
	require.NoError(t, err)
}

func (p *PostgresAdapter) CountBookPages(t require.TestingT) int64 {
	var count int64
	err := p.pool.QueryRow(context.Background(), "SELECT count(*) FROM book_pages").Scan(&count)
	require.NoError(t, err)
	return count
}

func (p *PostgresAdapter) GetBookPageByID(t require.TestingT, id string) TestBookPageRecord {
	var pRec TestBookPageRecord
	pRec.ID = id
	err := p.pool.QueryRow(
		context.Background(),
		"SELECT book_id, page_number, content, is_bookmarked, highlight, attached_image_file_id, created_at, updated_at FROM book_pages WHERE id = $1",
		id,
	).Scan(&pRec.BookID, &pRec.PageNumber, &pRec.Content, &pRec.IsBookmarked, &pRec.Highlight, &pRec.AttachedImageFileID, &pRec.CreatedAt, &pRec.UpdatedAt)
	require.NoError(t, err)
	return pRec
}

func (p *PostgresAdapter) SeedBookPages(t require.TestingT, pages []TestBookPageRecord) {
	query := `INSERT INTO book_pages (id, book_id, page_number, content, is_bookmarked, highlight, attached_image_file_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	for _, pRec := range pages {
		_, err := p.pool.Exec(
			context.Background(),
			query,
			pRec.ID,
			pRec.BookID,
			pRec.PageNumber,
			pRec.Content,
			pRec.IsBookmarked,
			pRec.Highlight,
			pRec.AttachedImageFileID,
			pRec.CreatedAt,
			pRec.UpdatedAt,
		)
		require.NoError(t, err)
	}
}
