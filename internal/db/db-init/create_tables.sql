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
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rules (
	id SERIAL PRIMARY KEY,
	rule_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_id INTEGER REFERENCES users(id) NULL,
	name TEXT NOT NULL,
	conditions_json JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (user_id, name)
);

CREATE TABLE rule_bindings (
	id SERIAL PRIMARY KEY,
	-- if the referenced rulesets row is deleted, delete all rows that reference it in this table --
	ruleset_id INTEGER REFERENCES rulesets(id) ON DELETE CASCADE,
	-- if the referenced rules row is deleted, delete all rows that reference it in this table --
	rule_id INTEGER REFERENCES rules(id) ON DELETE CASCADE,
	target_uuid UUID NOT NULL,
	target_name TEXT NOT NULL
);

CREATE TABLE jobs (
	id SERIAL PRIMARY KEY,
	job_uuid UUID NOT NULL DEFAULT gen_random_uuid(),
	user_id INTEGER REFERENCES users(id) NULL, --nullable for guests--
	upload_id INTEGER REFERENCES uploads(id),
	ruleset_id INTEGER REFERENCES rulesets(id),
	status TEXT NOT NULL,
	result_path TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
);
