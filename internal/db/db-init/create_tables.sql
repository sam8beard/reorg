CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users ( 
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL, 
	password_hash TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE uploads (
	id SERIAL PRIMARY KEY,
	upload_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_id INTEGER REFERENCES users(id) NULL, --nullable for guests--
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE files (
	id SERIAL PRIMARY KEY,
	upload_id INTEGER REFERENCES uploads(id),
	s3_key TEXT NOT NULL,
	size INTEGER NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rulesets (
	id SERIAL PRIMARY KEY,
	ruleset_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_id INTEGER REFERENCES users(id) NULL, --nullable for guests--
	name TEXT NOT NULL, 
	rules_json JSONB NOT NUll,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE jobs (
	id SERIAL PRIMARY KEY,
	job_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_id INTEGER REFERENCES users(id) NULL, --nullable for guests--
	upload_id INTEGER REFERENCES uploads(id),
	ruleset_json JSONB NOT NULL,
	status TEXT NOT NULL,
	result_path TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
);
