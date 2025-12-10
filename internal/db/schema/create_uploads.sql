CREATE TABLE uploads (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id), --nullable for guests--
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),	
)
