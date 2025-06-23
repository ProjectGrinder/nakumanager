CREATE TABLE views (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    team_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (team_id) REFERENCES teams(id)
);

CREATE TABLE view_group_bys (
    view_id INTEGER NOT NULL,
    group_by TEXT NOT NULL,
    FOREIGN KEY (view_id) REFERENCES views(id)
);

CREATE TABLE view_issues (
    view_id INTEGER NOT NULL,
    issue_id INTEGER NOT NULL,
    PRIMARY KEY (view_id, issue_id),
    FOREIGN KEY (view_id) REFERENCES views(id),
    FOREIGN KEY (issue_id) REFERENCES issues(id)
);
