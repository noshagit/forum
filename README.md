# Forum

A simple web-based forum application with user authentication, posting, commenting, and liking features. The project is split into a Go backend and a static frontend.

## Features
- User registration and authentication
- Create, view, and delete posts
- Comment on posts
- Like posts
- User profile management

## Project Structure
```
forum/
├── back/           # Go backend (API, handlers, database)
│   ├── api/        # API endpoints (JS)
│   ├── database/   # SQLite DB and setup scripts
│   ├── handlers/   # Go HTTP handlers
│   └── main/       # Main Go entry point
├── front/          # Frontend (HTML, CSS, JS)
│   ├── comments/   # Comments UI
│   ├── images/     # Static images
│   ├── login/      # Login UI
│   ├── password/   # Password reset UI
│   ├── post-list/  # Post list UI
│   ├── pp/         # Profile pictures
│   ├── profil/     # Profile UI
│   └── register/   # Registration UI
└── README.md       # Project documentation
```

## Getting Started

### Prerequisites
- Go (1.23.0)
- Node.js (for frontend development, optional)
- SQLite3

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd back
   ```
2. Install Go dependencies:
   ```bash
   go mod tidy
   ```
3. Run the backend server:
   ```bash
   cd main
   go run main.go
   ```

## Usage
- Register a new account or log in.
- Create, view, and interact with posts.
- Manage your profile and password.

## License
MIT
---
Feel free to contribute or open issues for improvements!

