# Requirements Document — AA Big Book Study Group Secretary App

## 1. Purpose and Scope

A web application written in Go to support the secretary of an Alcoholics Anonymous Big Book study group. The secretary can manage one or more study groups, maintain a participant roster, distribute email communications, and record meeting notes. The primary users are the secretary (administrator) and, indirectly, study group participants who receive emails.

---

## 2. Data Models

### Study Group
- ID, name, description
- Created date, active/archived status

### Participant
- ID, foreign key to study group
- First name, last name
- Email address
- Telephone number
- Sobriety date
- Active/inactive status

### Meeting
- ID, foreign key to study group
- Meeting number (1–20, corresponding to the material sequence)
- Scheduled date and time
- Location (text: address or virtual link)
- Notes (free-form text entered after the meeting)
- Status: scheduled, completed, cancelled

### Email Log
- ID, foreign key to meeting or group
- Recipient list snapshot (JSON)
- Subject, body
- Sent timestamp, status (sent, failed)

---

## 3. Functional Requirements

### 3.1 Study Group Management
- Create a new study group with a name and description
- List all study groups
- Archive a study group (soft delete)
- Switch the active group context in the UI

### 3.2 Participant Management
- Add a participant to a study group
- Edit participant contact information and sobriety date
- Mark a participant inactive without deleting their record
- List all participants in a group, sortable by last name or sobriety date
- View a single participant's detail page

### 3.3 Meeting Management
- Create a meeting record specifying the meeting number, date/time, and location
- Edit the scheduled date/time and location before the meeting
- Enter meeting notes after the meeting concludes
- Mark a meeting as completed or cancelled
- List all meetings for a group in chronological order

### 3.4 Materials Pipeline
- Store one introduction markdown file and 20 meeting markdown files on the server filesystem (or embedded via `go:embed`)
- Retrieve the correct markdown file by meeting number
- Render markdown to HTML at send time for use in email bodies (use a library such as `goldmark`)

### 3.5 Email Communications

**Introduction email** — sent to a participant (or all participants) when they join a new group. Body is rendered from `intro.md`.

**Meeting notes + next meeting email** — sent after a meeting is marked complete. Contains:
- Meeting notes entered by the secretary
- Date, time, and location of the next scheduled meeting
- The study material for the upcoming meeting number, rendered from the corresponding markdown file

**Email composition screen** — before sending, the secretary reviews a preview of the rendered email and confirms the recipient list. Sending is not automatic.

**Email log** — every sent email is recorded with the recipients, subject, timestamp, and delivery status.

### 3.6 Dashboard
- Active group name and participant count
- Upcoming meeting date, number, and location
- List of recent email sends with status
- Quick links to add participant, record meeting notes, and send email

---

## 4. Non-Functional Requirements

### 4.1 Technology Constraints (Go-specific)
- Standard library `net/http` for the HTTP server; a lightweight router such as `chi` or `gorilla/mux` for path parameters
- `database/sql` with the `mattn/go-sqlite3` or `modernc.org/sqlite` driver for persistence (SQLite is appropriate for a single-secretary, low-traffic tool)
- `html/template` (or `text/template`) for server-rendered HTML pages
- `go:embed` to bundle markdown material files and static assets into the compiled binary
- `net/smtp` from the standard library for email, or `gomail.v2` for easier MIME handling; an SMTP relay (such as Mailgun, SendGrid, or a local mail server) configured via environment variables
- `goldmark` (or `blackfriday`) for markdown-to-HTML rendering
- Configuration loaded from environment variables or a `.env` file (use `godotenv` or standard `os.Getenv`)

### 4.2 Deployment
- Single compiled binary plus an SQLite database file
- Deployable to a Linux VPS, Fly.io, or Railway with minimal configuration
- SMTP credentials, base URL, and port configurable via environment variables
- No external services required beyond an SMTP relay

### 4.3 Security
- Simple session-based authentication (one secretary account): a login form sets a signed cookie
- No participant-facing login is required in the initial version
- All routes except `/login` require a valid session
- CSRF protection on all forms (use the `gorilla/csrf` package or implement the double-submit cookie pattern)
- Participant email addresses are stored only in the database and never exposed in URLs

### 4.4 Performance and Scale
- Designed for a single group of up to 50 participants; no horizontal scaling requirement
- All pages should be server-side rendered; no JavaScript framework required

---

## 5. Project Structure (Suggested Go Package Layout)

```
/cmd/server           — main.go, entry point, server startup
/internal/group       — study group domain logic and repository
/internal/participant — participant domain logic and repository
/internal/meeting     — meeting domain logic and repository
/internal/email       — email composition and SMTP dispatch
/internal/materials   — markdown loading and rendering
/internal/auth        — session handling and login
/internal/db          — database connection, migrations
/web/templates        — HTML templates (go:embed)
/web/static           — CSS, minimal JS (go:embed)
/materials            — intro.md, meeting_01.md … meeting_20.md (go:embed)
```

---

## 6. Suggested Build Sequence

This ordering lets you get something working at each step, which is especially useful while learning Go.

1. Database schema and migrations (`/internal/db`)
2. Participant CRUD with list and detail pages
3. Study group creation and participant association
4. Meeting creation and the meeting list page
5. Meeting notes entry and the meeting detail page
6. Markdown file loading and rendering
7. Email composition, preview, and SMTP dispatch
8. Introduction email flow
9. Post-meeting email flow (notes + next meeting + material)
10. Authentication and session middleware
11. Email log and dashboard

---

## 7. Out of Scope for v1

- Participant self-registration or any participant-facing portal
- Automated scheduling or reminders
- File uploads (photos, documents)
- Multi-secretary roles or permissions
- SMS or messaging integrations