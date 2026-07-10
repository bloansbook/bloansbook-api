# Staff API

**Base path:** `/api/v1/staff`

The staff module manages all personnel records — from creation through to role management. Every person in the business (factory workers, managers, sales staff) has a staff record. Login-enabled staff can authenticate and access the system; no-login staff are managed on their behalf by supervisors.

All endpoints in this module require authentication. Permission requirements are noted per endpoint.

---

## Data Models

### Staff Object (full)

Returned by `GET /staff/:id`.

| Field                   | Type            | Description |
|-------------------------|-----------------|-------------|
| `id`                    | UUID            | Internal identifier |
| `staffId`               | string          | Business ID e.g. `BLN-0001` — used for login |
| `firstName`             | string          | First name |
| `lastName`              | string          | Last name |
| `email`                 | string \| null  | Contact email |
| `phone`                 | string \| null  | Phone number |
| `address`               | string \| null  | Home address |
| `dateOfBirth`           | string \| null  | ISO 8601 date |
| `dateOfHire`            | string          | ISO 8601 date |
| `emergencyContactName`  | string \| null  | Emergency contact full name |
| `emergencyContactPhone` | string \| null  | Emergency contact number |
| `bankName`              | string \| null  | Bank name for payroll |
| `bankAccountNumber`     | string \| null  | Bank account number |
| `bankAccountName`       | string \| null  | Account name as it appears on bank records |
| `department`            | string          | One of: `factory`, `admin`, `sales`, `management` |
| `jobTitle`              | string          | Descriptive title e.g. `Tailor`, `Sales Manager` |
| `payType`               | string          | `monthly` (only supported value in v1) |
| `baseSalary`            | string (decimal)| Monthly gross base salary in NGN |
| `status`                | string          | One of: `active`, `inactive`, `fired` |
| `firedAt`               | string \| null  | Timestamp of termination, if applicable |
| `hasLogin`              | boolean         | Whether this staff member can log in |
| `createdBy`             | StaffSummary    | Compact profile of the staff member who created this record |
| `roles`                 | RoleSummary[]   | Currently assigned roles |
| `createdAt`             | string          | ISO 8601 timestamp |
| `updatedAt`             | string          | ISO 8601 timestamp |

### Staff Summary Object

Compact staff shape used inside nested responses.

| Field        | Type           | Description |
|--------------|----------------|-------------|
| `staffId`    | string         | Business ID |
| `firstName`  | string         | First name |
| `lastName`   | string         | Last name |
| `email`      | string \| null | Contact email |
| `phone`      | string \| null | Phone number |
| `department` | string         | Department |
| `jobTitle`   | string         | Job title |
| `status`     | string         | `active`, `inactive`, or `fired` |

---

## Endpoints

### List All Staff

Returns a paginated list of all staff records.

```
GET /api/v1/staff
```

**Authentication:** Required.
**Permission:** `staff.view`

**Query Parameters**

| Parameter | Type    | Default | Description |
|-----------|---------|---------|-------------|
| `limit`   | integer | `10`    | Number of records to return |
| `offset`  | integer | `0`     | Number of records to skip |

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
    "limit": 10,
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
|-----------|------|-------------|
| `id`      | UUID | The staff member's internal `id` |

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
    "dateOfBirth": "1990-03-22",
    "dateOfHire": "2024-01-10",
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
|--------|---------|-------|
| `400`  | `Invalid request data` | `:id` is not a valid UUID |
| `500`  | `staff not found: ...`  | No record found for the given ID |

---

### Create Staff

Create a new staff member. The system automatically generates a `staffId` (e.g. `BLN-0002`), a temporary password, and sends a welcome email with login credentials if `hasLogin` is `true`.

```
POST /api/v1/staff
```

**Authentication:** Required.
**Permission:** `staff.create`

**Request Body**

| Field                   | Type            | Required | Description |
|-------------------------|-----------------|----------|-------------|
| `firstName`             | string          | Yes      | First name |
| `lastName`              | string          | Yes      | Last name |
| `email`                 | string          | No       | Contact email — required if `hasLogin: true` |
| `phone`                 | string          | No       | Phone number |
| `address`               | string          | No       | Home address |
| `dateOfBirth`           | string          | No       | ISO 8601 date e.g. `1990-03-22` |
| `dateOfHire`            | string          | Yes      | ISO 8601 date e.g. `2026-01-15` |
| `emergencyContactName`  | string          | No       | Emergency contact full name |
| `emergencyContactPhone` | string          | No       | Emergency contact phone number |
| `bankName`              | string          | No       | Bank name |
| `bankAccountNumber`     | string          | No       | Bank account number |
| `bankAccountName`       | string          | No       | Account name as on bank records |
| `department`            | string          | Yes      | One of: `factory`, `admin`, `sales`, `management` |
| `jobTitle`              | string          | Yes      | Descriptive job title |
| `payType`               | string          | Yes      | `monthly` |
| `baseSalary`            | string (decimal)| Yes      | Monthly gross base salary in NGN |
| `status`                | string          | Yes      | One of: `active`, `inactive`, `fired` |
| `hasLogin`              | boolean         | No       | Default: `false`. Set to `true` to create a Supabase Auth account and send login credentials |

```json
{
  "firstName": "Funmi",
  "lastName": "Adeyemi",
  "email": "funmi@example.com",
  "phone": "+2348055001122",
  "dateOfHire": "2026-07-01",
  "department": "sales",
  "jobTitle": "Sales Executive",
  "payType": "monthly",
  "baseSalary": "120000.00",
  "status": "active",
  "hasLogin": true
}
```

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
|--------|---------|-------|
| `400`  | `Invalid request data` | Malformed or missing required fields |
| `401`  | `Authentication required` | No valid token provided |
| `403`  | `You do not have permission to perform this action` | Missing `staff.create` permission |
| `500`  | `failed to create staff: ...` | Database error or staff_id conflict |
| `500`  | `failed to send welcome email: ...` | Email delivery failed after record was saved |

**Side Effects**

- A unique `staffId` is auto-generated in the format `BLN-NNNN`.
- A random password is generated and hashed before storage.
- If `hasLogin: true`, a Supabase Auth account is created with a synthetic email (`{staffId}@bloansbook.local`).
- A welcome email containing the `staffId` and temporary password is sent to the staff member's `email` address.

---

### Update Staff

Partially update an existing staff record. Only fields included in the request body are updated — all other fields remain unchanged.

```
PATCH /api/v1/staff/:id
```

**Authentication:** Required.
**Permission:** `staff.update`

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | UUID | The staff member's internal `id` |

**Request Body**

All fields are optional. Include only the fields you want to change.

| Field                   | Type            | Description |
|-------------------------|-----------------|-------------|
| `firstName`             | string          | First name |
| `lastName`              | string          | Last name |
| `email`                 | string          | Contact email |
| `phone`                 | string          | Phone number |
| `address`               | string          | Home address |
| `dateOfBirth`           | string          | ISO 8601 date |
| `emergencyContactName`  | string          | Emergency contact name |
| `emergencyContactPhone` | string          | Emergency contact phone |
| `bankName`              | string          | Bank name |
| `bankAccountNumber`     | string          | Bank account number |
| `bankAccountName`       | string          | Account name |
| `department`            | string          | Department |
| `jobTitle`              | string          | Job title |
| `payType`               | string          | Pay type |
| `baseSalary`            | string (decimal)| Base salary |
| `status`                | string          | `active`, `inactive`, or `fired` |

```json
{
  "phone": "+2348077001199",
  "baseSalary": "150000.00",
  "status": "active"
}
```

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
      "email": "funmi@example.com",
      "phone": "+2348077001199",
      "department": "sales",
      "jobTitle": "Sales Executive",
      "status": "active"
    },
    "createdAt": "2026-07-10T11:30:00Z",
    "updatedAt": "2026-07-10T14:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | `:id` is not a valid UUID or body is malformed |
| `401`  | `Authentication required` | No valid token |
| `403`  | `You do not have permission to perform this action` | Missing `staff.update` permission |
| `500`  | `cannot update non-existent staff: ...` | Staff record not found |

---

## Role Management Endpoints

These endpoints manage which roles are assigned to a staff member. All role changes are recorded permanently in `staff_role_history`. All require the `auth.manage_roles` permission.

---

### Assign Role

Assign a role to a staff member. A staff member can hold multiple roles simultaneously.

```
POST /api/v1/staff/:id/roles
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | UUID | The staff member's internal `id` |

**Request Body**

| Field    | Type   | Required | Description |
|----------|--------|----------|-------------|
| `roleId` | UUID   | Yes      | The ID of the role to assign |
| `reason` | string | No       | Optional reason for the audit trail |

```json
{
  "roleId": "b3f2a1d0-...",
  "reason": "Promoted to HR Manager"
}
```

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
|--------|---------|-------|
| `400`  | `Invalid request data` | Invalid UUID or malformed body |
| `500`  | `role already assigned to this staff member` | Role is already held by this staff member |
| `500`  | `staff not found: ...` | Staff record does not exist |

---

### Revoke Role

Remove a role from a staff member.

```
DELETE /api/v1/staff/:id/roles
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | UUID | The staff member's internal `id` |

**Request Body**

| Field    | Type   | Required | Description |
|----------|--------|----------|-------------|
| `roleId` | UUID   | Yes      | The ID of the role to revoke |
| `reason` | string | No       | Optional reason for the audit trail |

```json
{
  "roleId": "b3f2a1d0-...",
  "reason": "Role restructure — moved to admin"
}
```

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

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | Invalid UUID or malformed body |
| `500`  | `role not assigned to this staff member` | The staff member does not hold this role |

---

### Update Role

Atomically swap one role for another in a single transaction. The old role is revoked and the new role is assigned together — either both succeed or neither does. Two history entries are written.

```
PUT /api/v1/staff/:id/roles
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | UUID | The staff member's internal `id` |

**Request Body**

| Field       | Type   | Required | Description |
|-------------|--------|----------|-------------|
| `oldRoleId` | UUID   | Yes      | The role to remove |
| `newRoleId` | UUID   | Yes      | The role to assign |
| `reason`    | string | No       | Reason applies to both the revoke and assign history entries |

```json
{
  "oldRoleId": "b3f2a1d0-...",
  "newRoleId": "c4e3b2a1-...",
  "reason": "Promoted to Production Supervisor"
}
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Role assigned successfully",
  "statusCode": 200,
  "data": {
    "staffId": "018e1c2d-...",
    "roleId": "c4e3b2a1-...",
    "roleName": "production_supervisor",
    "action": "assigned"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | Invalid UUID or malformed body |
| `500`  | `old role not assigned to this staff member` | `oldRoleId` is not currently held |
| `500`  | `failed to assign new role: ...` | `newRoleId` does not exist |

---

### Get Role History

Retrieve the complete role change log for a staff member, ordered newest first.

```
GET /api/v1/staff/:id/roles/history
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | UUID | The staff member's internal `id` |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Role fetched successfully",
  "statusCode": 200,
  "data": [
    {
      "id": "uuid...",
      "roleId": "c4e3b2a1-...",
      "roleName": "production_supervisor",
      "action": "assigned",
      "performedBy": "018e1c2d-...",
      "reason": "Promoted to Production Supervisor",
      "createdAt": "2026-07-10T14:00:00Z"
    },
    {
      "id": "uuid...",
      "roleId": "b3f2a1d0-...",
      "roleName": "hr_manager",
      "action": "revoked",
      "performedBy": "018e1c2d-...",
      "reason": "Promoted to Production Supervisor",
      "createdAt": "2026-07-10T14:00:00Z"
    }
  ]
}
```

---

## Business Rules

- `staffId` is auto-generated at creation in the format `{PREFIX}-NNNN` (e.g. `BLN-0001`). It never changes after creation.
- Staff records are never deleted. Use `status: inactive` to deactivate or `status: fired` to terminate.
- Staff with `has_login: false` have no Supabase Auth account and cannot authenticate.
- If `has_login: true` is set at creation, the welcome email is mandatory — email delivery failure will roll back the operation.
- All role assignment, revocation, and update operations are transactional. The `staff_role_history` table provides a permanent, append-only audit trail of every role change.
- The caller's identity (`performedBy`) is always captured from the JWT — it is never passed in the request body.
