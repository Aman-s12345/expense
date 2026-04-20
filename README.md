# Expense Tracker

A minimal full-stack personal finance tool to record and review expenses.

**Live:** [https://expense.amans.site](https://expense.amans.site)
**API:** [https://expense-tracker-api-1p0y.onrender.com](https://expense-tracker-api-1p0y.onrender.com)

---

## Features

- Create expense entries with amount, category, description, and date
- View a list of all expenses
- Filter expenses by category
- Sort by date (newest first / oldest first)
- See the total of currently visible expenses (e.g., "Total: ₹1,250.00")
- Category-wise summary with visual breakdown
- Guest login (no signup needed) and email/password authentication
- Idempotent expense creation — safe against double-clicks and network retries

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (Fiber), PostgreSQL, sqlx |
| Frontend | React 18, Vite, CSS Modules |
| Hosting | Render (Docker for backend, static site for frontend) |

---

## Project Structure

```
expense/
├── backend/
│   ├── main.go                    # Entry point
│   ├── Dockerfile                 # Multi-stage build
│   ├── config/                    # Environment-based configuration
│   ├── db/
│   │   ├── db.go                  # PostgreSQL connection pool
│   │   ├── migrations/001_init.sql
│   │   └── models/                # Expense, User, Session structs
│   ├── services/
│   │   ├── auth/                  # Register, login, guest, session management
│   │   └── expense/               # Create & list expenses
│   ├── middlewares/               # Auth middleware, rate limiting
│   ├── providers/                 # Dependency injection (wires stores → services)
│   ├── routes/                    # HTTP handlers (auth/, expense/)
│   └── cmd/migrate/               # Migration runner
├── frontend/
│   ├── src/
│   │   ├── lib/                   # API client, utility functions
│   │   ├── context/               # AuthContext (login state)
│   │   ├── hooks/                 # useExpenses, useCreateExpense
│   │   ├── components/            # Header, AddExpenseForm, ExpenseList, CategorySummary
│   │   └── pages/                 # LoginPage, DashboardPage
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
└── render.yaml                    # Render blueprint (DB + backend + frontend)
```

---

## Key Design Decisions

### Money as integers (paisa)

All monetary values are stored as `BIGINT` in the database, representing paisa (1/100th of a rupee). This avoids floating-point precision issues entirely. The frontend converts to/from rupees for display and input — `15050` paisa is shown as `₹150.50`.

### Idempotent expense creation

The client generates a UUID (`idempotency_key`) for each expense submission. This key is sent with the POST request. The database enforces a `UNIQUE(user_id, idempotency_key)` constraint, and the API uses `ON CONFLICT DO NOTHING`. If a user double-clicks submit or the browser retries due to a network hiccup, the same expense is returned instead of creating a duplicate.

### Session-based auth (not JWT)

I chose opaque session tokens stored in PostgreSQL over JWTs. Sessions can be revoked server-side on logout, there's no token-refresh complexity, and for a small app like this, a DB lookup per request is perfectly fine. Guest users get a session with a longer TTL (30 days vs 7 days for registered users).

### PostgreSQL over SQLite

The assignment allows any persistence mechanism. I chose PostgreSQL because it's what I'd use in production — it handles concurrent writes correctly, supports constraints (`CHECK`, `UNIQUE`) that enforce data integrity at the DB level, and Render offers a free managed instance. The `CHECK (amount > 0)` constraint means invalid data can never enter the database regardless of application bugs.

### Service / Store pattern

Each domain (auth, expense) follows a 4-file pattern: `service.go` (interface), `service_impl.go` (business logic), `store.go` (data access interface), `store_impl.go` (SQL queries). This separates business rules from database concerns and makes the code testable — you can mock the store to unit-test the service.

---

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/auth/register` | No | Register with email + password |
| `POST` | `/api/auth/login` | No | Login with email + password |
| `POST` | `/api/auth/guest` | No | Create a guest session |
| `POST` | `/api/auth/logout` | Yes | End current session |
| `GET` | `/api/auth/me` | Yes | Get current user |
| `POST` | `/api/expenses/` | Yes | Create an expense |
| `GET` | `/api/expenses/` | Yes | List expenses |
| `GET` | `/api/health` | No | Health check |

**GET /api/expenses/** supports query parameters:
- `category` — filter by category (e.g., `?category=Food`)
- `sort` — `date_desc` (default, newest first) or `date_asc`

The response includes a computed `total_paisa` for the filtered set.

---

## Running Locally

### Backend

```bash
cd backend
cp .env.example .env        # edit DATABASE_URL if needed
go mod tidy
go run cmd/migrate/main.go  # run migrations
go run .                     # starts on :3000
```

### Frontend

```bash
cd frontend
npm install
npm run dev                  # starts on :5173, proxies /api to :3000
```

---

## Trade-offs Due to Timebox

- **No pagination** — the expense list returns all records. For a personal tracker this is fine for hundreds of entries; at scale, cursor-based pagination would be needed.
- **No edit/delete** — the assignment only requires create and list. Adding these would be straightforward (PUT/DELETE endpoints, soft-delete column).
- **No refresh tokens** — sessions expire after 7 days and the user re-authenticates. A refresh-token flow would improve UX for long sessions.
- **Minimal frontend styling** — CSS Modules with a simple dark theme. No component library, no animations beyond basic transitions. Prioritized correctness over polish.

## What I Intentionally Did Not Do

- **No ORM** — I used raw SQL via `sqlx` for clarity and control. ORMs add abstraction that can hide what's actually happening at the database level, which matters when you're thinking about money handling and data correctness.
- **No client-side caching** — the expense list always fetches fresh from the server. For a small dataset this is simpler and avoids cache-invalidation bugs.
- **No WebSocket / real-time updates** — unnecessary for a single-user expense tracker.
- **No input sanitization beyond trim** — PostgreSQL parameterized queries handle SQL injection; React handles XSS. Additional sanitization would be theater.