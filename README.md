# NexusBoard

![CI Badge](https://github.com/ArcherStand/nexus-board/actions/workflows/ci.yml/badge.svg)

A real-time collaborative application built with **Go**, **React**, and **WebSockets**.  
This project demonstrates a full-stack, microservice-based architecture designed for concurrency and modern web standards.

---

## 🚀 Key Features

- **Secure User Authentication**  
  Users can register and log in via a dedicated Go microservice, with authentication managed by JWTs.

- **Real-Time Messaging**  
  A live message board allows all connected users to communicate instantly without page reloads.

- **WebSocket Communication**  
  Leverages WebSockets for persistent, bi-directional communication between the client and server.

- **Microservice Architecture**  
  Backend is split into two distinct services (authentication and board logic) for separation of concerns and scalability.

- **CI/CD Pipeline**  
  Includes a GitHub Actions workflow to automatically build and test the backend and frontend on every push.

---

## 🛠 Tech Stack & Architecture

This project uses a modern tech stack to handle real-time events and user management efficiently.

### **Backend (Golang)**

- **Services**
  - **auth-service**: Handles user registration and login. Uses bcrypt for password hashing and issues JWTs.
  - **board-service**: Manages live WebSocket connections, authenticates clients using JWTs, and broadcasts messages.

- **Frameworks & Libraries**
  - [Gin](https://gin-gonic.com/) – HTTP routing and middleware
  - [GORM](https://gorm.io/) & [glebarez/sqlite](https://github.com/glebarez/sqlite) – Database interaction (pure-Go SQLite driver)
  - [Gorilla WebSocket](https://github.com/gorilla/websocket) – WebSocket handling
  - [gin-contrib/cors](https://github.com/gin-contrib/cors) – CORS management

---

### **Frontend (React + TypeScript)**

- **Framework**: React 18+ with [Vite](https://vitejs.dev/) as the build tool
- **Language**: TypeScript
- **State Management**: React Hooks (`useState`, `useEffect`)
- **API Communication**
  - [axios](https://axios-http.com/) – RESTful API calls to the auth-service
  - [react-use-websocket](https://github.com/robtaussig/react-use-websocket) – Custom hook for managing WebSocket connections

---

## 🖥 Getting Started

You will need **Go** and **Node.js** installed.

### 1️⃣ Clone the Repository
```bash
git clone https://github.com/ArcherStand/nexus-board.git
cd nexus-board
```

### 2️⃣ Run the Backend Services

Open **two separate terminals** for the backend.

**Terminal 1 – Start the auth-service**
```bash
cd backend/auth-service

# Set JWT secret key (Windows CMD example)
set JWT_SECRET_KEY=my_development_secret_key

# Run the service
go run .
```

**Terminal 2 – Start the board-service**
```bash
cd backend/board-service
go run .
```

### 3️⃣ Run the Frontend

Open a **third terminal**.

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be available at: **[http://localhost:5173](http://localhost:5173)**

---

## 📌 Future Improvements

- Implement different chat rooms via `/ws/board/:boardId`
- Store message history in the database
- Add "task board" functionality with draggable cards
- Deploy services to a cloud provider using Docker

---
