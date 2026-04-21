CREATE TABLE IF NOT EXISTS study_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    archived INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_group_id INTEGER NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    telephone TEXT NOT NULL DEFAULT '',
    sobriety_date TEXT,
    active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS meetings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_group_id INTEGER NOT NULL,
    meeting_number INTEGER NOT NULL,
    scheduled_at TEXT NOT NULL,
    location TEXT NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'scheduled',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE RESTRICT,
    CHECK (meeting_number BETWEEN 1 AND 20),
    CHECK (status IN ('scheduled', 'completed', 'cancelled'))
);

CREATE TABLE IF NOT EXISTS email_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    study_group_id INTEGER,
    meeting_id INTEGER,
    recipient_snapshot TEXT NOT NULL,
    subject TEXT NOT NULL,
    body TEXT NOT NULL,
    sent_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    status TEXT NOT NULL DEFAULT 'sent',
    FOREIGN KEY (study_group_id) REFERENCES study_groups(id) ON DELETE SET NULL,
    FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE SET NULL,
    CHECK (status IN ('sent', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_participants_study_group_id ON participants(study_group_id);
CREATE INDEX IF NOT EXISTS idx_meetings_study_group_id ON meetings(study_group_id);
CREATE INDEX IF NOT EXISTS idx_meetings_scheduled_at ON meetings(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_email_logs_study_group_id ON email_logs(study_group_id);
CREATE INDEX IF NOT EXISTS idx_email_logs_meeting_id ON email_logs(meeting_id);