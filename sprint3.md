## User Stories

### From Sprint 1
1. As a student, I want to create a task with a deadline and priority so that others understand its urgency.
2. As a student, I want to view available (open) tasks so that I can find a task to work on.
3. As a student, I want to claim an open task so that I can take responsibility for completing it.
4. As a student, I want to see the status of a task so that I can track its progress.
5. As a student, I want to edit or delete a task I created so that I can manage my tasks.

### From Sprint 2
6. As a student, I want to register an account so that I can create and manage tasks securely.
7. As a student, I want to log in to my account so that I can access the system.
8. As a student, I want to search for tasks by keyword so that I can find specific tasks quickly.
9. As a student, I want to filter tasks by priority, status, deadline range, and creator so that I can find relevant tasks.

### New for Sprint 3
10. As a student, I want to create and update my profile so that others can see my experience and skills.
11. As a student, I want to view other students' profiles so that I can learn about potential collaborators.
12. As a student, I want to leave feedback on completed tasks so that I can rate the quality of work done.
13. As a student, I want to see feedback on tasks so that I can evaluate task quality before claiming.

---
### Backend - New Features
- **User Profiles**: Users can view and update their profile with full name, bio, major, year, and skills
- **Task Feedback/Reviews**: Users can rate completed tasks (1-5 stars) and leave comments
- **Auto-profile creation**: When a user registers, an empty profile is automatically created
- **Cascade delete**: Deleting a task also removes all associated feedback
- **New database tables**: Added `profiles` and `feedback` tables with foreign key constraints

### Backend - New Unit Tests
- 13 new unit tests for profile and feedback functionality
- Total: 47+ backend unit tests

### Integration
- Frontend-backend integration via REST API and CORS middleware
- Frontend calls backend API for all CRUD operations

---

## Backend Unit Tests

All tests are in `handlers/handler_test.go`. Run with:
```bash
go test ./handlers/ -v
```

### Register Tests (4 tests)
| Test Name | Description |
|-----------|-------------|
| TestRegister_Success | Registers a new user, verifies 201 and correct username |
| TestRegister_MissingFields | Missing data returns 400 |
| TestRegister_DuplicateUsername | Duplicate username returns 409 |
| TestRegister_ShortPassword | Password < 6 chars returns 400 |

### Login Tests (4 tests)
| Test Name | Description |
|-----------|-------------|
| TestLogin_Success | Correct credentials return 200 |
| TestLogin_WrongPassword | Wrong password returns 401 |
| TestLogin_NonExistentUser | Non-existent user returns 401 |
| TestLogin_MissingFields | Missing data returns 400 |

### CreateTask Tests (5 tests)
| Test Name | Description |
|-----------|-------------|
| TestCreateTask_Success | Creates task, verifies all fields |
| TestCreateTask_MissingFields | Missing data returns 400 |
| TestCreateTask_InvalidDeadline | Bad date format returns 400 |
| TestCreateTask_DefaultPriority | No priority defaults to "normal" |
| TestCreateTask_UnregisteredUser | Unregistered user returns 400 |

### GetTasks Tests (11 tests)
| Test Name | Description |
|-----------|-------------|
| TestGetTasks_Empty | Empty database returns empty array |
| TestGetTasks_FilterByStatus | Filters by status=open |
| TestGetTasks_FilterByCreatedUser | Filters by created_by |
| TestGetTasks_FilterByClaimedUser | Filters by claimed_by |
| TestGetTasks_Search | Searches by keyword "ML" |
| TestSearchTasks_ByKeyword | Searches description for "ML" |
| TestGetTasks_SearchCaseInsensitive | Case-insensitive search |
| TestGetTasks_FilterByPriority | Filters by priority=high |
| TestGetTasks_FilterByDeadlineRange | Filters by deadline_before |
| TestGetTasks_SortByDeadline | Sorts earliest first |

### GetTask Tests (2 tests)
| Test Name | Description |
|-----------|-------------|
| TestGetTask_Success | Gets task by ID |
| TestGetTask_NotFound | Non-existent ID returns 404 |

### UpdateTask Tests (2 tests)
| Test Name | Description |
|-----------|-------------|
| TestUpdateTask_Success | Creator updates title |
| TestUpdateTask_Forbidden | Non-creator gets 403 |

### DeleteTask Tests (2 tests)
| Test Name | Description |
|-----------|-------------|
| TestDeleteTask_Success | Creator deletes task |
| TestDeleteTask_Forbidden | Non-creator gets 403 |

### ClaimTask Tests (2 tests)
| Test Name | Description |
|-----------|-------------|
| TestClaimTask_Success | Claims open task |
| TestClaimTask_AlreadyClaimed | Already claimed returns 409 |

### UpdateTaskStatus Tests (4 tests)
| Test Name | Description |
|-----------|-------------|
| TestUpdateTaskStatus_Success | Updates to in_progress |
| TestUpdateTaskStatus_InvalidStatus | Invalid status returns 400 |
| TestUpdateTaskStatus_Forbidden | Unauthorized user gets 403 |
| TestMarkTaskAsCompleted | Marks task as done |

### Profile Tests (5 tests) — Sprint 3 NEW
| Test Name | Description |
|-----------|-------------|
| TestGetProfile_Success | Gets profile after registration |
| TestGetProfile_NotFound | Non-existent profile returns 404 |
| TestUpdateProfile_Success | Updates full_name, bio, major, year, skills |
| TestUpdateProfile_Forbidden | Editing another user's profile returns 403 |
| TestUpdateProfile_MissingUsername | Missing username param returns 400 |

### Feedback Tests (8 tests) — Sprint 3 NEW
| Test Name | Description |
|-----------|-------------|
| TestAddFeedback_Success | Adds 5-star feedback with comment |
| TestAddFeedback_InvalidRating | Rating > 5 returns 400 |
| TestAddFeedback_TaskNotDone | Feedback on non-completed task returns 400 |
| TestAddFeedback_DuplicateFeedback | Duplicate feedback returns 409 |
| TestAddFeedback_Forbidden | Non-creator/claimer gets 403 |
| TestGetFeedback_Success | Gets feedback list for a task |
| TestGetFeedback_Empty | No feedback returns empty array |
| TestGetFeedback_TaskNotFound | Non-existent task returns 404 |

---

## Frontend Unit Tests

## Frontend Unit Tests

All frontend unit tests are located in:

- `frontend/src/pages/LoginPage.test.jsx`
- `frontend/src/utils/validation.test.js`
- `frontend/src/pages/TaskDetailsModal.test.jsx`
- `frontend/src/pages/SearchFilterTasks.test.jsx`
- `frontend/src/pages/ForgotPasswordModal.test.jsx`
- `frontend/src/pages/ClaimTaskFlow.test.jsx`
- `frontend/src/pages/TaskCategoryRendering.test.jsx`

Frontend test setup/config files:

- `frontend/src/test/setupTests.js`
- `frontend/vite.config.js`

## Frontend Unit Test Summary

Frontend unit tests validate both Sprint 2 and Sprint 3 functionality.

- **Total frontend unit tests passed:** 38+
- **Total frontend test files passed:** 7

### Login Page Tests (2 tests)

| Test Name | Description |
|-----------|-------------|
| shows validation errors for invalid email and UFID | Verifies invalid email and UFID show validation messages |
| logs in a valid user and navigates to /home | Verifies valid login redirects to the home page |

### Validation Helper Tests (6 tests)

| Test Name | Description |
|-----------|-------------|
| accepts a valid @ufl.edu email | Confirms valid UFL email passes validation |
| rejects a non-UFL email | Confirms non-UFL email fails validation |
| accepts a valid 8 digit UFID | Confirms valid UFID passes validation |
| rejects an invalid UFID | Confirms invalid UFID fails validation |
| accepts non-empty trimmed text | Confirms meaningful text is accepted |
| rejects empty trimmed text | Confirms empty text is rejected |

### Task Details Modal Tests (10 tests)

| Test Name | Description |
|-----------|-------------|
| opens task details modal on task click | Verifies task click opens detail modal |
| shows creator information correctly | Verifies creator details are displayed correctly |
| shows task information correctly | Verifies task details are displayed correctly |
| shows category badge correctly | Verifies Academic or Personal badge is displayed |
| shows correct task type for Academic and Personal | Verifies correct category rendering |
| closes modal on close button click | Verifies modal closes properly |
| renders Claim Task button | Verifies Claim Task button is visible |
| claim button updates claimed state | Verifies Claim Task updates task claim state |
| handles missing creator data gracefully | Verifies fallback values are shown |
| ensures modal stays within viewport | Verifies proper scrolling and layout |

### Search and Filter Tests (8 tests)

| Test Name | Description |
|-----------|-------------|
| searches tasks by title keyword | Verifies search by title works |
| searches tasks by description keyword | Verifies search by description works |
| search is case insensitive | Verifies search works regardless of case |
| filters tasks by priority | Verifies priority filter works |
| filters tasks by category | Verifies category filter works |
| filters tasks by status | Verifies status filter works |
| combined filters work correctly | Verifies multiple filters work together |
| clears all filters correctly | Verifies clear filters restores all tasks |

### Forgot Password Modal Tests (6 tests)

| Test Name | Description |
|-----------|-------------|
| opens forgot password modal from login page | Verifies Forgot Password opens modal |
| accepts valid email input | Verifies email input works correctly |
| rejects invalid email format | Verifies invalid email shows error |
| prevents empty submission | Verifies empty input is not allowed |
| shows success confirmation after submit | Verifies success message is shown |
| closes modal correctly | Verifies modal closes properly |

### Claim Task Flow Tests (3 tests)

| Test Name | Description |
|-----------|-------------|
| disables claim button for already claimed task | Verifies Claim Task is disabled or hidden for claimed tasks |
| shows success message after claiming task | Verifies success message appears after claiming |
| updates task status after claim | Verifies claimed task status is updated in UI |

### Task Category Rendering Tests (3 tests)

| Test Name | Description |
|-----------|-------------|
| renders Academic task badge correctly | Verifies Academic badge is displayed properly |
| renders Personal task badge correctly | Verifies Personal badge is displayed properly |
| opens correct modal data for selected category task | Verifies correct task data is loaded based on clicked category task |

---

## Backend API Documentation

### Base URL
```
http://localhost:8080/api
```

### All Endpoints

| Method | Endpoint | Description | Sprint |
|--------|----------|-------------|--------|
| POST | /api/register | Register new user | 2 |
| POST | /api/login | Login user | 2 |
| GET | /api/profile/{username} | Get user profile | 3 |
| PUT | /api/profile/{username}?username=xxx | Update own profile | 3 |
| POST | /api/tasks | Create task | 1 |
| GET | /api/tasks | List/search/filter tasks | 1+2 |
| GET | /api/tasks/{id} | Get single task | 1 |
| PUT | /api/tasks/{id}?username=xxx | Edit task (creator only) | 1 |
| DELETE | /api/tasks/{id}?username=xxx | Delete task (creator only) | 1 |
| POST | /api/tasks/{id}/claim | Claim open task | 1 |
| PATCH | /api/tasks/{id}/status?username=xxx | Update task status | 1 |
| POST | /api/tasks/{id}/feedback?username=xxx | Add feedback to completed task | 3 |
| GET | /api/tasks/{id}/feedback | Get feedback for a task | 3 |

### Sprint 3 New Endpoints

#### GET /api/profile/{username}
Get a user's profile.

**Example:** `GET /api/profile/Chinmai`

**Success Response (200):**
```json
{
    "id": 1,
    "username": "Chinmai",
    "full_name": "Chinmai Reddy",
    "bio": "CS student at UF",
    "major": "Computer Science",
    "year": "Senior",
    "skills": "Go, Python, React",
    "created_at": "2026-03-01T12:00:00Z",
    "updated_at": "2026-03-01T12:00:00Z"
}
```

**Error:** 404 if profile not found

---

#### PUT /api/profile/{username}?username=xxx
Update your own profile. Only the profile owner can edit.

**Example:** `PUT /api/profile/Chinmai?username=Chinmai`

**Request Body:**
```json
{
    "full_name": "Chinmai Reddy",
    "bio": "CS student at UF",
    "major": "Computer Science",
    "year": "Senior",
    "skills": "Go, Python, React"
}
```

**Error Responses:**
- 400: Missing username query parameter
- 403: Trying to edit another user's profile

---

#### POST /api/tasks/{id}/feedback?username=xxx
Add feedback to a completed task. Only the task creator or claimer can leave feedback. Rating must be 1-5.

**Example:** `POST /api/tasks/1/feedback?username=Alice`

**Request Body:**
```json
{
    "rating": 5,
    "comment": "Great work on this task!"
}
```

**Success Response (201):**
```json
{
    "id": 1,
    "task_id": 1,
    "username": "Alice",
    "rating": 5,
    "comment": "Great work on this task!",
    "created_at": "2026-03-01T12:00:00Z"
}
```

**Error Responses:**
- 400: Task not completed (must be status "done")
- 400: Rating not between 1-5
- 403: User is not creator or claimer
- 409: User already submitted feedback for this task

---

#### GET /api/tasks/{id}/feedback
Get all feedback for a task.

**Example:** `GET /api/tasks/1/feedback`

**Success Response (200):**
```json
[
    {
        "id": 1,
        "task_id": 1,
        "username": "Alice",
        "rating": 5,
        "comment": "Great work!",
        "created_at": "2026-03-01T12:00:00Z"
    }
]
```

### GET /api/tasks Query Parameters (unchanged from Sprint 2)
| Parameter | Example | Description |
|-----------|---------|-------------|
| status | ?status=open | Filter by status |
| priority | ?priority=high | Filter by priority |
| created_by | ?created_by=Chinmai | Filter by creator |
| claimed_by | ?claimed_by=Alice | Filter by claimer |
| search | ?search=ML | Search title and description |
| deadline_before | ?deadline_before=2026-04-01 | Tasks due before date |
| deadline_after | ?deadline_after=2026-03-01 | Tasks due after date |
| sort | ?sort=deadline | Sort: deadline/priority/newest/oldest |

