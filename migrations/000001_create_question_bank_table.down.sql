-- 1. Hapus tabel child terlebih dahulu (karena memiliki Foreign Key ke parent)
DROP TABLE IF EXISTS question_options;

-- 2. Hapus tabel parent
DROP TABLE IF EXISTS questions;