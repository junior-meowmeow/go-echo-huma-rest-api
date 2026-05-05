package database

import (
	"github.com/stretchr/testify/require"
)

func (p *PostgresAdapter) CleanFiles(_ require.TestingT) {
}

func (p *PostgresAdapter) CountFiles(_ require.TestingT) int64 {
	return 0
}

func (p *PostgresAdapter) CleanReviews(_ require.TestingT) {
}

func (p *PostgresAdapter) CountReviews(_ require.TestingT) int64 {
	return 0
}

func (p *PostgresAdapter) GetAllReviews(_ require.TestingT) []TestReviewRecord {
	return nil
}

func (p *PostgresAdapter) SeedReviews(_ require.TestingT, _ []TestReviewRecord) {
}

func (p *PostgresAdapter) CleanUsers(_ require.TestingT) {
}

func (p *PostgresAdapter) CountUsers(_ require.TestingT) int64 {
	return 0
}

func (p *PostgresAdapter) GetUserByUsername(_ require.TestingT, _ string) TestUserRecord {
	return TestUserRecord{}
}
