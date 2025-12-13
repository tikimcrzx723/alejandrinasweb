CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    username VARCHAR(30) UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated BOOL NOT NULL,
    is_block BOOL NOT NULL,
    banner TEXT DEFAULT '',
    avatar TEXT DEFAULT '',
    created_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM now())::int,
    updated_at INTEGER DEFAULT 0,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS roles(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS roles_permissions(
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS users_roles(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry INTEGER NOT NULL,
    scope text NOT NULL
);

INSERT INTO roles(name) VALUES('USER');
INSERT INTO roles(name) VALUES('ADMIN');

INSERT INTO permissions(name) VALUES('CREATE');
INSERT INTO permissions(name) VALUES('READ');
INSERT INTO permissions(name) VALUES('UPDATE');
INSERT INTO permissions(name) VALUES('DELETE');

INSERT INTO roles_permissions(role_id, permission_id) VALUES(1, 2);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 1);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 2);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 3);
INSERT INTO roles_permissions(role_id, permission_id) VALUES(2, 4);
