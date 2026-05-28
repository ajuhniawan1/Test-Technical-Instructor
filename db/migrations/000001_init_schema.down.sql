-- Migration DOWN digunakan untuk rollback schema.
-- Urutan drop dibalik dari urutan create karena ada foreign key.

DROP TABLE IF EXISTS submission_histories;
DROP TABLE IF EXISTS submission_reviews;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS assignments;
DROP TABLE IF EXISTS class_talents;
DROP TABLE IF EXISTS class_trainers;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS users;
