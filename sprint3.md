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

