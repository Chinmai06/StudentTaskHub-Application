# sprint 2

## User Stories 
These are the user stories for sprint2

1. As a student, I want to log in with my email and password so that I can access my task manager and view my tasks.
2. As a student, I want to create a new account if I am a new user so that I can register and start using the application.
3. As a student, I want to use an eye icon to show or hide my password while entering it so that I can type it correctly without mistakes.
4. As a student, I want to see all the tasks I created so that I can track and manage them.
5. As a student, I want to view the tasks I have claimed so that I know which tasks I need to complete.
6. As a student, I want to mark a task as completed after finishing it so that the system updates its status.
7. As a student, I want to search tasks by keywords so that I can quickly find the tasks I need.
8. As a student, I want to update my profile details successfully so that my personal information stays accurate and up to date.
9. As a student, I want a back button on the Tasks page so that I can easily navigate to the previous page without confusion.
10. As a student, I want the logout button to be clearly visible and properly placed so that I can securely log out of my account when needed.


| Issue # | Title                                                     | Description                                                                           | Status    |
| ------- | --------------------------------------------------------- | ------------------------------------------------------------------------------------- | --------- |
| 27      | Logout button not visible / improperly placed             | Fix UI issue where logout button is missing or incorrectly positioned                 | Completed |
| 26      | Back button missing on Tasks page                         | Add navigation back button on Tasks page for better usability                         | Completed |
| 25      | Profile details update functionality not working properly | Debug and fix profile update API/UI integration                                       | Completed |
| 24      | Password Visibility Toggle                                | Implement toggle to show/hide password in input fields                                | Completed |
| 23      | Account Creation                                          | Implement user registration functionality                                             | Completed |
| 22      | Sign In Functionality                                     | Implement login authentication for users                                              | Completed |
| 21      | Implement search tasks API                                | Create API to search tasks by title/description keywords                              | Completed |
| 20      | Mark task as completed                                    | Update task status to "done" for claimed tasks                                        | Completed |
| 19      | Fetch tasks claimed by user                               | Retrieve tasks assigned to a specific user                                            | Completed |
| 18      | Fetch tasks created by user                               | Retrieve tasks created by a specific user                                             | Completed |
| 17      | Implement API to fetch tasks created by a user            | Backend endpoint for user-specific task retrieval                                     | Completed |
| 16      | Unit Testing                                              | Write and validate unit tests for all backend APIs (handlers, edge cases, validation) | Completed |
| 15      | Create users table in database                            | Design and implement users table schema                                               | Completed |
| 14      | Implement basic user authentication                       | Add login/register with password validation                                           | Completed |

<<<<<<< HEAD
In Backend Prespective all the 6 issues are completed

In Frontend Prespective all the 5 issues are completed
=======
In Backend Prespective all the 8 issues are completed

In Frontend Prespective all the 6 issues are completed
## Frontend Unit Tests

To be added by frontend team

## Cypress Tests

To be added by frontend team

## Backend Unit Tests

All tests are in `handlers/handler_test.go`. Run with:

```bash
go test ./handlers/ -v
```

### Register Tests (4 tests)

| Test Name | Description |
|-----------|-------------|
| TestRegister_Success | Registers a new user and verifies 201 response with correct username |
| TestRegister_MissingFields | Sends incomplete data, expects 400 error |
| TestRegister_DuplicateUsername | Registers same username twice, expects 409 conflict |
| TestRegister_ShortPassword | Sends password under 6 characters, expects 400 error |

### Login Tests (4 tests)

| Test Name | Description |
|-----------|-------------|
| TestLogin_Success | Logs in with correct credentials, verifies 200 response |
| TestLogin_WrongPassword | Uses wrong password, expects 401 unauthorized |
| TestLogin_NonExistentUser | Tries to login as non-existent user, expects 401 |
| TestLogin_MissingFields | Sends incomplete data, expects 400 error |

### CreateTask Tests (5 tests)

| Test Name | Description |
|-----------|-------------|
| TestCreateTask_Success | Creates task with all fields, verifies title/priority/status |
| TestCreateTask_MissingFields | Sends incomplete data, expects 400 error |
| TestCreateTask_InvalidDeadline | Sends bad date format, expects 400 error |
| TestCreateTask_DefaultPriority | Creates task without priority, verifies default is "normal" |
| TestCreateTask_UnregisteredUser | Non-existent user tries to create task, expects 400 error |

### GetTasks Tests (11 tests)

| Test Name | Description |
|-----------|-------------|
| TestGetTasks_Empty | Gets tasks from empty database, expects empty array |
| TestGetTasks_FilterByStatus | Filters tasks by status=open, verifies correct count |
| TestGetTasks_FilterByCreatedUser | Filters tasks by created_by=nikita, verifies only nikita's tasks returned |
| TestGetTasks_FilterByClaimedUser | Claims a task as alice, filters by claimed_by=alice, verifies correct result |
| TestGetTasks_Search | Searches tasks by keyword "ML", verifies only matching task returned |
| TestSearchTasks_ByKeyword | Creates 3 tasks, searches "ML", verifies match in description |
| TestGetTasks_SearchCaseInsensitive | Searches with lowercase "ml", verifies case-insensitive match |
| TestGetTasks_FilterByPriority | Filters by priority=high, verifies only high priority returned |
| TestGetTasks_FilterByDeadlineRange | Filters tasks before a date, verifies correct filtering |
| TestGetTasks_SortByDeadline | Sorts tasks by deadline, verifies earliest first |

### GetTask Tests (2 tests)

| Test Name | Description |
|-----------|-------------|
| TestGetTask_Success | Gets task by ID, verifies correct task returned |
| TestGetTask_NotFound | Gets non-existent task ID, expects 404 error |

### UpdateTask Tests (2 tests)

| Test Name | Description |
|-----------|-------------|
| TestUpdateTask_Success | Creator updates task title, verifies change |
| TestUpdateTask_Forbidden | Non-creator tries to edit, expects 403 forbidden |

### DeleteTask Tests (2 tests)

| Test Name | Description |
|-----------|-------------|
| TestDeleteTask_Success | Creator deletes task, verifies 200 response |
| TestDeleteTask_Forbidden | Non-creator tries to delete, expects 403 forbidden |

### ClaimTask Tests (2 tests)

| Test Name | Description |
|-----------|-------------|
| TestClaimTask_Success | User claims open task, verifies status changes to "claimed" |
| TestClaimTask_AlreadyClaimed | Second user tries to claim already claimed task, expects 409 conflict |

### UpdateTaskStatus Tests (4 tests)

| Test Name | Description |
|-----------|-------------|
| TestUpdateTaskStatus_Success | Claimer updates status to in_progress, verifies change |
| TestUpdateTaskStatus_InvalidStatus | Sends invalid status value, expects 400 error |
| TestUpdateTaskStatus_Forbidden | Unauthorized user tries to update status, expects 403 forbidden |
| TestMarkTaskAsCompleted | Claims task then marks as done, verifies status is "done" |

---
>>>>>>> b377c5f912481f10342a0f5832da8901bfa59bef

## Backend API Documentation

### URL

```url
http://localhost:8080/api
```

### How to Run

```bash
# Delete old database (schema changed from Sprint 1)
del studenttaskhub.db

# Install dependencies
go mod tidy

# Start the server
go run main.go

# Run unit tests (separate terminal)
go test ./handlers/ -v
```

---

### User Endpoints

#### POST /api/register

Register a new user account.

**Request Body:**

```json
{
    "username": "",
    "email": "chinmai@ufl.edu",
    "password": "pass123"
}
```

**Success Response (201):**

```json
{
    "message": "User registered successfully",
    "username": "Chinmai"
}
```

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing fields (username, email, or password) |
| 400 | Password shorter than 6 characters |
| 409 | Username or email already exists |

---

#### POST /api/login

Authenticate a user.

**Request Body:**

```json
{
    "username": "Chinmai",
    "password": "pass123"
}
```

**Success Response (200):**

```json
{
    "message": "Login successful",
    "username": "Chinmai"
}
```

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing username or password |
| 401 | Invalid username or password |

---

### Task Endpoints

#### POST /api/tasks

Create a new task. The `created_by` user must be registered.

**Request Body:**

```json
{
    "title": "ML Assignment 3",
    "description": "Neural network implementation",
    "deadline": "2026-03-01",
    "priority": "high",
    "created_by": "Chinmai"
}
```

**Success Response (201):**

```json
{
    "id": 1,
    "title": "ML Assignment 3",
    "description": "Neural network implementation",
    "deadline": "2026-03-01",
    "priority": "high",
    "status": "open",
    "created_by": "Chinmai",
    "created_at": "2026-02-13T12:00:00Z",
    "updated_at": "2026-02-13T12:00:00Z"
}
```

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing title, deadline, or created_by |
| 400 | Invalid deadline format (must be YYYY-MM-DD) |
| 400 | User does not exist (not registered) |

---

#### GET /api/tasks

`http://localhost:8080/api/tasks`

Retrieve all tasks with optional filtering, searching, and sorting.

**Query Parameters:**

| Parameter | Description | Example |
|-----------|--------------|---------|
| status | Filter by task status | ?status=open |
| priority |  Filter by priority | ?priority=high |
| created_by | Filter by creator | ?created_by=Chinmai |
| claimed_by | Filter by claimer | ?claimed_by=Alice |
| search | Search title and description (case-insensitive) | ?search=ML |
| deadline_before | Tasks due on or before date (YYYY-MM-DD) | ?deadline_before=2026-04-01 |
| deadline_after |  Tasks due on or after date (YYYY-MM-DD) | ?deadline_after=2026-03-01 |
| sort | Sort order: deadline, priority, newest, oldest | ?sort=deadline |

All parameters can be combined: `?status=open&priority=high&sort=deadline&search=ML`

**Success Response (200):**

```json
[
    {
        "id": 1,
        "title": "ML Assignment 3",
        "description": "Neural network implementation",
        "deadline": "2026-03-01",
        "priority": "high",
        "status": "open",
        "created_by": "Chinmai",
        "created_at": "2026-02-13T12:00:00Z",
        "updated_at": "2026-02-13T12:00:00Z"
    }
]
```

---

#### GET /api/tasks/{id}

Retrieve a single task by its ID.

**Example:** `http://localhost:8080/api/tasks/1`

**Success Response (200):** Single task object

```json
{
    "id": 2,
    "title": "Fix login bug",
    "description": "Frontend issue",
    "deadline": "2025-12-31",
    "priority": "high",
    "status": "open",
    "created_by": "john",
    "created_at": "2026-03-22T03:37:20.738224-05:00",
    "updated_at": "2026-03-22T03:37:20.738224-05:00"
}
```

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Invalid task ID |
| 404 | Task not found |

---

#### PUT /api/tasks/{id}?username=xxx

Update a task. Only the task creator can edit.

**Example:** `http://localhost:8080/api/tasks/2?username=john`

**Request Body (all fields optional unchanged fields keep their values):**

```json
{
    "title": "Updated Title",
    "description": "Updated description",
    "deadline": "2026-04-01",
    "priority": "normal"
}
```

**Success Response (200):** Updated task object

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing username query parameter |
| 400 | Invalid priority or deadline format |
| 403 | User is not the task creator |
| 404 | Task not found |

---

#### DELETE /api/tasks/{id}?username=xxx

Delete a task. Only the task creator can delete.

**Example:** `http://localhost:8080/api/tasks/1?username=Chinmai`

**Success Response (200):**

```json
{
    "message": "Task deleted successfully"
}
```

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing username query parameter |
| 403 | User is not the task creator |
| 404 | Task not found |

---

#### POST /api/tasks/{id}/claim

Claim an open task. Only tasks with status "open" can be claimed. The `claimed_by` user must be registered.

**Example:** `http://localhost:8080/api/tasks/1/claim`

**Request Body:**

```json
{
    "claimed_by": "jhon"
}
```

**Success Response (200):** Updated task object with status "claimed"

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing claimed_by |
| 400 | User does not exist (not registered) |
| 404 | Task not found |
| 409 | Task is not open for claiming (already claimed) |

---

#### PATCH /api/tasks/{id}/status?username=xxx

Update the status of a task. Only the task creator or claimer can update status.

**Example:** `http://localhost:8080/api/tasks/1/status?username=Alice`

**Request Body:**

```json
{
    "status": "in_progress"
}
```

**Valid Status Values:**

| Status | Description |
|--------|-------------|
| open | Task is available for claiming |
| claimed | Task has been claimed by a user |
| in_progress | Task is being worked on |
| done | Task is completed |

**Success Response (200):** Updated task object

**Error Responses:**

| Status | Reason |
|--------|--------|
| 400 | Missing username query parameter |
| 400 | Invalid status value |
| 403 | User is not the creator or claimer |
| 404 | Task not found |
