# Roles & Permissions API

**Base path:** `/api/v1/roles`

The roles module manages the access control layer of the system. Roles are named groups of permissions. Every login-enabled staff member holds one or more roles. The middleware enforces permissions on every protected endpoint by checking the caller's role-derived permission set on each request.

All endpoints in this module require authentication and the `auth.manage_roles` permission. In practice, only `super_admin` holds this permission by default.

---

## How RBAC Works

```
staff member → staff_roles → roles → role_permissions → permissions
```

1. A **permission** is a single action code such as `invoices.post` or `staff.create`.
2. A **role** is a named group of permissions e.g. `hr_manager`, `sales_manager`.
3. A **staff_roles** entry links a staff member to a role.
4. On each authenticated request, the middleware loads the caller's full permission set (derived through their roles) into the request context.
5. Each route declares the specific permission it requires. The `RequirePermission` guard checks the context and blocks the request if the permission is absent.

---

## Data Models

### Role Object

Returned by `GET /roles/:id` and `GET /roles`.

| Field         | Type               | Description |
|---------------|--------------------|-------------|
| `id`          | UUID               | Internal identifier |
| `name`        | string             | Machine-readable role name e.g. `super_admin` |
| `permissions` | PermissionSummary[]| All permissions assigned to this role |

### Permission Summary Object

| Field    | Type   | Description |
|----------|--------|-------------|
| `id`     | UUID   | Internal identifier |
| `code`   | string | Permission code e.g. `invoices.post` |
| `module` | string | Module the permission belongs to e.g. `invoices` |

### Roles Object (raw — returned by create)

| Field         | Type    | Description |
|---------------|---------|-------------|
| `id`          | UUID    | Internal identifier |
| `name`        | string  | Role name |
| `description` | string  | Human-readable description |
| `isSystem`    | boolean | If `true`, this role cannot be deleted |
| `createdAt`   | string  | ISO 8601 timestamp |
| `updatedAt`   | string  | ISO 8601 timestamp |

---

## Seeded Roles

The following 7 system roles are seeded at migration time. All have `isSystem: true` and cannot be deleted via the API.

| Role Name               | Description |
|-------------------------|-------------|
| `super_admin`           | Full system access. Only one active super_admin at a time |
| `hr_manager`            | Manages staff records, payroll, and attendance |
| `admin`                 | General administration — customers, suppliers, invoices, bills |
| `sales_manager`         | Manages sales orders, invoices, and customer relationships |
| `production_supervisor` | Oversees production jobs, materials, and attendance |
| `inventory_officer`     | Manages inventory transactions and purchase orders |
| `auditor`               | Read-only access to audit log and reports |

---

## Endpoints

### List All Roles

Returns all roles with their assigned permissions.

```
GET /api/v1/roles
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Query Parameters**

| Parameter | Type    | Default | Description |
|-----------|---------|---------|-------------|
| `limit`   | integer | `10`    | Number of roles to return |
| `offset`  | integer | `0`     | Number of roles to skip |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Roles fetched successfully",
  "statusCode": 200,
  "data": [
    {
      "id": "b3f2a1d0-...",
      "name": "hr_manager",
      "permissions": [
        { "id": "uuid...", "code": "staff.create", "module": "staff" },
        { "id": "uuid...", "code": "staff.view", "module": "staff" },
        { "id": "uuid...", "code": "payroll.view", "module": "payroll" }
      ]
    }
  ]
}
```

---

### Get Role by ID

Retrieve a single role with all of its assigned permissions.

```
GET /api/v1/roles/:id
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Path Parameters**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id`      | UUID | Role's internal `id` |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Role fetched successfully",
  "statusCode": 200,
  "data": {
    "id": "b3f2a1d0-...",
    "name": "hr_manager",
    "permissions": [
      { "id": "uuid...", "code": "staff.create", "module": "staff" },
      { "id": "uuid...", "code": "staff.view", "module": "staff" },
      { "id": "uuid...", "code": "staff.update", "module": "staff" },
      { "id": "uuid...", "code": "payroll.view", "module": "payroll" }
    ]
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | `:id` is not a valid UUID |
| `500`  | `role not found: ...` | No role found for the given ID |

---

### Create Role

Create a new custom role. System roles are seeded at migration time — this endpoint is used to create additional business-specific roles.

```
POST /api/v1/roles
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Request Body**

| Field         | Type    | Required | Description |
|---------------|---------|----------|-------------|
| `name`        | string  | Yes      | Unique, machine-readable role name (lowercase, underscores) |
| `isSystem`    | boolean | Yes      | Set to `false` for custom roles |
| `description` | string  | No       | Human-readable description |

```json
{
  "name": "logistics_coordinator",
  "isSystem": false,
  "description": "Manages supplier deliveries and logistics coordination"
}
```

**Success Response — `201 Created`**

```json
{
  "status": "success",
  "message": "Role created successfully",
  "statusCode": 201,
  "data": {
    "id": "d5e4c3b2-...",
    "name": "logistics_coordinator",
    "description": "Manages supplier deliveries and logistics coordination",
    "isSystem": false,
    "createdAt": "2026-07-10T12:00:00Z",
    "updatedAt": "2026-07-10T12:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | Malformed body |
| `400`  | `role name is required` | `name` field is empty |
| `500`  | `failed to create role: ...` | Database error or name conflict (names are unique) |

---

### Create Permission

Create a new permission code. Permissions are seeded at migration time for all standard modules — this endpoint is used when adding permissions for new custom features.

```
POST /api/v1/roles/permissions
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Request Body**

| Field         | Type   | Required | Description |
|---------------|--------|----------|-------------|
| `code`        | string | Yes      | Unique permission code in `module.action` format e.g. `logistics.manage` |
| `module`      | string | Yes      | The module this permission belongs to |
| `description` | string | No       | Human-readable description of what this permission allows |

```json
{
  "code": "logistics.manage",
  "module": "logistics",
  "description": "Manage logistics and delivery scheduling"
}
```

**Success Response — `201 Created`**

```json
{
  "status": "success",
  "message": "Permission created successfully",
  "statusCode": 201,
  "data": {
    "id": "e6f5d4c3-...",
    "code": "logistics.manage",
    "module": "logistics",
    "description": "Manage logistics and delivery scheduling",
    "createdAt": "2026-07-10T12:05:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | Malformed body |
| `400`  | `permission code and module are required` | `code` or `module` is empty |
| `500`  | `failed to create permission: ...` | Database error or code conflict (codes are unique) |

---

### Assign Permission to Role

Link a permission to a role. After this call, any staff member holding that role immediately gains the permission on their next request.

```
POST /api/v1/roles/assign-permission
```

**Authentication:** Required.
**Permission:** `auth.manage_roles`

**Request Body**

| Field          | Type | Required | Description |
|----------------|------|----------|-------------|
| `roleId`       | UUID | Yes      | The role to assign the permission to |
| `permissionId` | UUID | Yes      | The permission to assign |

```json
{
  "roleId": "d5e4c3b2-...",
  "permissionId": "e6f5d4c3-..."
}
```

**Success Response — `201 Created`**

```json
{
  "status": "success",
  "message": "Permission assigned to role successfully",
  "statusCode": 201,
  "data": {
    "created_at": "2026-07-10T12:10:00Z",
    "role_id": "d5e4c3b2-...",
    "permission_id": "e6f5d4c3-...",
    "role_name": "logistics_coordinator",
    "permission_code": "logistics.manage",
    "permission_module": "logistics"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|--------|---------|-------|
| `400`  | `Invalid request data` | Malformed body or invalid UUIDs |
| `500`  | `failed to validate role: ...` | `roleId` does not exist |
| `500`  | `failed to assign permission to role: ...` | Database error |

---

## Seeded Permissions

The following permissions are seeded at migration time (`000015_seed_permissions`). The `super_admin` role is pre-assigned all of them. Other roles start with no permissions and are configured by Super Admin via this API.

| Module            | Permission Codes |
|-------------------|-----------------|
| `auth`            | `auth.manage_accounts`, `auth.reset_password`, `auth.manage_roles` |
| `staff`           | `staff.create`, `staff.view`, `staff.update`, `staff.deactivate`, `staff.terminate` |
| `customers`       | `customers.create`, `customers.view`, `customers.update` |
| `suppliers`       | `suppliers.create`, `suppliers.view`, `suppliers.update` |
| `products`        | `products.create`, `products.view`, `products.update` |
| `materials`       | `materials.create`, `materials.view`, `materials.update` |
| `inventory`       | `inventory.view`, `inventory.adjust`, `inventory.approve_adjustment`, `inventory.issue_to_job` |
| `purchase_orders` | `purchase_orders.create`, `purchase_orders.view`, `purchase_orders.approve`, `purchase_orders.deliver`, `purchase_orders.cancel` |
| `bills`           | `bills.create`, `bills.view`, `bills.approve`, `bills.reverse` |
| `invoices`        | `invoices.create`, `invoices.view`, `invoices.post`, `invoices.reverse`, `invoices.cancel`, `invoices.apply_discount` |
| `payments`        | `payments.create`, `payments.view` |
| `job_costing`     | `job_costing.view`, `job_costing.add_labor`, `job_costing.add_overhead` |
| `payroll`         | `payroll.create_run`, `payroll.view`, `payroll.approve`, `payroll.mark_paid` |
| `attendance`      | `attendance.import`, `attendance.view`, `attendance.approve_exception`, `attendance.approve_overtime` |
| `tasks`           | `tasks.create`, `tasks.view`, `tasks.update`, `tasks.assign` |
| `errands`         | `errands.create`, `errands.view`, `errands.approve` |
| `reports`         | `reports.view`, `reports.export` |
| `audit`           | `audit.view` |
| `settings`        | `settings.view`, `settings.update` |
| `periods`         | `periods.close`, `periods.reopen` |

---

## Business Rules

- Role names must be unique across the system.
- Permission codes must be unique. The recommended format is `module.action`.
- The seven system roles (`is_system: true`) cannot be deleted via the API. Attempting to delete them is blocked at the database level.
- Assigning a permission to a role takes effect immediately — the permission is resolved fresh on each authenticated request, so there is no need to re-login.
- `super_admin` is pre-assigned all 68 seeded permissions. Removing a permission from `super_admin` is technically possible via the database but not recommended.
- Role assignment and revocation from individual staff members is managed through the [Staff API](./staff.md).
