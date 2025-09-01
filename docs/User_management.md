## User & Auth API

This document describes the API that is currently implemented. All endpoints return JSON. Base path: `/api/v1`.

### Auth model

- Authentication: Bearer JWT in `Authorization: Bearer <access_token>` for protected endpoints.
- Token refresh: Uses long-lived refresh tokens. See Refresh Token.
- Errors: `{ "error": string }` with appropriate HTTP status code.
- Rate limiting: Global limiter is enabled; repeated requests may return 429.

---

## Authentication endpoints

### Register

- POST `/api/v1/auth/register`
- Body:
  - `username` string (3–32)
  - `email` valid email
  - `password` string (8–32, must include upper, lower, digit, symbol)
  - `fullname` string (3–50)
- Responses:
  - 201: `{ "message": "User created successfully. Please check your email to verify your account." }`
  - 409: `{ "error": "..." }` if conflict (e.g., email/username taken)

### Login

- POST `/api/v1/auth/login`
- Body:
  - `email` string, `password` string
- Responses:
  - 200: `{ "user": User, "access_token": string, "refresh_token": string }`
  - 401: `{ "error": "Invalid credentials or unverified email" }`

### Verify email

- GET `/api/v1/auth/verify-email?verifier=<id>&token=<token>`
- Responses:
  - 200: `{ "message": "Email verified successfully", "user": User }`
  - 400: `{ "error": "invalid token or expired token" }`

### Request verification email

- POST `/api/v1/auth/request-verification-email`
- Body: `{ "user_id": string }`
- Responses:
  - 200: `{ "message": "Verification email sent successfully" }`
  - 400/404/500: `{ "error": "..." }`

### Forgot password

- POST `/api/v1/auth/forgot-password`
- Body: `{ "email": string }`
- Response: 200 with message (even if email does not exist)

### Reset password

- POST `/api/v1/auth/reset-password`
- Body: `{ "token": string, "verifier": string, "password": string(min 8) }`
- Responses:
  - 200: `{ "message": "Password reset successfully" }`
  - 400: `{ "error": "Invalid or expired reset token" }`

### Refresh token

- POST `/api/v1/auth/refresh-token`
- Body: `{ "refresh_token": string }`
- Responses:
  - 200: `{ "access_token": string, "refresh_token": string }`
  - 401: `{ "error": "Invalid or expired refresh token" }`

### Logout

- POST `/api/v1/logout`
- Body: `{ "refresh_token": string }`
- Responses:
  - 200: `{ "message": "Logged out successfully" }`
  - 400/500: `{ "error": "..." }`

### Google OAuth

- GET `/api/v1/auth/google/login` → Redirects to Google OAuth consent.
- GET `/api/v1/auth/google/callback?code=...&state=...`
  - Success 200: `{ "message": "login successful", "access token": string, "refresh token": string }`

---

## User endpoints

### Get public profile

- GET `/api/v1/users/profile/:id`
- Response:
  - 200: `User`
  - 404: `{ "error": "User not found" }`

### Get current user (me)

- GET `/api/v1/me` (auth required)
- Header: `Authorization: Bearer <access_token>`
- Response: 200: `User`

### Update profile

- PUT `/api/v1/me` (auth required)
- Body (any subset):
  - `username` string (3–32)
  - `fullname` string (<=50)
- Response: 200: `User`

---

## Subscriptions (me)

### List subscriptions

- GET `/api/v1/me/subscriptions` (auth required)
- Response: 200: `{ "subscriptions": SubscriptionDetail[], "total_subscriptions": number, "subscription_limit": number }`

### Add subscription

- POST `/api/v1/me/subscriptions` (auth required)
- Body: `{ "source_key": string }` (source slug)
- Responses:
  - 201: (no body)
  - 403: `{ "error": "subscription limit reached" }`
  - 400/500: `{ "error": "..." }`

### Remove subscription

- DELETE `/api/v1/me/subscriptions/:source_slug` (auth required)
- Response: 200: `{ "message": "Successfully unsubscribed" }`

---

## Admin endpoints (overview)

These routes are protected; only admins should be allowed. Current implementation checks `userID == "admin"` in code (subject to change).

### Create topic

- POST `/api/v1/topics` (auth required)
- Body:
  - `slug` string
  - `label` { `en`: string, `am`: string }
- Responses: 201 on success; errors return `{ "error": "..." }`

### Create source

- POST `/api/v1/sources` (auth required)
- Body:
  - `slug` string
  - `name` string
  - `description` string
  - `url` string
  - `logo_url` string
  - `languages` string (e.g., "EN", "AM")
  - `topics` string[]
  - `reliability_score` number
- Responses: 201 on success; errors return `{ "error": "..." }`

---

## Schemas

### User

```
{
	"id": string,
	"username": string,
	"fullname": string,
	"email": string,
	"role": string,
	"first_name": string|null,
	"last_name": string|null,
	"avatar_url": string|null,
	"created_at": RFC3339 string,
	"preferences": {
		"lang": string,
		"topics": string[],
		"subscribed_sources": string[],
		"brief_type": string,
		"data_saver": boolean,
		"notifications": { "daily_brief": boolean, "breaking_news": boolean }
	}
}
```

### SubscriptionDetail

```
{
	"source_slug": string,
	"source_name": string,
	"subscribed_at": string (optional),
	"topics": string[]
}
```

---

## Notes

- CORS is enabled for all origins and common methods/headers.
- All times are returned as RFC3339 strings.
- Error messages are intentionally generic for security on auth flows.
