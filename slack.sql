-- Table for storing user information
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    user_code VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing team information
CREATE TABLE IF NOT EXISTS teams (
    team_id SERIAL PRIMARY KEY,
    team_type VARCHAR(50) NOT NULL,
    team_intro TEXT,
    team_name VARCHAR(100) NOT NULL,
    team_leader INT REFERENCES users(user_id) ON DELETE CASCADE,
    team_description TEXT,
    num_members INT DEFAULT 0,
    team_etc TEXT,
    message_ts VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing tags
CREATE TABLE IF NOT EXISTS tags (
    tag_id SERIAL PRIMARY KEY,
    tag_name VARCHAR(50) NOT NULL UNIQUE
);

-- Junction table to associate users with teams (many-to-many relationship)
CREATE TABLE IF NOT EXISTS user_teams (
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    team_id INT REFERENCES teams(team_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, team_id)
);

-- Junction table to associate teams with tags (many-to-many relationship)
CREATE TABLE IF NOT EXISTS team_tags (
    team_id INT REFERENCES teams(team_id) ON DELETE CASCADE,
    tag_id INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (team_id, tag_id)
);
