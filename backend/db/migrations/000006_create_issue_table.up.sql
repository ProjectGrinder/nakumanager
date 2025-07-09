CREATE TABLE issues (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    priority TEXT CHECK(priority IN ('low', 'medium', 'high')),
    status TEXT CHECK(status IN ('todo', 'doing', 'done')) NOT NULL,
    project_id TEXT,
    team_id TEXT NOT NULL,
    start_date DATETIME,
    end_date DATETIME,
    label TEXT,
    owner_id TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE issue_assignees (
    issue_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (issue_id, user_id),
    FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
