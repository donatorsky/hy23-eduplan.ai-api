CREATE TABLE main.universities
(
	id          INTEGER             NOT NULL
		CONSTRAINT universities_pk
			PRIMARY KEY AUTOINCREMENT,
	name        TEXT COLLATE NOCASE NOT NULL
		CONSTRAINT universities_name_uq
			UNIQUE,
	link        TEXT                NOT NULL,
	voivodeship TEXT                NOT NULL,
	city        TEXT                NOT NULL
);

CREATE TABLE main.specializations
(
	id            INTEGER
		CONSTRAINT specializations_pk
			PRIMARY KEY AUTOINCREMENT,
	university_id INTEGER NOT NULL
		CONSTRAINT specializations_universities_id_fk
			REFERENCES main.universities
			ON UPDATE CASCADE ON DELETE CASCADE,
	name          TEXT    NOT NULL,
	link          TEXT    NOT NULL,
	level         TEXT    NOT NULL,
	type          TEXT    NOT NULL,
	profile       TEXT    NOT NULL,
	description   TEXT    NOT NULL
);
