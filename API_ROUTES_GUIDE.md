# API Routes Guide

This document outlines the suggested REST API routes for the IQ Theory application based on the service layer interfaces. All routes are prefixed with `/api` as configured in the main router.

## Base URL Structure

```
https://yourdomain.com/api/{endpoint}
```

## Authentication

Most endpoints require authentication. The following endpoints are public:

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/health`

For authenticated endpoints, include the authorization header:

```
Authorization: Bearer {jwt_token}
```

---

## 🔐 Authentication & User Management

### Authentication Routes

| Method | Route                       | Description          | Service Method                 |
| ------ | --------------------------- | -------------------- | ------------------------------ |
| `POST` | `/api/auth/register`        | Register a new user  | `UserService.CreateUser`       |
| `POST` | `/api/auth/login`           | Authenticate user    | `UserService.AuthenticateUser` |
| `POST` | `/api/auth/logout`          | Logout user          | N/A (client-side)              |
| `POST` | `/api/auth/change-password` | Change user password | `UserService.ChangePassword`   |
| `POST` | `/api/auth/verify-email`    | Verify email address | `UserService.VerifyEmail`      |

### User Management Routes

| Method   | Route                                   | Description                     | Service Method                  |
| -------- | --------------------------------------- | ------------------------------- | ------------------------------- |
| `GET`    | `/api/users/me`                         | Get current user profile        | `UserService.GetUserByID`       |
| `PUT`    | `/api/users/me`                         | Update current user profile     | `UserService.UpdateProfile`     |
| `DELETE` | `/api/users/me`                         | Delete current user account     | `UserService.DeleteUser`        |
| `GET`    | `/api/users/{userID}`                   | Get user by ID (public profile) | `UserService.GetUserByID`       |
| `GET`    | `/api/users/search?username={username}` | Search user by username         | `UserService.GetUserByUsername` |

---

## 👥 Friendship Management

| Method   | Route                                 | Description                 | Service Method                           |
| -------- | ------------------------------------- | --------------------------- | ---------------------------------------- |
| `GET`    | `/api/friends`                        | Get user's friends list     | `FriendshipService.GetUserFriends`       |
| `GET`    | `/api/friends/requests`               | Get pending friend requests | `FriendshipService.GetPendingRequests`   |
| `POST`   | `/api/friends/request`                | Send friend request         | `FriendshipService.SendFriendRequest`    |
| `PUT`    | `/api/friends/{friendshipID}/accept`  | Accept friend request       | `FriendshipService.AcceptFriendRequest`  |
| `PUT`    | `/api/friends/{friendshipID}/decline` | Decline friend request      | `FriendshipService.DeclineFriendRequest` |
| `DELETE` | `/api/friends/{friendshipID}`         | Remove friend               | `FriendshipService.RemoveFriend`         |

### Request/Response Examples

**Send Friend Request**

```json
POST /api/friends/request
{
  "addresseeID": "123e4567-e89b-12d3-a456-426614174000"
}
```

---

## 🏫 Group Management

### Group Operations

| Method   | Route                         | Description            | Service Method                    |
| -------- | ----------------------------- | ---------------------- | --------------------------------- |
| `GET`    | `/api/groups`                 | Get user's groups      | `GroupService.GetUserGroups`      |
| `POST`   | `/api/groups`                 | Create new group       | `GroupService.CreateGroup`        |
| `GET`    | `/api/groups/{groupID}`       | Get group details      | `GroupService.GetGroupByID`       |
| `PUT`    | `/api/groups/{groupID}`       | Update group           | `GroupService.UpdateGroup`        |
| `DELETE` | `/api/groups/{groupID}`       | Delete group           | `GroupService.DeleteGroup`        |
| `GET`    | `/api/groups/join/{joinCode}` | Get group by join code | `GroupService.GetGroupByJoinCode` |

### Group Membership

| Method   | Route                                           | Description             | Service Method                  |
| -------- | ----------------------------------------------- | ----------------------- | ------------------------------- |
| `GET`    | `/api/groups/{groupID}/members`                 | Get group members       | `GroupService.GetGroupMembers`  |
| `POST`   | `/api/groups/{groupID}/join`                    | Join group by ID        | `GroupService.JoinGroupByID`    |
| `POST`   | `/api/groups/join`                              | Join group by join code | `GroupService.JoinGroup`        |
| `DELETE` | `/api/groups/{groupID}/leave`                   | Leave group             | `GroupService.LeaveGroup`       |
| `DELETE` | `/api/groups/{groupID}/members/{memberID}`      | Remove member (admin)   | `GroupService.RemoveMember`     |
| `PUT`    | `/api/groups/{groupID}/members/{memberID}/role` | Update member role      | `GroupService.UpdateMemberRole` |

### Request/Response Examples

**Create Group**

```json
POST /api/groups
{
  "name": "Music Theory Class",
  "description": "Beginner music theory study group",
  "maxMembers": 25
}
```

**Join Group by Code**

```json
POST /api/groups/join
{
  "joinCode": "ABC123"
}
```

---

## 🎵 Quiz System

### Quiz Configuration

| Method | Route                           | Description                       | Service Method                           |
| ------ | ------------------------------- | --------------------------------- | ---------------------------------------- |
| `GET`  | `/api/quiz/configurations`      | Get available quiz configurations | `QuizService.GetAvailableConfigurations` |
| `GET`  | `/api/quiz/clef-types`          | Get available clef types          | `QuizService.GetClefTypes`               |
| `GET`  | `/api/quiz/duration-options`    | Get duration options              | `QuizService.GetDurationOptions`         |
| `GET`  | `/api/quiz/ledger-line-options` | Get ledger line options           | `QuizService.GetLedgerLineOptions`       |

### Quiz Sessions

| Method | Route                                     | Description              | Service Method                    |
| ------ | ----------------------------------------- | ------------------------ | --------------------------------- |
| `POST` | `/api/quiz/sessions`                      | Create new quiz session  | `QuizService.CreateQuizSession`   |
| `GET`  | `/api/quiz/sessions`                      | Get user's quiz sessions | `QuizService.GetUserQuizSessions` |
| `GET`  | `/api/quiz/sessions/{sessionID}`          | Get quiz session details | `QuizService.GetQuizSession`      |
| `POST` | `/api/quiz/sessions/{sessionID}/start`    | Start quiz session       | `QuizService.StartQuizSession`    |
| `POST` | `/api/quiz/sessions/{sessionID}/complete` | Complete quiz session    | `QuizService.CompleteQuizSession` |
| `POST` | `/api/quiz/sessions/{sessionID}/abandon`  | Abandon quiz session     | `QuizService.AbandonQuizSession`  |

### Questions & Answers

| Method | Route                                          | Description         | Service Method                  |
| ------ | ---------------------------------------------- | ------------------- | ------------------------------- |
| `GET`  | `/api/quiz/sessions/{sessionID}/next-question` | Get next question   | `QuizService.GetNextQuestion`   |
| `POST` | `/api/quiz/sessions/{sessionID}/answers`       | Submit answer       | `QuizService.SubmitAnswer`      |
| `GET`  | `/api/quiz/sessions/{sessionID}/answers`       | Get session answers | `QuizService.GetSessionAnswers` |

### Results & Scoring

| Method | Route                                    | Description             | Service Method                      |
| ------ | ---------------------------------------- | ----------------------- | ----------------------------------- |
| `GET`  | `/api/quiz/sessions/{sessionID}/score`   | Calculate session score | `QuizService.CalculateSessionScore` |
| `GET`  | `/api/quiz/sessions/{sessionID}/results` | Get session results     | `QuizService.GetSessionResults`     |

### Request/Response Examples

**Create Quiz Session**

```json
POST /api/quiz/sessions
{
  "clef": "treble",
  "durationSeconds": 300,
  "maxLedgerLines": 2
}
```

**Submit Answer**

```json
POST /api/quiz/sessions/{sessionID}/answers
{
  "questionNumber": 1,
  "userAnswer": "A4",
  "timeTakenMs": 1250
}
```

---

## 🏆 Leaderboards

| Method | Route                               | Description                  | Service Method                            |
| ------ | ----------------------------------- | ---------------------------- | ----------------------------------------- |
| `GET`  | `/api/leaderboard/global`           | Get global leaderboard       | `LeaderboardService.GetGlobalLeaderboard` |
| `GET`  | `/api/leaderboard/user/{userID}`    | Get user ranking             | `LeaderboardService.GetUserRanking`       |
| `GET`  | `/api/leaderboard/groups/{groupID}` | Get group leaderboard        | `LeaderboardService.GetGroupLeaderboard`  |
| `POST` | `/api/leaderboard/refresh`          | Refresh leaderboards (admin) | `LeaderboardService.RefreshLeaderboards`  |

### Query Parameters

**Global Leaderboard**

```
GET /api/leaderboard/global?clef=treble&duration=300&maxLedgerLines=2&limit=50
```

**Group Leaderboard**

```
GET /api/leaderboard/groups/{groupID}?clef=bass&duration=180&maxLedgerLines=1&limit=20
```

---

## 🔧 System Routes

| Method | Route          | Description           | Purpose           |
| ------ | -------------- | --------------------- | ----------------- |
| `GET`  | `/api/health`  | Health check endpoint | System monitoring |
| `GET`  | `/api/version` | API version info      | Version tracking  |

---

## 📊 HTTP Status Codes

### Success Codes

- `200 OK` - Successful GET, PUT requests
- `201 Created` - Successful POST requests that create resources
- `204 No Content` - Successful DELETE requests

### Client Error Codes

- `400 Bad Request` - Invalid request format or parameters
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - User lacks permission for the operation
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., duplicate email)
- `422 Unprocessable Entity` - Validation errors

### Server Error Codes

- `500 Internal Server Error` - Unexpected server error

---

## 🔒 Security Considerations

### Rate Limiting

Implement rate limiting for:

- Authentication endpoints (login/register): 5 attempts per minute
- Quiz session creation: 10 per hour per user
- Friend requests: 20 per hour per user

### Input Validation

- Validate all UUID parameters
- Sanitize user inputs
- Enforce maximum payload sizes
- Validate quiz configuration parameters

### Authorization Rules

- Users can only access their own data
- Group admins can manage group members
- Quiz sessions belong to specific users
- Leaderboard data is read-only for regular users

---

## 📝 Error Response Format

All error responses follow this format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": {
      "field": "email",
      "reason": "Email address is required"
    }
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR` - Input validation failed
- `AUTHENTICATION_ERROR` - Authentication failed
- `AUTHORIZATION_ERROR` - Insufficient permissions
- `RESOURCE_NOT_FOUND` - Requested resource not found
- `RESOURCE_CONFLICT` - Resource already exists
- `RATE_LIMIT_EXCEEDED` - Too many requests

---

## 🚀 Implementation Priorities

### Phase 1 (MVP)

1. Authentication routes (`/api/auth/*`)
2. Basic user management (`/api/users/me`)
3. Quiz configuration (`/api/quiz/clef-types`, `/api/quiz/duration-options`)
4. Basic quiz sessions (`/api/quiz/sessions`)

### Phase 2 (Social Features)

1. Friendship management (`/api/friends/*`)
2. Group management (`/api/groups/*`)
3. Group leaderboards

### Phase 3 (Advanced Features)

1. Global leaderboards
2. Advanced quiz analytics
3. Real-time features (WebSocket endpoints)

---

## 📚 Additional Resources

- [OpenAPI/Swagger Documentation](./api-spec.yaml) (TODO)
- [Authentication Guide](./auth-guide.md) (TODO)
- [Rate Limiting Documentation](./rate-limiting.md) (TODO)
- [WebSocket API Documentation](./websocket-api.md) (TODO)
