# TaskFlow API Documentation

Base URL: `http://localhost:8080`

## Authentication

Most endpoints require authentication via JWT Bearer Token. Include the token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

---

## Public Endpoints

### Health Check

Check if the API is running.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok",
  "message": "TaskFlow API is running"
}
```

---

### Register

Register a new user account.

**Endpoint:** `POST /api/auth/register`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Request Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | User's name (min: 2 characters) |
| email | string | Yes | Valid email address |
| password | string | Yes | Password (min: 8 characters) |

**Success Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "invalid input"
}
```

---

### Login

Authenticate and receive a JWT token.

**Endpoint:** `POST /api/auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Request Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | Valid email address |
| password | string | Yes | User password |

**Success Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "invalid credentials"
}
```

---

## Protected Endpoints

All endpoints in this section require authentication via JWT Bearer Token.

---

### Get Current User

Get details of the authenticated user.

**Endpoint:** `GET /api/me`

**Response (200 OK):**
```json
{
  "user_id": 1,
  "email": "john@example.com"
}
```

---

### Projects

#### Create Project

Create a new project.

**Endpoint:** `POST /api/projects`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "My Project",
  "description": "Project description"
}
```

**Request Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Project name (min: 2 characters) |
| description | string | No | Project description |

**Success Response (201 Created):**
```json
{
  "id": 1,
  "name": "My Project",
  "description": "Project description",
  "owner_id": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input
- `500 Internal Server Error` - Server error

---

#### Get My Projects

Get all projects owned by the authenticated user.

**Endpoint:** `GET /api/projects`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "My Project",
    "description": "Project description",
    "owner_id": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

---

#### Get Project By ID

Get a specific project with its members and tasks.

**Endpoint:** `GET /api/projects/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | int64 | Project ID |

**Success Response (200 OK):**
```json
{
  "id": 1,
  "name": "My Project",
  "description": "Project description",
  "owner_id": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "members": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "tasks": [
    {
      "id": 1,
      "project_id": 1,
      "title": "Task Title",
      "description": "Task description",
      "status": "todo",
      "assignee_id": null,
      "due_date": null,
      "created_by": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request` - Invalid project ID
- `403 Forbidden` - Access denied
- `404 Not Found` - Project not found
- `500 Internal Server Error` - Server error

---

#### Add Project Member

Add a member to a project.

**Endpoint:** `POST /api/projects/:id/members`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | int64 | Project ID |

**Request Body:**
```json
{
  "user_id": 2
}
```

**Request Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| user_id | int64 | Yes | ID of user to add |

**Success Response (200 OK):**
```json
{
  "message": "member added successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid project ID or input
- `500 Internal Server Error` - Server error

---

### Tasks

#### Create Task

Create a new task in a project.

**Endpoint:** `POST /api/projects/:id/tasks`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | int64 | Project ID |

**Request Body:**
```json
{
  "title": "Task Title",
  "description": "Task description",
  "status": "todo",
  "assignee_id": 2,
  "due_date": "2024-12-31T23:59:59Z"
}
```

**Request Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | Yes | Task title (min: 2 characters) |
| description | string | No | Task description |
| status | string | No | Task status: `todo`, `in_progress`, `done` |
| assignee_id | int64 | No | ID of user to assign task to |
| due_date | string (ISO 8601) | No | Task due date |

**Success Response (201 Created):**
```json
{
  "task": {
    "id": 1,
    "project_id": 1,
    "title": "Task Title",
    "description": "Task description",
    "status": "todo",
    "assignee_id": 2,
    "due_date": "2024-12-31T23:59:59Z",
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid project ID or input
- `500 Internal Server Error` - Server error

---

#### Update Task

Update an existing task.

**Endpoint:** `PATCH /api/tasks/:task_id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| task_id | int64 | Task ID |

**Request Body (all fields are optional):**
```json
{
  "title": "Updated Task Title",
  "description": "Updated description",
  "status": "in_progress",
  "assignee_id": 3,
  "due_date": "2024-12-31T23:59:59Z"
}
```

**Request Schema:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| title | string | No | Updated task title |
| description | string | No | Updated description |
| status | string | No | Updated status: `todo`, `in_progress`, `done` |
| assignee_id | int64 | No | Updated assignee ID |
| due_date | string (ISO 8601) | No | Updated due date |

**Success Response (200 OK):**
```json
{
  "task": {
    "id": 1,
    "project_id": 1,
    "title": "Updated Task Title",
    "description": "Updated description",
    "status": "in_progress",
    "assignee_id": 3,
    "due_date": "2024-12-31T23:59:59Z",
    "created_by": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid task ID or input
- `500 Internal Server Error` - Server error

---

#### Delete Task

Delete a task.

**Endpoint:** `DELETE /api/tasks/:task_id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| task_id | int64 | Task ID |

**Success Response (200 OK):**
```json
{
  "message": "task deleted successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid task ID
- `500 Internal Server Error` - Server error

---

## Task Status Values

The `status` field accepts the following values:
- `todo` - Task is pending
- `in_progress` - Task is being worked on
- `done` - Task is completed
```
