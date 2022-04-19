CREATE TABLE IF NOT EXISTS auth (
	id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	phone varchar(20) NOT NULL,
	organization varchar(32) NOT NULL,
	code integer NOT NULL,
	request_time bigint NOT NULL,
	typeOfSign integer)
;

CREATE TABLE IF NOT EXISTS usersa0eebc999c0b4ef8bb6d6bb9bd380a11 (
	id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	name varchar(40) NOT NULL,
	phone varchar(20) NOT NULL,
	password varchar(255) NOT NULL) 
;

CREATE TABLE IF NOT EXISTS permissionsa0eebc999c0b4ef8bb6d6bb9bd380a11 (
	id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY)
;

CREATE TABLE IF NOT EXISTS users_permissionsa0eebc999c0b4ef8bb6d6bb9bd380a11 (
	id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	user_id integer,
	permission_id integer,

	CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES usersa0eebc999c0b4ef8bb6d6bb9bd380a11(id),
	CONSTRAINT permission_fk FOREIGN KEY (permission_id) REFERENCES permissionsa0eebc999c0b4ef8bb6d6bb9bd380a11(id),
	CONSTRAINT user_permission_unique UNIQUE (user_id, permission_id))
;

DROP TABLE IF EXISTS usersa0eebc999c0b4ef8bb6d6bb9bd380a11

SELECT * FROM usersa0eebc999c0b4ef8bb6d6bb9bd380a11