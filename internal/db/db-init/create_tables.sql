CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users ( 
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL, 
	password_hash TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Uploads table
CREATE TABLE uploads (
	id SERIAL PRIMARY KEY,
	upload_uuid UUID NOT NULL UNIQUE,
	user_id INTEGER REFERENCES users(id) NULL, -- nullable for guests
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Files table
CREATE TABLE files (
	id SERIAL PRIMARY KEY,
	file_uuid UUID NOT NULL UNIQUE,
	upload_id INTEGER REFERENCES uploads(id) ON DELETE CASCADE,
	upload_uuid UUID REFERENCES uploads(upload_uuid) ON DELETE CASCADE,
	file_name TEXT NOT NULL,
	s3_key TEXT NOT NULL,
	size INTEGER NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Rulesets table
CREATE TABLE rulesets (
	id SERIAL PRIMARY KEY,
	ruleset_uuid UUID NOT NULL UNIQUE,
	user_id INTEGER REFERENCES users(id) NULL, -- nullable for guests
	name TEXT NOT NULL, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Rules table
CREATE TABLE rules (
	id SERIAL PRIMARY KEY,
	rule_uuid UUID NOT NULL UNIQUE,
	user_id INTEGER REFERENCES users(id) NULL,
	name TEXT NOT NULL,
	conditions_json JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (user_id, name)
);

-- Rule bindings table
CREATE TABLE rule_bindings (
	id SERIAL PRIMARY KEY,
	ruleset_id INTEGER REFERENCES rulesets(id) ON DELETE CASCADE,
	rule_id INTEGER REFERENCES rules(id) ON DELETE CASCADE,
	target_uuid UUID NOT NULL UNIQUE,
	target_name TEXT NOT NULL
);

-- Jobs table
CREATE TABLE jobs (
	id SERIAL PRIMARY KEY,
	job_uuid UUID NOT NULL UNIQUE,
	user_id INTEGER REFERENCES users(id) NULL, -- nullable for guests
	upload_id INTEGER REFERENCES uploads(id),
	ruleset_id INTEGER REFERENCES rulesets(id),
	status TEXT NOT NULL,
	result_path TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
);
