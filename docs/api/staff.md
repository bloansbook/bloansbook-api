# Staff API

**Base path:** `/api/v1/staff`

The staff module manages all personnel records — from creation through role management to termination. Every person in the business has a staff record. Login-enabled staff can authenticate and access the system; no-login staff are managed on their behalf by supervisors.

All endpoints require authentication. Permission requirements are noted per endpoint.

---

## Data Models

### Staff Object (full)

Returned by `GET /staff/:id`.

| Field | Type | Description |
|---|---|---|
| `id` | UUID | Internal identifier |
| `staffId` | string | Business ID e.g. `BLN-0001` — used for login |
| `firstName` | string | First name |
| `lastName` | string | Last name |
| `email` | string \| null | Contact email |
| `phone` | string \| null | Phone number |
| `address` | string \| null | Home address |
| `dateOfBirth` | string \| null | ISO 8601 date |
| `dateOfHire` | string | ISO 8601 date |
| `emergencyContactName` | string \| null | Emergency contact full name |
| `emergencyContactPhone` | string \| null | Emergency contact number |
| `bankName` | string \| null | Bank name for payroll |
| `bankAccountNumber` | string \| null | Bank account number |
| `bankAccountName` | string \| null | Account name as on bank records |
| `department` | string | One of: `factory`, `admin`, `sales`, `management` |
| `jobTitle` | string | Descriptive title e.g. `Tailor`, `Sales Manager` |
| `payType` | string | `monthly` (only supported value in v1) |
| `baseSalary` | string (decimal) | Monthly gross base salary in NGN |
| `status` | string | One of: `active`, `inactive`, `fired` |
| `firedAt` | string \| null | Timestamp of termination, if applicable |
| `hasLogin` | boolean | Whether this staff member can log in |
| `createdBy` | StaffSummary | Compact profile of the creator |
| `roles` | RoleSummary[] | Currently assigned roles |
| `createdAt` | string | ISO 8601 timestamp |
| `updatedAt` | string | ISO 8601 timestamp |

### Staff Summary Object

| Field | Type | Description |
|---|---|---|
| `staffId` | string | Business ID |
| `firstName` | string | First name |
| `lastName` | string | Last name |
| `email` | string \| null | Contact email |
| `phone` | string \| null | Phone number |
| `department` | string | Department |
| `jobTitle` | string | Job title |
| `status` | string | `active`, `inactive`, or `fired` |

### FireStaffResponse Object

Returned by `POST /staff/:id/fire`.

| Field | Type | Description |
|---|---|---|
| `id` | UUID | `fired_staff` record identifier |
| `staffId` | UUID | Internal ID of the terminated staff member |
| `terminationReason` | string | Reason provided at termination |
| `recordedBy` | UUID | Internal ID of the staff member who recorded the termination |
| `recordedAt` | string | ISO 8601 timestamp — also used as `fired_at` on the staff record |
| `createdAt` | string | ISO 8601 timestamp |

### OverrideTerminationResponse Object

Returned by `PATCH /staff/:id/fire/override`.

| Field | Type | Description |
|---|---|---|
| `id` | UUID | `fired_staff` record identifier |
| `staffId` | UUID | Internal ID of the staff member |
| `terminationReason` | string | Original termination reason |
| `isOverridden` | boolean | Always `true` after a successful override |
| `overriddenBy` | UUID | Internal ID of the Super Admin who overrode |
| `overriddenAt` | string | ISO 8601 timestamp of the override |
| `overrideReason` | string | Mandatory reason provided for the override |
| `recordedAt` | string | Original termination timestamp |
| `createdAt` | string | Record creation timestamp |

---

## Endpoints

### List All Staff

Returns a filtered, sorted, paginated list of all staff records.

```
GET /api/v1/staff
```

**Authentication:** Required.
**Permission:** `staff.view`

**Query Parameters**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `limit` | integer | `20` | Number of records to return |
| `offset` | integer | `0` | Number of records to skip |
| `search` | string | — | ILIKE match against `firstName`, `lastName`, `staffId` |
| `status` | string | — | Filter by status: `active`, `inactive`, or `fired` |
| `department` | string | — | Filter by department: `factory`, `admin`, `sales`, `management` |
| `sortBy` | string | `createdAt` | Column to sort by. Allowed values: `createdAt`, `firstName`, `lastName`, `staffId`, `department`, `status` |
| `sortOrder` | string | `desc` | `asc` or `desc` |

**Example requests**

```
GET /api/v1/staff?status=active&department=factory&limit=20&offset=0
GET /api/v1/staff?search=john&sortBy=firstName&sortOrder=asc
GET /api/v1/staff?status=fired&limit=10
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Staff records fetched successfully",
  "statusCode": 200,
  "data": {
    "data_items": [
      {
        "id": "018e1c2d-...",
        "staffId": "BLN-0001",
        "firstName": "Amara",
        "lastName": "Okafor",
        "department": "management",
        "jobTitle": "General Manager",
        "status": "active",
        "hasLogin": true,
        "roles": [],
        "createdAt": "2026-01-15T09:00:00Z",
        "updatedAt": "2026-01-15T09:00:00Z"
      }
    ],
    "count": 1,
    "total_count": 42,
    "limit": 20,
    "offset": 0
  }
}
```

---

### Get Staff by ID

Retrieve the full profile of a single staff member including their current roles.

```
GET /api/v1/staff/:id
```

**Authentication:** Required.
**Permission:** `staff.view`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The staff member's internal `id` |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Staff record fetched successfully",
  "statusCode": 200,
  "data": {
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
    "hasLogin": true,
    "createdBy": {
      "staffId": "SYSTEM-0000",
      "firstName": "System",
      "lastName": "Account",
      "department": "system",
      "jobTitle": "system_account",
      "status": "active"
    },
    "roles": [
      { "id": "uuid...", "name": "super_admin" }
    ],
    "createdAt": "2026-01-15T09:00:00Z",
    "updatedAt": "2026-01-15T09:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | `:id` is not a valid UUID |
| `500` | `staff not found: ...` | No record found for the given ID |

---

### Create Staff

Create a new staff member. The system auto-generates a `staffId`, a temporary password, and sends a welcome email if `hasLogin: true`.

```
POST /api/v1/staff
```

**Authentication:** Required.
**Permission:** `staff.create`

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `firstName` | string | Yes | First name |
| `lastName` | string | Yes | Last name |
| `email` | string | No | Contact email — required if `hasLogin: true` |
| `phone` | string | No | Phone number |
| `address` | string | No | Home address |
| `dateOfBirth` | string | No | ISO 8601 date |
| `dateOfHire` | string | Yes | ISO 8601 date |
| `emergencyContactName` | string | No | Emergency contact name |
| `emergencyContactPhone` | string | No | Emergency contact phone |
| `bankName` | string | No | Bank name |
| `bankAccountNumber` | string | No | Bank account number |
| `bankAccountName` | string | No | Account name as on bank records |
| `department` | string | Yes | One of: `factory`, `admin`, `sales`, `management` |
| `jobTitle` | string | Yes | Job title |
| `payType` | string | Yes | `monthly` |
| `baseSalary` | string (decimal) | Yes | Monthly gross base salary in NGN |
| `status` | string | Yes | One of: `active`, `inactive`, `fired` |
| `hasLogin` | boolean | No | Default: `false`. Set `true` to create Supabase Auth account and send credentials |

**Success Response — `201 Created`**

```json
{
  "status": "success",
  "message": "Staff record created successfully",
  "statusCode": 201,
  "data": {
    "id": "018e1c2d-...",
    "staff_info": {
      "staffId": "BLN-0002",
      "firstName": "Funmi",
      "lastName": "Adeyemi",
      "email": "funmi@example.com",
      "phone": "+2348055001122",
      "department": "sales",
      "jobTitle": "Sales Executive",
      "status": "active"
    },
    "createdAt": "2026-07-10T11:30:00Z",
    "updatedAt": "2026-07-10T11:30:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | Missing or malformed fields |
| `401` | `Authentication required` | No valid token |
| `403` | `You do not have permission to perform this action` | Missing `staff.create` |
| `500` | `failed to create staff: ...` | Database error or `staffId` conflict |
| `500` | `failed to send welcome email: ...` | Email delivery failed (record was saved) |

---

### Update Staff

Partially update a staff record. Only fields included in the request body are changed.

```
PATCH /api/v1/staff/:id
```

**Authentication:** Required.
**Permission:** `staff.update`

All body fields are optional. Include only what you want to change.

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Staff record updated successfully",
  "statusCode": 200,
  "data": {
    "id": "018e1c2d-...",
    "staff": {
      "staffId": "BLN-0002",
      "firstName": "Funmi",
      "lastName": "Adeyemi",
      "department": "sales",
      "jobTitle": "Sales Executive",
      "status": "active"
    },
    "createdAt": "2026-07-10T11:30:00Z",
    "updatedAt": "2026-07-10T14:00:00Z"
  }
}
```

---

### Terminate Staff (Fire)

Terminates a staff member. This is a permanent, auditable action:
- Inserts a record into `fired_staff`
- Deletes all current `staff_roles` and writes a `revoked` entry per role into `staff_role_history`
- Sets `staff.status = 'fired'` and `staff.fired_at` to the recorded termination time
- The terminated staff member cannot log in after this — the auth middleware blocks them on their next request

```
POST /api/v1/staff/:id/fire
```

**Authentication:** Required.
**Permission:** `staff.terminate`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The staff member's internal `id` |

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `terminationReason` | string | Yes | Mandatory reason for termination e.g. `Dismissed for misconduct` |

```json
{
  "terminationReason": "Dismissed for misconduct"
}
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Staff member has been terminated",
  "statusCode": 200,
  "data": {
    "id": "fired-staff-uuid...",
    "staffId": "018e1c2d-...",
    "terminationReason": "Dismissed for misconduct",
    "recordedBy": "caller-uuid...",
    "recordedAt": "2026-07-27T10:00:00Z",
    "createdAt": "2026-07-27T10:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | `:id` is not a valid UUID or body is malformed |
| `401` | `Authentication required` | No valid token |
| `403` | `You do not have permission to perform this action` | Missing `staff.terminate` |
| `409` | `This staff member has already been terminated` | Staff status is already `fired` |
| `500` | `staff not found: ...` | No record for the given ID |

---

### Override Termination

Overrides an incorrect termination record. Restores the staff member to `active` status and clears `fired_at`. The `fired_staff` record is **not deleted** — it is permanently marked as overridden with a mandatory reason and the identity of the Super Admin who overrode it.

This endpoint should only be accessible to super_admin in practice — wire it behind the `staff.terminate` permission which only super_admin holds by default.

```
PATCH /api/v1/staff/:id/fire/override
```

**Authentication:** Required.
**Permission:** `staff.terminate`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The staff member's internal `id` |

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `overrideReason` | string | Yes | Mandatory reason for overriding the termination |

```json
{
  "overrideReason": "Termination was recorded in error — staff member was misidentified"
}
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Termination record has been overridden and staff member restored",
  "statusCode": 200,
  "data": {
    "id": "fired-staff-uuid...",
    "staffId": "018e1c2d-...",
    "terminationReason": "Dismissed for misconduct",
    "isOverridden": true,
    "overriddenBy": "super-admin-uuid...",
    "overriddenAt": "2026-07-28T09:00:00Z",
    "overrideReason": "Termination was recorded in error — staff member was misidentified",
    "recordedAt": "2026-07-27T10:00:00Z",
    "createdAt": "2026-07-27T10:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | `:id` is not a valid UUID or body is malformed |
| `401` | `Authentication required` | No valid token |
| `403` | `You do not have permission to perform this action` | Missing `staff.terminate` |
| `404` | `Staff record not found` | No active (non-overridden) termination record exists for this staff member |
| `500` | `staff not found: ...` | No staff record for the given ID |

---

## Role Management Endpoints

These endpoints manage which roles are assigned to a staff member. All role changes are recorded permanently in `staff_role_history`. All require the `auth.manage_roles` permission.

---

### Assign Role

```
POST /api/v1/staff/:id/roles
```

**Permission:** `auth.manage_roles`

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `roleId` | UUID | Yes | The role to assign |
| `reason` | string | No | Optional reason for the audit trail |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Role assigned successfully",
  "statusCode": 200,
  "data": {
    "staffId": "018e1c2d-...",
    "roleId": "b3f2a1d0-...",
    "roleName": "hr_manager",
    "action": "assigned"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `500` | `role already assigned to this staff member` | Duplicate assignment |
| `500` | `staff not found: ...` | No staff record |

---

### Revoke Role

```
DELETE /api/v1/staff/:id/roles
```

**Permission:** `auth.manage_roles`

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `roleId` | UUID | Yes | The role to revoke |
| `reason` | string | No | Optional reason |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Role revoked successfully",
  "statusCode": 200,
  "data": {
    "staffId": "018e1c2d-...",
    "roleId": "b3f2a1d0-...",
    "roleName": "hr_manager",
    "action": "revoked"
  }
}
```

---

### Update Role (Swap)

Atomically revokes the old role and assigns the new one in a single transaction.

```
PUT /api/v1/staff/:id/roles
```

**Permission:** `auth.manage_roles`

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `oldRoleId` | UUID | Yes | The role to remove |
| `newRoleId` | UUID | Yes | The role to assign |
| `reason` | string | No | Applies to both history entries |

---

### Get Role History

Returns the complete role change log for a staff member, newest first.

```
GET /api/v1/staff/:id/roles/history
```

**Permission:** `auth.manage_roles`

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Role fetched successfully",
  "statusCode": 200,
  "data": [
    {
      "id": "uuid...",
      "roleId": "b3f2a1d0-...",
      "roleName": "hr_manager",
      "action": "revoked",
      "performedBy": "018e1c2d-...",
      "reason": "Termination from job: Dismissed for misconduct",
      "createdAt": "2026-07-27T10:00:00Z"
    }
  ]
}
```

---

## Business Rules

- `staffId` is auto-generated at creation in the format `{PREFIX}-NNNN` e.g. `BLN-0001`. It never changes.
- Staff records are never deleted. Use `status: inactive` to deactivate or `POST /:id/fire` to terminate.
- Termination (`POST /:id/fire`) is a multi-step atomic transaction — all steps succeed or none do.
- An incorrect termination is corrected via `PATCH /:id/fire/override`, not by deletion. The original record is preserved with `is_overridden = true`.
- The `fired_staff` table may only have one non-overridden record per staff member. Attempting to fire an already-fired staff member returns `409 Conflict`.
- All role assignment, revocation, and swap operations are transactional. The `staff_role_history` table is append-only.
- Firing a staff member automatically revokes all their roles and writes one `revoked` history entry per role.
- The caller's identity (`recordedBy`, `performedBy`, `overriddenBy`) is always taken from the JWT — never from the request body.
- A fired staff member's login is blocked immediately on their next request — the auth middleware checks `status = 'active'` on every request.
- When an override is applied, `staff.status` is restored to `active` and `fired_at` is cleared, restoring the staff member's login access.
