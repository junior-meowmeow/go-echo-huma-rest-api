-- +goose Up
-- +goose StatementBegin
CREATE TABLE books (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    author TEXT NOT NULL DEFAULT '',
    isbn TEXT NOT NULL DEFAULT '',
    genre TEXT NOT NULL DEFAULT '',
    cover_image_file_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_books_created_at ON books (created_at DESC);

CREATE TABLE book_pages (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    page_number BIGINT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    is_bookmarked BOOLEAN NOT NULL DEFAULT FALSE,
    highlight TEXT NOT NULL DEFAULT '',
    attached_image_file_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    UNIQUE (book_id, page_number)
);

CREATE INDEX idx_book_pages_book_id_page_number ON book_pages (book_id, page_number);
CREATE INDEX idx_book_pages_book_id ON book_pages (book_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book_pages;
DROP TABLE IF EXISTS books;
-- +goose StatementEnd
