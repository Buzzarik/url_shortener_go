CREATE TABLE IF NOT EXISTS urls (
	id bigserial PRIMARY KEY,
	url text NOT NULL,
	alias text UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_urls_alias ON urls (alias);