# Authentication API

**Base path:** `/api/v1/auth`

The authentication module handles staff login and session management. Login is the only public endpoint in the system — every other endpoint requires a valid Bearer token issued by this module.

Authentication is backed by Supabase Auth. When a staff member logs in, the API exchanges their credentials with Supabase and returns a JWT access token. That token must be sent on all subsequent requests in the `Authorization` header.

---

## How Authentication Works

1. Staff member sends their `staffId` and `password` to `POST /auth/login`.
2. The API validates the staff record exists and that `has_login = true`.
3. Credentials are forwarded to Supabase Auth for verification.
4. On success, a JWT access token is returned alongside the staff profile and their assigned roles.
5. The client stores the token and sends it as `Authorization: Bearer <token>` on every subsequent request.
6. The auth middleware on each protected route verifies the token with Supabase, resolves the staff record, and loads their permissions into the request context.

---

## Endpoints

### Login

Authenticate a staff member and receive an access token.

```
POST /api/v1/auth/login
```

**Authentication:** None — public endpoint.

**Request Body**

| Field      | Type   | Required | Description                                        |
|------------|--------|----------|----------------------------------------------------|
| `staffId`  | string | Yes      | The staff member's business ID e.g. `BLN-0001`     |
| `password` | string | Yes      | The staff member's password                        |

```json
{
  "staffId": "BLN-0001",
  "password": "your-password"
}
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Login successful",
  "statusCode": 200,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "staff": {
      "staffId": "BLN-0001",
      "firstName": "Amara",
      "lastName": "Okafor",
      "email": "amara@example.com",
      "phone": "+2348012345678",
      "department": "management",
      "jobTitle": "General Manager",
      "status": "active"
    },
    "roles": ["super_admin"]
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400` | `Invalid request data` | Missing or malformed request body |
| `401` | `Staff not found: ...` | No staff record matches the provided `staffId` |
| `401` | `Staff has no login access: ...` | Staff record exists but `has_login = false` |
| `401` | `Invalid credentials: ...` | Password is incorrect |

---

## Response Envelope

All responses from this API — success and error — follow the same envelope shape:

```json
{
  "status": "success" | "error",
  "message": "Human-readable message",
  "statusCode": 200,
  "data": { ... } | null
}
```

---

## Using the Token

Include the access token on every protected request:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

The token is a Supabase-issued JWT. It encodes the staff member's Supabase UID, which the middleware uses to resolve their full staff record and permission set on each request.

**Token expiry:** Managed by Supabase. When a token expires, the client will receive a `401` with message `Invalid session token`. The client must re-authenticate via `POST /auth/login`.

---

## Business Rules

- Only staff with `has_login = true` can authenticate. Staff without login access (e.g. factory workers) are excluded.
- Staff with `status = fired` or `status = inactive` are blocked by the auth middleware even with a valid token.
- Login credentials (`staffId` and initial password) are emailed to the staff member automatically when their record is created.
- Passwords are hashed with bcrypt before storage. The plaintext password is never stored.
- Each staff member has a synthetic Supabase email in the format `{staffId}@bloansbook.local` used internally for Supabase Auth — it is not a real email address.
