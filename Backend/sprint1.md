# Sprint 1 - StudentTaskHub

## User Stories

1. **As a student, I want to create a task with a deadline and priority so that others understand its urgency.**
   - Acceptance Criteria:
     - User can create a task with title, description, and deadline (YYYY-MM-DD)
     - User can optionally mark a task as "high" priority (radio button); otherwise it defaults to "normal"
     - Task is created with status "open" by default
     - Task stores who created it

2. **As a student, I want to view available (open) tasks so that I can find a task to work on.**
   - Acceptance Criteria:
     - User can view all tasks
     - User can filter tasks by status (e.g., only open tasks)
     - User can sort tasks by deadline or priority

3. **As a student, I want to claim an open task so that I can take responsibility for completing it.**
   - Acceptance Criteria:
     - User can claim a task that has status "open"
     - Claimed task status changes to "claimed" and records who claimed it
     - Already claimed tasks cannot be claimed again

4. **As a student, I want to see the status of a task so that I can track its progress.**
   - Acceptance Criteria:
     - Each task displays its current status (open, claimed, in_progress, done)
     - Task creator or claimer can update the status

5. **As a student, I want to edit or delete a task I created so that I can manage my tasks.**
   - Acceptance Criteria:
     - Only the task creator can edit task details (title, description, deadline, priority)
     - Only the task creator can delete a task
     - Other users receive a "forbidden" error if they try to edit/delete

## Planned Issues (Backend)

| Issue # | Title | Description | Status |
|---------|-------|-------------|--------|
| 1 | Project Setup | Initialize Go project with SQLite, set up project structure | ✅ Completed |
| 2 | Database Schema | Create users and tasks tables with proper constraints | ✅ Completed |
| 3 | User Registration & Login API | POST /api/register and POST /api/login endpoints | ✅ Completed |
| 4 | Create Task API | POST /api/tasks endpoint with validation | ✅ Completed |
| 5 | Get Tasks API (with filters) | GET /api/tasks with status filter and sort by deadline/priority | ✅ Completed |
| 6 | Get Single Task API | GET /api/tasks/{id} endpoint | ✅ Completed |
| 7 | Update Task API | PUT /api/tasks/{id} with creator-only authorization | ✅ Completed |
| 8 | Delete Task API | DELETE /api/tasks/{id} with creator-only authorization | ✅ Completed |
| 9 | Claim Task API | POST /api/tasks/{id}/claim endpoint | ✅ Completed |
| 10 | Update Task Status API | PATCH /api/tasks/{id}/status endpoint | ✅ Completed |
| 11 | CORS Middleware | Enable cross-origin requests for frontend integration | ✅ Completed |

## Successfully Completed

All 11 planned backend issues were completed. The backend API supports:
- User registration and login
- Full CRUD operations on tasks (Create, Read, Update, Delete)
- Task claiming functionality
- Task status updates
- Filtering tasks by status, creator, and claimer
- Sorting tasks by deadline or priority
- Authorization checks (only creators can edit/delete their tasks)
- CORS support for frontend integration

## Issues Not Completed

N/A - All planned Sprint 1 issues were completed.

## API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/register` | Register a new user |
| POST | `/api/login` | Login with username/password |
| POST | `/api/tasks` | Create a new task |
| GET | `/api/tasks` | List tasks (filter: `?status=open&sort=deadline`) |
| GET | `/api/tasks/{id}` | Get a single task |
| PUT | `/api/tasks/{id}?username=xxx` | Update a task (creator only) |
| DELETE | `/api/tasks/{id}?username=xxx` | Delete a task (creator only) |
| POST | `/api/tasks/{id}/claim` | Claim an open task |
| PATCH | `/api/tasks/{id}/status?username=xxx` | Update task status |