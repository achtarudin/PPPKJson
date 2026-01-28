-- Create exam_sessions table
CREATE TABLE IF NOT EXISTS exam_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    session_code VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'NOT_STARTED',
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    duration INTEGER DEFAULT 120, -- Duration in minutes
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for exam_sessions
CREATE INDEX IF NOT EXISTS idx_exam_sessions_user_id ON exam_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_exam_sessions_session_code ON exam_sessions(session_code);
CREATE INDEX IF NOT EXISTS idx_exam_sessions_status ON exam_sessions(status);
CREATE INDEX IF NOT EXISTS idx_exam_sessions_deleted_at ON exam_sessions(deleted_at);

-- Create exam_questions table (assigns specific questions to each exam session)
CREATE TABLE IF NOT EXISTS exam_questions (
    id BIGSERIAL PRIMARY KEY,
    exam_session_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    category VARCHAR(50) NOT NULL,
    order_number INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_exam_questions_exam_session 
    FOREIGN KEY (exam_session_id) 
    REFERENCES exam_sessions(id) 
    ON DELETE CASCADE,
    
    CONSTRAINT fk_exam_questions_question 
    FOREIGN KEY (question_id) 
    REFERENCES questions(id) 
    ON DELETE CASCADE
);

-- Create indexes for exam_questions
CREATE INDEX IF NOT EXISTS idx_exam_questions_exam_session_id ON exam_questions(exam_session_id);
CREATE INDEX IF NOT EXISTS idx_exam_questions_question_id ON exam_questions(question_id);
CREATE INDEX IF NOT EXISTS idx_exam_questions_category ON exam_questions(category);
CREATE INDEX IF NOT EXISTS idx_exam_questions_order_number ON exam_questions(order_number);
CREATE INDEX IF NOT EXISTS idx_exam_questions_deleted_at ON exam_questions(deleted_at);

-- Create user_answers table
CREATE TABLE IF NOT EXISTS user_answers (
    id BIGSERIAL PRIMARY KEY,
    exam_session_id BIGINT NOT NULL,
    exam_question_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    question_option_id BIGINT NOT NULL,
    score INTEGER NOT NULL,
    answered_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_user_answers_exam_session 
    FOREIGN KEY (exam_session_id) 
    REFERENCES exam_sessions(id) 
    ON DELETE CASCADE,
    
    CONSTRAINT fk_user_answers_exam_question 
    FOREIGN KEY (exam_question_id) 
    REFERENCES exam_questions(id) 
    ON DELETE CASCADE,
    
    CONSTRAINT fk_user_answers_question 
    FOREIGN KEY (question_id) 
    REFERENCES questions(id) 
    ON DELETE CASCADE,
    
    CONSTRAINT fk_user_answers_question_option 
    FOREIGN KEY (question_option_id) 
    REFERENCES question_options(id) 
    ON DELETE CASCADE
);

-- Create indexes for user_answers
CREATE INDEX IF NOT EXISTS idx_user_answers_exam_session_id ON user_answers(exam_session_id);
CREATE INDEX IF NOT EXISTS idx_user_answers_exam_question_id ON user_answers(exam_question_id);
CREATE INDEX IF NOT EXISTS idx_user_answers_question_id ON user_answers(question_id);
CREATE INDEX IF NOT EXISTS idx_user_answers_question_option_id ON user_answers(question_option_id);
CREATE INDEX IF NOT EXISTS idx_user_answers_deleted_at ON user_answers(deleted_at);

-- Create exam_results table (results per category)
CREATE TABLE IF NOT EXISTS exam_results (
    id BIGSERIAL PRIMARY KEY,
    exam_session_id BIGINT NOT NULL,
    category VARCHAR(50) NOT NULL,
    total_questions INTEGER DEFAULT 5,
    total_answered INTEGER NOT NULL,
    total_score INTEGER NOT NULL,
    max_score INTEGER DEFAULT 20, -- 5 questions * 4 max score
    percentage DECIMAL(5,2) NOT NULL,
    grade VARCHAR(5),
    is_passed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_exam_results_exam_session 
    FOREIGN KEY (exam_session_id) 
    REFERENCES exam_sessions(id) 
    ON DELETE CASCADE
);

-- Create indexes for exam_results
CREATE INDEX IF NOT EXISTS idx_exam_results_exam_session_id ON exam_results(exam_session_id);
CREATE INDEX IF NOT EXISTS idx_exam_results_category ON exam_results(category);
CREATE INDEX IF NOT EXISTS idx_exam_results_deleted_at ON exam_results(deleted_at);

-- Create exam_summaries table (overall exam summary)
CREATE TABLE IF NOT EXISTS exam_summaries (
    id BIGSERIAL PRIMARY KEY,
    exam_session_id BIGINT NOT NULL UNIQUE,
    user_id VARCHAR(50) NOT NULL,
    total_questions INTEGER DEFAULT 20, -- 4 categories * 5 questions
    total_answered INTEGER NOT NULL,
    total_score INTEGER NOT NULL,
    max_score INTEGER DEFAULT 80, -- 20 questions * 4 max score
    overall_percentage DECIMAL(5,2) NOT NULL,
    overall_grade VARCHAR(5),
    is_passed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Foreign key constraints
    CONSTRAINT fk_exam_summaries_exam_session 
    FOREIGN KEY (exam_session_id) 
    REFERENCES exam_sessions(id) 
    ON DELETE CASCADE
);

-- Create indexes for exam_summaries
CREATE INDEX IF NOT EXISTS idx_exam_summaries_exam_session_id ON exam_summaries(exam_session_id);
CREATE INDEX IF NOT EXISTS idx_exam_summaries_user_id ON exam_summaries(user_id);
CREATE INDEX IF NOT EXISTS idx_exam_summaries_deleted_at ON exam_summaries(deleted_at);

-- Add unique constraint to prevent duplicate exam sessions per user per day
CREATE UNIQUE INDEX IF NOT EXISTS idx_exam_sessions_user_date ON exam_sessions(user_id, DATE(created_at)) WHERE deleted_at IS NULL;
