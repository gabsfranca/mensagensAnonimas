# 📬 Anonymous Reporting Platform

> Full stack platform for receiving, tracking and triaging anonymous workplace reports.

**This system was built for a real compliance requirement and is actively used in production.**

---

## 📌 Context

Brazilian labor regulation **NR1** requires companies to provide a safe, anonymous channel 
for employees to report workplace incidents and concerns.

Employees submit reports anonymously, attach media, and track progress using only a short 
protocol ID — no account, no personal data collected. The committee accesses a separate 
authenticated panel to analyze, classify and respond to each case.

The anonymity model is intentional: reporters never create accounts. Each submission generates 
a unique short ID that acts as the only link between the reporter and their case — ensuring 
privacy while still allowing two-way communication.

---

## ✨ Features

### Public area
- Onboarding screen with usage guidelines and legal context
- Anonymous report submission with text and attachments (image, video, audio)
- Automatic generation of a short protocol ID for tracking
- Report status and observation history lookup by protocol ID
- Ability to add new observations to an existing report

### Admin area
- Admin registration and login
- JWT authentication with protected routes
- Report listing, filtering and detail view
- Status management: `received → under review → resolved`
- Tag assignment and removal for incident classification
- Internal observation logging throughout the investigation
- Dashboard with report counts by category

---

## 🛠️ Tech Stack

| Layer | Tech |
|---|---|
| Frontend | SolidJS, TypeScript, Vite, ApexCharts |
| Backend | Go, Gin, GORM, JWT |
| Database | PostgreSQL |
| Infra | Docker, Docker Compose, Nginx |

---

## 🏗️ Architecture

The project is split into two main services:

```text
frontend/   → user interface and admin panel
backend/    → API, business logic, authentication and persistence
```

Backend layer structure:

- `handler/` — HTTP entry points and response serialization
- `service/` — business rules
- `repo/` — data access
- `models/` — entities and relationships
- `middleware/` — authentication and route protection
- `storage/` — local file persistence for uploaded attachments

---

## 🔄 Product Flow

1. User accepts the channel's terms of use
2. Submits an anonymous report with or without attachments
3. Backend persists the report and generates a unique short protocol ID
4. User uses the ID to track investigation progress
5. Admin team accesses the authenticated panel
6. Report is classified, updated and commented throughout the investigation

---

## 🗃️ Data Model

- `Report` — report content, status, short protocol ID, attachments and tags
- `Media` — files submitted with the report
- `Observation` — history of updates and comments
- `Tag` — incident categorization
- `Admin` — users with access to the internal panel

---

## 🚀 Running locally

### Recommended: Docker Compose

```bash
git clone https://github.com/gabsfranca/mensagensAnonimas
cd mensagensAnonimas
docker compose up --build
```

Services:
- Frontend: `http://localhost`
- Backend: `http://localhost:8080`
- PostgreSQL: dedicated container

### Manual setup

#### Backend

Requirements: Go `1.23+`, PostgreSQL

```env
PORT=:8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=reports
URL=http://localhost:3000
```

```bash
cd backend
go mod download
go run ./cmd
```

#### Frontend

Requirements: Node.js, npm

```env
VITE_BACKEND_URL=http://localhost:8080
```

```bash
cd frontend
npm install
npm run dev
```

---

## 🔌 API Endpoints

### Public

| Method | Endpoint | Description |
|---|---|---|
| POST | `/send-anonymous-message` | Submit a new report |
| GET | `/reports/:id/status` | Get report status |
| GET | `/reports/:id/observations` | Get report observations |
| POST | `/reports/:id/observations` | Add observation to report |

### Admin (authenticated)

| Method | Endpoint | Description |
|---|---|---|
| POST | `/register` | Register admin |
| POST | `/login` | Login |
| POST | `/logout` | Logout |
| GET | `/messages` | List all reports |
| GET | `/messages/:id` | Get report detail |
| PATCH | `/messages/:id/status` | Update report status |
| POST | `/messages/reports/:id/tags` | Add tag to report |
| DELETE | `/messages/:messageId/tags/:tagId` | Remove tag |
| GET | `/messages/tags` | List all tags |
| GET | `/messages/tags/stats` | Tag statistics |

---

## 💡 Technical highlights

- Short unique protocol ID generated per report — reporters track cases without exposing UUIDs
- Automatic DB migration and tag seed on application startup
- Connection retry logic for improved robustness in containerized environments
- Multimedia attachment support with file type validation and size limits
- Admin panel with filters, observations and tag-based classification

---

## 🔮 Future improvements

- Automated tests for backend and frontend
- Session refresh and more robust security policy
- Cloud storage for uploaded files
- CI/CD pipeline with build, lint and test stages
- Advanced pagination and search in the admin panel
