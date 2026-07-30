# Authentication API

**Base path:** `/api/v1/auth`

The authentication module handles staff login, session management, and current-user resolution. Login is the only fully public endpoint — refresh is also public (Supabase validates the token itself). Every other endpoint requires a valid session token.

Authentication is backed by Supabase Auth. Tokens are delivered and stored as **HttpOnly cookies**, never in the response body. This means the frontend does not need to manually store or attach the token — the browser handles it automatically on every request.

---

## How Authentication Works

1. Staff member sends `staffId` and `password` to `POST /auth/login`.
2. The API checks that the staff record exists and that `has_login = true`.
3. The API checks that the staff member is not fired or inactive.
4. Credentials are forwarded to Supabase Auth for verification.
5. On success, two HttpOnly cookies are set on the response:
   - `access_token` — short-lived Supabase JWT (~1 hour). Sent automatically by the browser on every subsequent request.
   - `refresh_token` — longer-lived token (7 days). Scoped to `/api/v1/auth/refresh` only.
6. The auth middleware on each protected route reads the `access_token` cookie, verifies it with Supabase, resolves the staff record, and loads their permissions into the request context.

---

## Cookie Strategy

| Cookie | Path | Expiry | Purpose |
|---|---|---|---|
| `access_token` | `/` | ~1 hour (Supabase-managed) | Sent on every API request |
| `refresh_token` | `/api/v1/auth/refresh` | 7 days | Sent only to the refresh endpoint |

Both cookies are `HttpOnly`, `Secure`, and `SameSite=Lax`.

**CORS note:** Because auth uses cookies, `AllowCredentials: true` must be set on the backend and the frontend must call all API endpoints with `credentials: "include"`. In development (`ALLOWED_ORIGINS=*`) credentials are disabled — set a real origin in staging and production.

---

## Token Refresh Strategy

The frontend does **not** interact with Supabase directly. All token operations go through this API.

When the access token expires, the backend returns `401` with message `Invalid session token`. The frontend should:
1. Call `POST /api/v1/auth/refresh` (no body needed — the `refresh_token` cookie is sent automatically).
2. On success, a new `access_token` cookie is set. Retry the original request.
3. If refresh fails (refresh token also expired), redirect to login.

The refresh endpoint is **public** — no access token required — because the refresh token itself proves the caller's identity (Supabase validates it).

---

## Endpoints

### Login

```
POST /api/v1/auth/login
```

**Authentication:** None — public endpoint.

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `staffId` | string | Yes | Business ID e.g. `BLN-0001` |
| `password` | string | Yes | Staff member's password |

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
  "data": null
}
```

Tokens are delivered as HttpOnly cookies, not in the response body. `data` is always `null` on login.

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | Missing or malformed body |
| `401` | `Staff not found: ...` | No staff record for the given `staffId` |
| `401` | `staff does not have login access` | Staff record exists but `has_login = false` |
| `401` | `Invalid credentials: ...` | Password is wrong |
| `403` | `Access denied` | Staff status is `inactive` or `fired` — caught by middleware on the first protected request after login |

---

### Refresh Token

Exchange a valid refresh token for a new access token. The `refresh_token` cookie is sent automatically by the browser.

```
POST /api/v1/auth/refresh
```

**Authentication:** None — public. Supabase validates the refresh token.

**Request Body:** None required.

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "",
  "statusCode": 200,
  "data": null
}
```

Both cookies are rotated — a new `access_token` and `refresh_token` are set.

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `401` | `Authentication required` | No `refresh_token` cookie present |
| `401` | `Invalid session token` | Refresh token is expired or invalid |

---

### Get Current User (`/me`)

Returns the full profile of the currently authenticated staff member, their roles (with grouped permissions), and a flat list of all permission codes.

```
GET /api/v1/auth/me
```

**Authentication:** Required (`access_token` cookie).

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "User profile fetched successfully",
  "statusCode": 200,
  "data": {
    "me": {
      "id": "018e1c2d-...",
      "staffId": "BLN-0001",
      "firstName": "Amara",
      "lastName": "Okafor",
      "email": "amara@example.com",
      "phone": "+2348012345678",
      "address": "12 Ikeja Road, Lagos",
      "dateOfBirth": "1990-03-22T00:00:00Z",
      "dateOfHire": "2024-01-10T00:00:00Z",
      "emergencyContactName": "Chisom Okafor",
      "emergencyContactPhone": "+2348099887766",
      "bankName": "Access Bank",
      "bankAccountNumber": "1234567890",
      "bankAccountName": "AMARA OKAFOR",
      "department": "management",
      "jobTitle": "General Manager",
      "payType": "monthly",
      "baseSalary": "350000.00",
      "status": "active",
      "firedAt": null,
      "createdAt": "2026-01-15T09:00:00Z",
      "updatedAt": "2026-01-15T09:00:00Z"
    },
    "roles": [
      {
        "id": "b3f2a1d0-...",
        "name": "hr_manager",
        "permissions": [
          {
            "module": "payroll",
            "codes": ["payroll.approve", "payroll.view"]
          },
          {
            "module": "staff",
            "codes": ["staff.create", "staff.update", "staff.view"]
          }
        ]
      }
    ],
    "permissions": [
      "payroll.approve",
      "payroll.view",
      "staff.create",
      "staff.update",
      "staff.view"
    ]
  }
}
```

**Response fields explained**

| Field | Type | Description |
|---|---|---|
| `me` | object | Full staff profile of the authenticated user |
| `roles` | array | All assigned roles. Each role contains `permissions` grouped by module — useful for building permission UIs |
| `permissions` | string[] | Flat, deduplicated list of every permission code the user holds across all roles — use this for `hasPermission("staff.create")` checks |

**Frontend usage pattern**

```js
// After calling /me, store permissions as a Set for O(1) lookups
const perms = new Set(data.permissions)
const hasPermission = (code) => perms.has(code)

// Use to show/hide UI elements
if (hasPermission("staff.create")) { /* show Create Staff button */ }
```

**Important:** Frontend permission checks are UI-only. The backend enforces all permissions server-side on every request regardless of what the frontend does.

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `401` | `Authentication required` | No `access_token` cookie |
| `401` | `Invalid session token` | Token expired or invalid |
| `403` | `Access denied` | Staff status is `inactive` or `fired` |

---

## Auth Middleware Behaviour

Every protected route runs the following checks in order:

1. Reads `access_token` cookie — `401` if missing.
2. Calls Supabase to verify the token — `401` if invalid.
3. Looks up the staff record by `supabase_uid` — `401` if not found.
4. Checks `staff.status == "active"` — `403` if inactive or fired.
5. Loads the staff member's permission codes and role names into the request context.

Steps 4 and 5 run on **every** request — permission changes and firings take effect immediately on the next request, with no need to re-login.

---

## Business Rules

- Only staff with `has_login = true` can authenticate.
- Staff with `status = fired` are blocked at the `fired_staff` level — the middleware denies access even with a valid token.
- Staff with `status = inactive` are blocked by the same middleware check.
- Each staff member has a synthetic Supabase email `{staffId}@bloansbook.local` — not a real email.
- Passwords are hashed with bcrypt. The plaintext is never stored.
- Login credentials are emailed to the staff member automatically when their record is created with `has_login = true`.

---

## CORS & Credential Configuration

| Environment | `ALLOWED_ORIGINS` value | `AllowCredentials` |
|---|---|---|
| Local dev | *(unset or `*`)* | `false` — cookies still set but browser won't send them cross-origin |
| Staging | `http://localhost:3000` | `true` |
| Production | `https://app.bloansbook.com` | `true` |

Multiple origins are comma-separated:
```
ALLOWED_ORIGINS=http://localhost:3000,https://app.bloansbook.com
```
