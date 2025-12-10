CREATE TABLE users ( 
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL, 
	password_hash TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE uploads (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES users(id), --nullable for guests--
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()	
);
