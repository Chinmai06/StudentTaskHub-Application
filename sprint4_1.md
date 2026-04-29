# Sprint 4 - StudentTaskHub

## Work Completed in Sprint 4

### Backend - New Features
- **In-App Notifications**: Users receive notifications when someone claims their task, updates task status, or leaves feedback
- **Notification Management**: Users can view notifications, mark as read (single or all), and check unread count
- **Privacy**: Only notification owner can view and manage their notifications
- **Input Sanitization**: All user inputs are trimmed and length-limited to prevent abuse
- **Better Validation**: Username must be 3+ characters, email must contain @, improved error messages
- **Cascade Delete**: Deleting a task also removes associated notifications and feedback

### Backend - New Unit Tests
- 10 new tests for notification and validation functionality
- Total: 55+ backend unit tests

### Documentation
- Updated README.md with full setup instructions
- Sprint4.md with all tests and API documentation

---

## User Stories

### From Sprint 1-3
1. As a student, I want to create a task with a deadline and priority so that others understand its urgency.
2. As a student, I want to view available (open) tasks so that I can find a task to work on.
3. As a student, I want to claim an open task so that I can take responsibility for completing it.
4. As a student, I want to see the status of a task so that I can track its progress.
5. As a student, I want to edit or delete a task I created so that I can manage my tasks.
6. As a student, I want to register an account so that I can create and manage tasks securely.
7. As a student, I want to log in to my account so that I can access the system.
8. As a student, I want to search for tasks by keyword so that I can find specific tasks quickly.
9. As a student, I want to filter tasks by priority, status, deadline range, and creator.
10. As a student, I want to create and update my profile so others can see my skills.
11. As a student, I want to view other students' profiles.
12. As a student, I want to leave feedback on completed tasks.
13. As a student, I want to see feedback on tasks.

### New for Sprint 4
14. As a student, I want to receive notifications when someone claims my task so I stay informed.
15. As a student, I want to be notified when a task's status changes so I can track progress.
16. As a student, I want to be notified when someone leaves feedback on my task.
17. As a student, I want to mark notifications as read so I can manage my inbox.
18. As a student, I want to see how many unread notifications I have.
19. As a student, I want my notifications to be private so only I can see them.

---

## Backend Unit Tests

Run with:
```bash
go test ./handlers/ -v
```

### Register Tests
| Test Name | Description |
|-----------|-------------|
| TestRegister_Success | Registers user, verifies 201 |
| TestRegister_MissingFields | Missing data returns 400 |
| TestRegister_DuplicateUsername | Duplicate returns 409 |
| TestRegister_ShortPassword | Password < 6 chars returns 400 |
| TestRegister_ShortUsername | Username < 3 chars returns 400 (Sprint 4) |
| TestRegister_InvalidEmail | Email without @ returns 400 (Sprint 4) |

### Login Tests
| Test Name | Description |
|-----------|-------------|
| TestLogin_Success | Correct credentials return 200 |
| TestLogin_WrongPassword | Wrong password returns 401 |
| TestLogin_NonExistentUser | Non-existent user returns 401 |
| TestLogin_MissingFields | Missing data returns 400 |

### CreateTask Tests
| Test Name | Description |
|-----------|-------------|
| TestCreateTask_Success | Creates task with all fields |
| TestCreateTask_MissingFields | Missing data returns 400 |
| TestCreateTask_InvalidDeadline | Bad date returns 400 |
| TestCreateTask_DefaultPriority | No priority defaults to Medium |
| TestCreateTask_UnregisteredUser | Unregistered user returns 400 |
| TestCreateTask_WithCategoryAndLocation | Verifies category and location stored |

### GetTasks Tests
| Test Name | Description |
|-----------|-------------|
| TestGetTasks_Empty | Empty DB returns empty array |
| TestGetTasks_FilterByStatus | Filters by status=open |
| TestGetTasks_FilterByCreatedUser | Filters by created_by |
| TestGetTasks_Search | Searches by keyword |
| TestGetTasks_SearchCaseInsensitive | Case-insensitive search |
| TestGetTasks_FilterByPriority | Filters by priority=High |
| TestGetTasks_FilterByDeadlineRange | Filters by deadline_before |
| TestGetTasks_SortByDeadline | Sorts earliest first |
| TestGetTasks_SortByPriority | Sorts High > Medium > Low |

### GetTask Tests
| Test Name | Description |
|-----------|-------------|
| TestGetTask_Success | Gets task by ID |
| TestGetTask_NotFound | 404 for non-existent |

### UpdateTask Tests
| Test Name | Description |
|-----------|-------------|
| TestUpdateTask_Success | Creator updates task |
| TestUpdateTask_Forbidden | Non-creator gets 403 |

### DeleteTask Tests
| Test Name | Description |
|-----------|-------------|
| TestDeleteTask_Success | Creator deletes task |
| TestDeleteTask_Forbidden | Non-creator gets 403 |

### ClaimTask Tests
| Test Name | Description |
|-----------|-------------|
| TestClaimTask_Success | Claims open task |
| TestClaimTask_AlreadyClaimed | Already claimed returns 409 |

### UpdateTaskStatus Tests
| Test Name | Description |
|-----------|-------------|
| TestUpdateTaskStatus_Success | Updates to in_progress |
| TestUpdateTaskStatus_InvalidStatus | Invalid status returns 400 |
| TestUpdateTaskStatus_Forbidden | Unauthorized gets 403 |
| TestMarkTaskAsCompleted | Marks task as done |

### Profile Tests
| Test Name | Description |
|-----------|-------------|
| TestGetProfile_Success | Gets profile |
| TestGetProfile_NotFound | 404 for non-existent |
| TestUpdateProfile_Success | Updates profile fields |
| TestUpdateProfile_Forbidden | Can't edit other's profile |

### Feedback Tests
| Test Name | Description |
|-----------|-------------|
| TestAddFeedback_Success | Adds 5-star feedback |
| TestAddFeedback_InvalidRating | Rating > 5 returns 400 |
| TestAddFeedback_TaskNotDone | Non-completed task returns 400 |
| TestAddFeedback_DuplicateFeedback | Duplicate returns 409 |
| TestGetFeedback_Success | Gets feedback list |
| TestGetFeedback_Empty | No feedback returns empty |

### Notification Tests (Sprint 4 - NEW)
| Test Name | Description |
|-----------|-------------|
| TestGetNotifications_Empty | No notifications returns empty array |
| TestGetNotifications_MissingUsername | Missing username returns 400 |
| TestClaimTask_CreatesNotification | Claiming task notifies creator |
| TestStatusUpdate_CreatesNotification | Status change notifies other party |
| TestMarkNotificationRead_Success | Marks single notification read |
| TestMarkNotificationRead_Forbidden | Can't mark other's notification |
| TestMarkAllNotificationsRead_Success | Marks all as read |
| TestGetUnreadCount | Returns correct unread count |
| TestRegister_ShortUsername | Username < 3 chars returns 400 |
| TestRegister_InvalidEmail | Invalid email returns 400 |

---

## Frontend Unit Tests

*Listed by frontend team*

### Cypress Tests

*Listed by frontend team*

---

## Backend API Documentation

### Base URL
```
http://localhost:8080/api
```

### Sprint 4 New Endpoints

#### GET /api/notifications?username=xxx
Get all notifications for a user (most recent first, max 50).

**Success Response (200):**
```json
[
    {
        "id": 1,
        "username": "Chinmai",
        "message": "alice claimed your task: ML Assignment",
        "task_id": 1,
        "is_read": false,
        "created_at": "2026-04-01T12:00:00Z"
    }
]
```

#### GET /api/notifications/unread-count?username=xxx
Get the count of unread notifications.

**Success Response (200):**
```json
{
    "unread_count": 3
}
```

#### PATCH /api/notifications/{id}/read?username=xxx
Mark a single notification as read. Only the notification owner can do this.

**Success Response (200):**
```json
{
    "message": "Notification marked as read"
}
```

**Error:** 403 if trying to read another user's notification

#### PATCH /api/notifications/read-all?username=xxx
Mark all of a user's notifications as read.

**Success Response (200):**
```json
{
    "message": "All notifications marked as read"
}
```

### When Notifications Are Created

| Event | Who Gets Notified | Message |
|-------|------------------|---------|
| Task claimed | Task creator | "{claimer} claimed your task: {title}" |
| Status updated by claimer | Task creator | "Task '{title}' status changed to {status} by {claimer}" |
| Status updated by creator | Task claimer | "Task '{title}' status changed to {status} by {creator}" |
| Feedback left | The other party | "{username} left feedback on task '{title}'" |

### All Endpoints Summary

| Method | Endpoint | Description | Sprint |
|--------|----------|-------------|--------|
| POST | /api/register | Register user | 2 |
| POST | /api/login | Login | 2 |
| GET | /api/profile/{username} | Get profile | 3 |
| PUT | /api/profile/{username}?username=xxx | Update profile | 3 |
| POST | /api/tasks | Create task | 1 |
| GET | /api/tasks | List/filter/search tasks | 1+2 |
| GET | /api/tasks/{id} | Get single task | 1 |
| PUT | /api/tasks/{id}?username=xxx | Edit task | 1 |
| DELETE | /api/tasks/{id}?username=xxx | Delete task | 1 |
| POST | /api/tasks/{id}/claim | Claim task | 1 |
| PATCH | /api/tasks/{id}/status?username=xxx | Update status | 1 |
| POST | /api/tasks/{id}/feedback?username=xxx | Add feedback | 3 |
| GET | /api/tasks/{id}/feedback | Get feedback | 3 |
| GET | /api/notifications?username=xxx | Get notifications | 4 |
| GET | /api/notifications/unread-count?username=xxx | Unread count | 4 |
| PATCH | /api/notifications/{id}/read?username=xxx | Mark read | 4 |
| PATCH | /api/notifications/read-all?username=xxx | Mark all read | 4 |