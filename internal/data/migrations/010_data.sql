-- +goose Up

INSERT INTO roles (name, description) VALUES ('admin', 'Administrator role with full access');
INSERT INTO user_roles (role, user) VALUES (1, 1);


--
-- TODO: remove this
-- 

INSERT INTO users (name, password_hash) VALUES ('admin', 'admin');


