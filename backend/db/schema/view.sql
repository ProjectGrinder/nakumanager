CREATE TABLE views (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_by TEXT NOT NULL,
    team_id TEXT NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE view_group_bys (
    view_id TEXT NOT NULL,
    group_by TEXT NOT NULL CHECK (group_by IN ('status', 'assignee', 'priority', 'project_id', 'label', 'team_id', 'end_date')),
    PRIMARY KEY (view_id, group_by),
    FOREIGN KEY (view_id) REFERENCES views(id) ON DELETE CASCADE
);

CREATE TABLE view_issues (
    view_id TEXT NOT NULL,
    issue_id TEXT NOT NULL,
    PRIMARY KEY (view_id, issue_id),
    FOREIGN KEY (view_id) REFERENCES views(id) ON DELETE CASCADE,
    FOREIGN KEY (issue_id) REFERENCES issues(id)
);
