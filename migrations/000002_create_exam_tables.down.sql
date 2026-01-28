-- Drop tables in reverse order due to foreign key constraints

DROP INDEX IF EXISTS idx_exam_sessions_user_date;

-- Drop exam_summaries table
DROP INDEX IF EXISTS idx_exam_summaries_deleted_at;
DROP INDEX IF EXISTS idx_exam_summaries_user_id;
DROP INDEX IF EXISTS idx_exam_summaries_exam_session_id;
DROP TABLE IF EXISTS exam_summaries;

-- Drop exam_results table
DROP INDEX IF EXISTS idx_exam_results_deleted_at;
DROP INDEX IF EXISTS idx_exam_results_category;
DROP INDEX IF EXISTS idx_exam_results_exam_session_id;
DROP TABLE IF EXISTS exam_results;

-- Drop user_answers table
DROP INDEX IF EXISTS idx_user_answers_deleted_at;
DROP INDEX IF EXISTS idx_user_answers_question_option_id;
DROP INDEX IF EXISTS idx_user_answers_question_id;
DROP INDEX IF EXISTS idx_user_answers_exam_question_id;
DROP INDEX IF EXISTS idx_user_answers_exam_session_id;
DROP TABLE IF EXISTS user_answers;

-- Drop exam_questions table
DROP INDEX IF EXISTS idx_exam_questions_deleted_at;
DROP INDEX IF EXISTS idx_exam_questions_order_number;
DROP INDEX IF EXISTS idx_exam_questions_category;
DROP INDEX IF EXISTS idx_exam_questions_question_id;
DROP INDEX IF EXISTS idx_exam_questions_exam_session_id;
DROP TABLE IF EXISTS exam_questions;

-- Drop exam_sessions table
DROP INDEX IF EXISTS idx_exam_sessions_deleted_at;
DROP INDEX IF EXISTS idx_exam_sessions_status;
DROP INDEX IF EXISTS idx_exam_sessions_session_code;
DROP INDEX IF EXISTS idx_exam_sessions_user_id;
DROP TABLE IF EXISTS exam_sessions;
