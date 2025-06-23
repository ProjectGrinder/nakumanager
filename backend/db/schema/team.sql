CREATE TABLE teams (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    leader_id INTEGER,
    FOREIGN KEY (leader_id) REFERENCES users(id)
);

CREATE TABLE team_members (
    team_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (team_id, user_id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
