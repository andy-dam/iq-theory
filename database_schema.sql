-- IQ Theory Database Schema
-- PostgreSQL Database Design

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false
);

-- Groups (for classrooms/teachers)
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    join_code VARCHAR(20) UNIQUE NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true,
    max_members INTEGER DEFAULT 100
);

-- Group memberships
CREATE TABLE group_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, group_id)
);

-- Friend relationships
CREATE TABLE friendships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    addressee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'blocked')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(requester_id, addressee_id),
    CHECK (requester_id != addressee_id)
);

-- Clef types (normalized parameter table)
CREATE TABLE clef_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL CHECK (name IN ('treble', 'bass', 'alto', 'tenor')),
    display_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true
);

-- Duration options (normalized parameter table)
CREATE TABLE duration_options (
    id SERIAL PRIMARY KEY,
    duration_seconds INTEGER UNIQUE NOT NULL CHECK (duration_seconds > 0),
    display_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true
);

-- Ledger line options (normalized parameter table)
CREATE TABLE ledger_line_options (
    id SERIAL PRIMARY KEY,
    max_lines INTEGER UNIQUE NOT NULL CHECK (max_lines >= 0 AND max_lines <= 3),
    display_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true
);

-- Quiz sessions (individual quiz attempts) - stores parameters directly
CREATE TABLE quiz_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Store the actual parameters instead of referencing a configuration
    clef VARCHAR(20) NOT NULL,
    duration_seconds INTEGER NOT NULL,
    max_ledger_lines INTEGER NOT NULL,
    
    score INTEGER NOT NULL DEFAULT 0,
    total_questions INTEGER NOT NULL,
    correct_answers INTEGER NOT NULL DEFAULT 0,
    time_taken_seconds INTEGER NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'abandoned')),
    accuracy_percentage DECIMAL(5,2) GENERATED ALWAYS AS (
        CASE 
            WHEN total_questions > 0 THEN (correct_answers::DECIMAL / total_questions) * 100
            ELSE 0
        END
    ) STORED,
    
    -- Add foreign key constraints to ensure valid parameters
    FOREIGN KEY (clef) REFERENCES clef_types(name),
    FOREIGN KEY (duration_seconds) REFERENCES duration_options(duration_seconds),
    FOREIGN KEY (max_ledger_lines) REFERENCES ledger_line_options(max_lines)
);

-- Individual quiz questions and answers
CREATE TABLE quiz_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_session_id UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_number INTEGER NOT NULL,
    note_image VARCHAR(100) NOT NULL, -- e.g., "A4.png"
    correct_note VARCHAR(10) NOT NULL, -- e.g., "A4"
    user_answer VARCHAR(10), -- User's guess
    is_correct BOOLEAN NOT NULL,
    time_taken_ms INTEGER NOT NULL,
    answered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(quiz_session_id, question_number)
);

-- View to dynamically generate all possible quiz configurations
CREATE VIEW available_quiz_configurations AS
SELECT 
    CONCAT(ct.display_name, ' - ', do.display_name, ' - ', llo.display_name) as configuration_name,
    ct.name as clef,
    ct.display_name as clef_display,
    do.duration_seconds,
    do.display_name as duration_display,
    llo.max_lines as max_ledger_lines,
    llo.display_name as ledger_display,
    ct.is_active AND do.is_active AND llo.is_active as is_available
FROM clef_types ct
CROSS JOIN duration_options do
CROSS JOIN ledger_line_options llo
WHERE ct.is_active = true 
  AND do.is_active = true 
  AND llo.is_active = true
ORDER BY ct.name, do.duration_seconds, llo.max_lines;

-- Leaderboards (materialized view for performance)
CREATE MATERIALIZED VIEW leaderboards AS
SELECT 
    qs.clef,
    qs.duration_seconds,
    qs.max_ledger_lines,
    CONCAT(
        (SELECT display_name FROM clef_types WHERE name = qs.clef), 
        ' - ', 
        (SELECT display_name FROM duration_options WHERE duration_seconds = qs.duration_seconds),
        ' - ',
        (SELECT display_name FROM ledger_line_options WHERE max_lines = qs.max_ledger_lines)
    ) as quiz_name,
    u.id as user_id,
    u.username,
    u.display_name,
    MAX(qs.score) as best_score,
    MAX(qs.accuracy_percentage) as best_accuracy,
    MIN(qs.time_taken_seconds) as fastest_time,
    COUNT(qs.id) as total_attempts,
    AVG(qs.score) as average_score,
    MAX(qs.completed_at) as last_attempt,
    RANK() OVER (
        PARTITION BY qs.clef, qs.duration_seconds, qs.max_ledger_lines 
        ORDER BY MAX(qs.score) DESC, MIN(qs.time_taken_seconds) ASC
    ) as global_rank
FROM quiz_sessions qs
JOIN users u ON qs.user_id = u.id
WHERE qs.status = 'completed'
GROUP BY qs.clef, qs.duration_seconds, qs.max_ledger_lines, u.id, u.username, u.display_name;

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_group_memberships_user_id ON group_memberships(user_id);
CREATE INDEX idx_group_memberships_group_id ON group_memberships(group_id);
CREATE INDEX idx_friendships_requester ON friendships(requester_id);
CREATE INDEX idx_friendships_addressee ON friendships(addressee_id);
CREATE INDEX idx_quiz_sessions_user_id ON quiz_sessions(user_id);
CREATE INDEX idx_quiz_sessions_clef ON quiz_sessions(clef);
CREATE INDEX idx_quiz_sessions_params ON quiz_sessions(clef, duration_seconds, max_ledger_lines);
CREATE INDEX idx_quiz_sessions_completed_at ON quiz_sessions(completed_at);
CREATE INDEX idx_quiz_answers_session_id ON quiz_answers(quiz_session_id);

-- Insert the base parameters (only 11 records total instead of 48 configurations)
INSERT INTO clef_types (name, display_name) VALUES
('treble', 'Treble Clef'),
('bass', 'Bass Clef'),
('alto', 'Alto Clef'),
('tenor', 'Tenor Clef');

INSERT INTO duration_options (duration_seconds, display_name) VALUES
(30, '30 seconds'),
(60, '1 minute'),
(120, '2 minutes');

INSERT INTO ledger_line_options (max_lines, display_name) VALUES
(0, 'No ledger lines'),
(1, 'Up to 1 ledger line'),
(2, 'Up to 2 ledger lines'),
(3, 'Up to 3 ledger lines');

-- Function to get quiz configuration display name
CREATE OR REPLACE FUNCTION get_quiz_config_name(
    p_clef VARCHAR(20),
    p_duration_seconds INTEGER,
    p_max_ledger_lines INTEGER
) RETURNS VARCHAR(150) AS $$
DECLARE
    clef_display VARCHAR(50);
    duration_display VARCHAR(50);
    ledger_display VARCHAR(50);
BEGIN
    SELECT ct.display_name INTO clef_display 
    FROM clef_types ct WHERE ct.name = p_clef;
    
    SELECT do.display_name INTO duration_display 
    FROM duration_options do WHERE do.duration_seconds = p_duration_seconds;
    
    SELECT llo.display_name INTO ledger_display 
    FROM ledger_line_options llo WHERE llo.max_lines = p_max_ledger_lines;
    
    RETURN CONCAT(clef_display, ' - ', duration_display, ' - ', ledger_display);
END;
$$ LANGUAGE plpgsql;

-- Function to refresh leaderboard materialized view
CREATE OR REPLACE FUNCTION refresh_leaderboards()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW leaderboards;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-refresh leaderboards when quiz sessions are completed
CREATE OR REPLACE FUNCTION trigger_refresh_leaderboards()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'completed' AND (OLD.status IS NULL OR OLD.status != 'completed') THEN
        PERFORM refresh_leaderboards();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER quiz_completion_trigger
    AFTER UPDATE ON quiz_sessions
    FOR EACH ROW
    EXECUTE FUNCTION trigger_refresh_leaderboards();
