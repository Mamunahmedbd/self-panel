# PPPoE Client Login System Architecture

## System Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         USER INTERACTION                             │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      LOGIN PAGE (login.templ)                        │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  [Username/Email Input]  ← form.UsernameOrEmail               │  │
│  │  [Password Input]        ← form.Password                      │  │
│  │  [Sign In Button]                                             │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   LOGIN HANDLER (login.go)                           │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  1. Bind Form Data                                            │  │
│  │  2. Validate Input (trim, check empty)                        │  │
│  │  3. Query ClientUser Table                                    │  │
│  │     WHERE username = ? OR email = ?                           │  │
│  │  4. Check Client Found                                        │  │
│  │  5. Verify Status = 'active'                                  │  │
│  │  6. Compare Password (plain-text)                             │  │
│  │  7. Get/Create User Record                                    │  │
│  │  8. Get/Create Profile Record                                 │  │
│  │  9. Store client_id in Session                                │  │
│  │ 10. Login User (Auth.Login)                                   │  │
│  │ 11. Redirect to Profile                                       │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    ▼                           ▼
┌──────────────────────────────┐  ┌──────────────────────────────┐
│   DATABASE: clients          │  │   DATABASE: users            │
│  ┌────────────────────────┐  │  │  ┌────────────────────────┐  │
│  │ id                     │  │  │  │ id                     │  │
│  │ username (unique)      │  │  │  │ name                   │  │
│  │ password (plain-text)  │  │  │  │ email (unique)         │  │
│  │ email                  │  │  │  │ password (dummy)       │  │
│  │ status (active/...)    │  │  │  └────────────────────────┘  │
│  │ name                   │  │  │                              │
│  │ balance                │  │  │   DATABASE: profiles         │
│  │ package_pool           │  │  │  ┌────────────────────────┐  │
│  │ ...                    │  │  │  │ user_id (FK)           │  │
│  └────────────────────────┘  │  │  │ bio                    │  │
└──────────────────────────────┘  │  │ fully_onboarded        │  │
                                  │  │ ...                    │  │
                                  │  └────────────────────────┘  │
                                  └──────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        SESSION STORAGE                               │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  authSessionKeyUserID: <user_id>                              │  │
│  │  authSessionKeyAuthenticated: true                            │  │
│  │  client_id: <client_id>                                       │  │
│  │  client_username: <username>                                  │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    AUTHENTICATED STATE                               │
│  ┌───────────────────────────────────────────────────────────────┐  │
│  │  page.IsAuth = true                                           │  │
│  │  page.AuthUser = <User object>                                │  │
│  │  Container.GetAuthenticatedClient() → <ClientUser object>     │  │
│  └───────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    ▼                           ▼
┌──────────────────────────────┐  ┌──────────────────────────────┐
│   NAVBAR (navbar.templ)      │  │   PRIVATE ROUTES             │
│  ┌────────────────────────┐  │  │  ┌────────────────────────┐  │
│  │ if page.IsAuth:        │  │  │  │ /auth/profile          │  │
│  │   Show User Menu       │  │  │  │ /auth/dashboard        │  │
│  │   - Profile Picture    │  │  │  │ /auth/preferences      │  │
│  │   - Name & Email       │  │  │  │ /auth/logout           │  │
│  │   - Dashboard Link     │  │  │  │ ...                    │  │
│  │   - Settings Link      │  │  │  │                        │  │
│  │   - Logout Link        │  │  │  │ Protected by:          │  │
│  │ else:                  │  │  │  │ RequireAuthentication()│  │
│  │   Show Login Button    │  │  │  └────────────────────────┘  │
│  └────────────────────────┘  │  └──────────────────────────────┘
└──────────────────────────────┘
```

## Data Flow

### Login Request Flow

```
Browser → POST /user/login
         ↓
    Form Binding
         ↓
    Input Validation
         ↓
    Database Query (clients table)
         ↓
    Status Check (active?)
         ↓
    Password Check (plain-text)
         ↓
    User/Profile Creation (if needed)
         ↓
    Session Creation
         ↓
    Auth Login
         ↓
    Redirect → /auth/profile
         ↓
    Navbar Updates (shows user menu)
```

### Authenticated Request Flow

```
Browser → GET /auth/dashboard
         ↓
    Middleware: RequireAuthentication()
         ↓
    Check Session (authSessionKeyAuthenticated)
         ↓
    Load User (from authSessionKeyUserID)
         ↓
    Set page.IsAuth = true
         ↓
    Set page.AuthUser = <User>
         ↓
    Handler: GetAuthenticatedClient()
         ↓
    Load ClientUser (from session client_id)
         ↓
    Render Page with Client Data
         ↓
    Browser (displays dashboard)
```

## Component Relationships

```
┌─────────────────────────────────────────────────────────────────┐
│                        FRONTEND LAYER                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ login.templ  │  │ navbar.templ │  │ pages/*.templ│          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                        ROUTING LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  router.go   │  │  login.go    │  │ other routes │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                       MIDDLEWARE LAYER                           │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  RequireAuthentication() → LoadAuthenticatedUser()       │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                        SERVICE LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Container   │  │  AuthClient  │  │ Helpers      │          │
│  │  - ORM       │  │  - Login()   │  │ - GetClient()│          │
│  │  - Auth      │  │  - Logout()  │  │ - IsAuth()   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                          ORM LAYER (Ent)                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ ClientUser   │  │    User      │  │   Profile    │          │
│  │  Entity      │  │   Entity     │  │   Entity     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      DATABASE LAYER (MariaDB)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   clients    │  │    users     │  │   profiles   │          │
│  │   table      │  │    table     │  │    table     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

## Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      SECURITY LAYERS                             │
│                                                                  │
│  Layer 1: Input Validation                                      │
│  ┌────────────────────────────────────────────────────────┐     │
│  │ - Trim whitespace                                      │     │
│  │ - Check empty fields                                   │     │
│  │ - Sanitize input                                       │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                  │
│  Layer 2: Authentication                                        │
│  ┌────────────────────────────────────────────────────────┐     │
│  │ - Username/Email lookup                                │     │
│  │ - Plain-text password comparison (PPPoE requirement)   │     │
│  │ - Status verification (active only)                    │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                  │
│  Layer 3: Session Management                                    │
│  ┌────────────────────────────────────────────────────────┐     │
│  │ - Secure cookie storage                                │     │
│  │ - Session encryption                                   │     │
│  │ - CSRF protection                                      │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                  │
│  Layer 4: Authorization                                         │
│  ┌────────────────────────────────────────────────────────┐     │
│  │ - Middleware checks                                    │     │
│  │ - Route protection                                     │     │
│  │ - Status verification                                  │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                  │
│  Layer 5: Logging & Monitoring                                  │
│  ┌────────────────────────────────────────────────────────┐     │
│  │ - Login attempts                                       │     │
│  │ - Failed authentications                               │     │
│  │ - Client activities                                    │     │
│  └────────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────────┘
```

## File Structure

```
self_panel/
├── ent/
│   └── schema/
│       └── client.go              ← ClientUser entity definition
│
├── pkg/
│   ├── context/
│   │   └── client.go              ← Session key constants
│   │
│   ├── routing/
│   │   ├── routes/
│   │   │   ├── router.go          ← Route registration
│   │   │   └── login.go           ← Login handler logic
│   │   └── routenames/
│   │       └── routenames.go      ← Route name constants
│   │
│   ├── services/
│   │   ├── container.go           ← Service container
│   │   ├── auth.go                ← Auth client
│   │   └── client_helpers.go     ← Client helper functions
│   │
│   └── types/
│       └── login.go               ← Login form structure
│
├── templates/
│   ├── pages/
│   │   └── login.templ            ← Login page template
│   └── components/
│       └── navbar.templ           ← Navbar component
│
└── docs/
    ├── CLIENT_LOGIN_IMPLEMENTATION.md  ← Full documentation
    ├── QUICK_REFERENCE.md              ← Quick reference guide
    └── ARCHITECTURE.md                 ← This file
```

---

**Architecture Version**: 1.0.0
**Last Updated**: 2025-12-22
