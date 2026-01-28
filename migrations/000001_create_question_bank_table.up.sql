-- 1. Membuat tabel 'questions'
CREATE TABLE IF NOT EXISTS questions (
    id BIGSERIAL PRIMARY KEY,           -- BIGSERIAL otomatis membuat Sequence untuk auto-increment
    category VARCHAR(100) NOT NULL,
    question_text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Membuat Index (Postgres menggunakan CONCURRENTLY jika ingin non-blocking, tapi untuk migrasi awal biasa saja cukup)
CREATE INDEX IF NOT EXISTS idx_questions_deleted_at ON questions(deleted_at);

-- ---------------------------------------------------------

-- 2. Membuat tabel 'question_options'
CREATE TABLE IF NOT EXISTS question_options (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL,        
    option_text TEXT NOT NULL,
    score INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Definisi Foreign Key dengan ON DELETE CASCADE
    CONSTRAINT fk_questions_options 
    FOREIGN KEY (question_id) 
    REFERENCES questions(id) 
    ON DELETE CASCADE
);

-- Membuat Index untuk question_options
CREATE INDEX IF NOT EXISTS idx_question_options_question_id ON question_options(question_id);
CREATE INDEX IF NOT EXISTS idx_question_options_deleted_at ON question_options(deleted_at);