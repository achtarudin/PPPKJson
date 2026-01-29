-- Test data for grade verification
-- Create exam sessions and results for coba1 (Grade A), coba2 (Grade B), coba3 (Grade C)

-- Clean existing test data
DELETE FROM exam_answers WHERE exam_session_id IN (
    SELECT id FROM exam_sessions WHERE user_id IN ('coba1', 'coba2', 'coba3')
);
DELETE FROM exam_results_by_category WHERE exam_session_id IN (
    SELECT id FROM exam_sessions WHERE user_id IN ('coba1', 'coba2', 'coba3')
);
DELETE FROM exam_results WHERE exam_session_id IN (
    SELECT id FROM exam_sessions WHERE user_id IN ('coba1', 'coba2', 'coba3')
);
DELETE FROM exam_questions WHERE exam_session_id IN (
    SELECT id FROM exam_sessions WHERE user_id IN ('coba1', 'coba2', 'coba3')
);
DELETE FROM exam_sessions WHERE user_id IN ('coba1', 'coba2', 'coba3');

-- Insert exam sessions
INSERT INTO exam_sessions (user_id, session_code, status, expires_at, duration, created_at, updated_at) VALUES
('coba1', 'EXAM_coba1_1769659900', 'COMPLETED', NOW() + INTERVAL '2 hours', 130, NOW(), NOW()),
('coba2', 'EXAM_coba2_1769659900', 'COMPLETED', NOW() + INTERVAL '2 hours', 130, NOW(), NOW()),
('coba3', 'EXAM_coba3_1769659900', 'COMPLETED', NOW() + INTERVAL '2 hours', 130, NOW(), NOW());

-- Insert exam results for COBA1 (Grade A - 100% = 690 points)
INSERT INTO exam_results (exam_session_id, user_id, total_questions, total_answered, total_score, max_score, overall_percentage, overall_grade, is_passed, completed_at, created_at, updated_at)
SELECT id, 'coba1', 145, 145, 690, 690, 100.0, 'A', true, NOW(), NOW(), NOW()
FROM exam_sessions WHERE user_id = 'coba1';

-- Insert category results for COBA1
INSERT INTO exam_results_by_category (exam_session_id, category, total_questions, total_answered, total_score, max_score, percentage, grade, is_passed, created_at, updated_at)
SELECT id, 'TEKNIS', 90, 90, 450, 450, 100.0, 'A', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba1'
UNION ALL
SELECT id, 'MANAJERIAL', 25, 25, 100, 100, 100.0, 'A', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba1'
UNION ALL
SELECT id, 'SOSIAL KULTURAL', 20, 20, 100, 100, 100.0, 'A', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba1'
UNION ALL
SELECT id, 'WAWANCARA', 10, 10, 40, 40, 100.0, 'A', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba1';

-- Insert exam results for COBA2 (Grade B - 90% = 621 points)
INSERT INTO exam_results (exam_session_id, user_id, total_questions, total_answered, total_score, max_score, overall_percentage, overall_grade, is_passed, completed_at, created_at, updated_at)
SELECT id, 'coba2', 145, 145, 621, 690, 90.0, 'B', true, NOW(), NOW(), NOW()
FROM exam_sessions WHERE user_id = 'coba2';

-- Insert category results for COBA2
INSERT INTO exam_results_by_category (exam_session_id, category, total_questions, total_answered, total_score, max_score, percentage, grade, is_passed, created_at, updated_at)
SELECT id, 'TEKNIS', 90, 90, 405, 450, 90.0, 'B', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba2'
UNION ALL
SELECT id, 'MANAJERIAL', 25, 25, 90, 100, 90.0, 'B', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba2'
UNION ALL
SELECT id, 'SOSIAL KULTURAL', 20, 20, 90, 100, 90.0, 'B', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba2'
UNION ALL
SELECT id, 'WAWANCARA', 10, 10, 36, 40, 90.0, 'B', true, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba2';

-- Insert exam results for COBA3 (Grade C - 80% = 552 points)
INSERT INTO exam_results (exam_session_id, user_id, total_questions, total_answered, total_score, max_score, overall_percentage, overall_grade, is_passed, completed_at, created_at, updated_at)
SELECT id, 'coba3', 145, 145, 552, 690, 80.0, 'C', false, NOW(), NOW(), NOW()
FROM exam_sessions WHERE user_id = 'coba3';

-- Insert category results for COBA3
INSERT INTO exam_results_by_category (exam_session_id, category, total_questions, total_answered, total_score, max_score, percentage, grade, is_passed, created_at, updated_at)
SELECT id, 'TEKNIS', 90, 90, 360, 450, 80.0, 'C', false, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba3'
UNION ALL
SELECT id, 'MANAJERIAL', 25, 25, 80, 100, 80.0, 'C', false, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba3'
UNION ALL
SELECT id, 'SOSIAL KULTURAL', 20, 20, 80, 100, 80.0, 'C', false, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba3'
UNION ALL
SELECT id, 'WAWANCARA', 10, 10, 32, 40, 80.0, 'C', false, NOW(), NOW() FROM exam_sessions WHERE user_id = 'coba3';

-- Verification queries
SELECT 'COBA1 Results' as test_user;
SELECT user_id, total_score, max_score, overall_percentage, overall_grade, is_passed FROM exam_results WHERE user_id = 'coba1';

SELECT 'COBA2 Results' as test_user;
SELECT user_id, total_score, max_score, overall_percentage, overall_grade, is_passed FROM exam_results WHERE user_id = 'coba2';

SELECT 'COBA3 Results' as test_user;
SELECT user_id, total_score, max_score, overall_percentage, overall_grade, is_passed FROM exam_results WHERE user_id = 'coba3';