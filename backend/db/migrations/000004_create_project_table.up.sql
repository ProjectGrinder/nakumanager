CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT,
    priority TEXT,
    workspace_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    leader_id TEXT,
    start_date DATETIME,
    end_date DATETIME,
    label TEXT,
    created_by TEXT NOT NULL,

    FOREIGN KEY (workspace_id) REFERENCES workspaces(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (leader_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE project_members (
    project_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (project_id, user_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
