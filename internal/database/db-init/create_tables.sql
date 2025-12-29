CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users ( 
	id UUID PRIMARY KEY,
	username TEXT UNIQUE NOT NULL, 
	password_hash TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Uploads table
CREATE TABLE uploads (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) NULL,
	guest_id UUID UNIQUE NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Files table
CREATE TABLE files (
	id UUID PRIMARY KEY,
	upload_id UUID REFERENCES uploads(id) ON DELETE CASCADE,
	name TEXT NOT NULL,
	s3_key TEXT NOT NULL,
	size INTEGER NOT NULL,
	mime_type TEXT NOT NULL,
	original_timestamp TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Rulesets table
CREATE TABLE rulesets (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) NULL, -- nullable for guests
	name TEXT NOT NULL, 
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Rules table
CREATE TABLE rules (
	id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(id) NULL,
	name TEXT NOT NULL,
	conditions_json JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	UNIQUE (user_id, name)
);

-- Rule bindings table
CREATE TABLE rule_bindings (
	id UUID PRIMARY KEY,
	ruleset_id UUID REFERENCES rulesets(id) ON DELETE CASCADE,
	rule_id UUID REFERENCES rules(id) ON DELETE CASCADE,
	target_id UUID NOT NULL UNIQUE,
	target_name TEXT NOT NULL
);

-- Tasks table
CREATE TABLE tasks (
	id UUID PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description TEXT,
	user_id UUID REFERENCES users(id) NULL, -- nullable for guests
	ruleset_id UUID REFERENCES rulesets(id),
	status VARCHAR(50) NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
);
