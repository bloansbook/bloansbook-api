# BloansBooks — Entity Relationship Document

**Version:** 1.0  
**Last Updated:** March 2026  
**Status:** In Progress

---

## How to Read This Document

Each group covers a set of related tables. For every table you will find:

- **Purpose** — what the table is for
- **Columns** — name, type, constraints, and notes
- **Relationships** — how it connects to other tables

**Column type conventions used throughout:**
| Type | Meaning |
|---|---|
| `UUID` | Auto-generated universally unique identifier via `gen_random_uuid()` |
| `TEXT` | Variable length string, no fixed limit |
| `BOOLEAN` | true / false |
| `DATE` | Calendar date, no time component |
| `TIMESTAMPTZ` | Timestamp with timezone — always stored in UTC |
| `NUMERIC(15,2)` | Exact decimal number — used for all money values |

---

## AUTH & SECURITY

**Tables:** `roles`, `permissions`, `role_permissions`, `fired_staff`

These tables form the access control foundation. They define what roles
exist, what each role can do, and which staff members are permanently
blocked from the system.

---

### Table: `roles`

**Purpose:**
Stores the named roles in the system. Each role groups a set of
permissions. Some roles are system-defined and cannot be deleted.
The Super Admin role is the only role that cannot be assigned or
revoked by any other user.

| Column        | Type        | Constraints                            | Notes                                                                                                                                           |
| ------------- | ----------- | -------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier — never exposed to users                                                                                                    |
| `name`        | TEXT        | NOT NULL, UNIQUE                       | Machine-readable role name. e.g. `super_admin`, `hr_manager`, `admin`, `sales_manager`, `production_supervisor`, `inventory_officer`, `auditor` |
| `description` | TEXT        | NOT NULL                               | Human-readable description of the role                                                                                                          |
| `is_system`   | BOOLEAN     | NOT NULL, DEFAULT false                | If true, this role cannot be deleted. Set to true for all 7 named roles at seed time                                                            |
| `created_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go backend, UTC                                                                                                                          |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `name`

**Business rules:**

- All 7 named roles (`super_admin`, `hr_manager`, `admin`, `sales_manager`, `production_supervisor`, `inventory_officer`, `auditor`) are seeded at migration time with `is_system = true`
- `is_system = true` roles cannot be deleted via the API
- There can only be one active Super Admin at a time — enforced by a unique partial index on the `staff` table (see Group 2) and by Go middleware

---

### Table: `permissions`

**Purpose:**
Stores every action that can be performed in the system. Each permission
has a unique code that the Go middleware checks against before processing
any request. Permissions are seeded at migration time and managed by
Super Admin through the UI.

| Column        | Type        | Constraints                            | Notes                                                                                                                              |
| ------------- | ----------- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                                |
| `code`        | TEXT        | NOT NULL, UNIQUE                       | The permission code checked by Go middleware. Format: `module.action` e.g. `invoice.create`, `payroll.approve`, `inventory.adjust` |
| `description` | TEXT        | NOT NULL                               | Human-readable description e.g. "Create and draft invoices"                                                                        |
| `module`      | TEXT        | NOT NULL                               | The system module this permission belongs to. e.g. `invoices`, `payroll`, `inventory`, `purchase_orders`, `staff`, `reports`       |
| `created_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go backend, UTC                                                                                                             |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `code`

**Seeded permission codes by module:**

| Module            | Permission Codes                                                                                                                 |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `auth`            | `auth.manage_accounts`, `auth.reset_password`, `auth.manage_roles`                                                               |
| `staff`           | `staff.create`, `staff.view`, `staff.update`, `staff.deactivate`, `staff.terminate`                                              |
| `customers`       | `customers.create`, `customers.view`, `customers.update`                                                                         |
| `suppliers`       | `suppliers.create`, `suppliers.view`, `suppliers.update`                                                                         |
| `products`        | `products.create`, `products.view`, `products.update`                                                                            |
| `materials`       | `materials.create`, `materials.view`, `materials.update`                                                                         |
| `inventory`       | `inventory.view`, `inventory.adjust`, `inventory.approve_adjustment`, `inventory.issue_to_job`                                   |
| `purchase_orders` | `purchase_orders.create`, `purchase_orders.view`, `purchase_orders.approve`, `purchase_orders.deliver`, `purchase_orders.cancel` |
| `bills`           | `bills.create`, `bills.view`, `bills.approve`, `bills.reverse`                                                                   |
| `invoices`        | `invoices.create`, `invoices.view`, `invoices.post`, `invoices.reverse`, `invoices.cancel`, `invoices.apply_discount`            |
| `payments`        | `payments.create`, `payments.view`                                                                                               |
| `job_costing`     | `job_costing.view`, `job_costing.add_labor`, `job_costing.add_overhead`                                                          |
| `payroll`         | `payroll.create_run`, `payroll.view`, `payroll.approve`, `payroll.mark_paid`                                                     |
| `attendance`      | `attendance.import`, `attendance.view`, `attendance.approve_exception`, `attendance.approve_overtime`                            |
| `tasks`           | `tasks.create`, `tasks.view`, `tasks.update`, `tasks.assign`                                                                     |
| `errands`         | `errands.create`, `errands.view`, `errands.approve`                                                                              |
| `reports`         | `reports.view`, `reports.export`                                                                                                 |
| `audit`           | `audit.view`                                                                                                                     |
| `settings`        | `settings.view`, `settings.update`                                                                                               |
| `periods`         | `periods.close`, `periods.reopen`                                                                                                |

---

### Table: `role_permissions`

**Purpose:**
Junction table that maps permissions to roles. Defines exactly what
each role can do. Managed by Super Admin through the UI after the
initial seed. Super Admin has all permissions by default.

| Column          | Type        | Constraints                                       | Notes                                          |
| --------------- | ----------- | ------------------------------------------------- | ---------------------------------------------- |
| `role_id`       | UUID        | NOT NULL, FK → `roles.id` ON DELETE CASCADE       | The role being assigned permissions            |
| `permission_id` | UUID        | NOT NULL, FK → `permissions.id` ON DELETE CASCADE | The permission being assigned                  |
| `created_at`    | TIMESTAMPTZ | NOT NULL, DEFAULT now()                           | When this permission was assigned to this role |

**Indexes:**

- `PRIMARY KEY (role_id, permission_id)` — composite, prevents duplicate assignments

**Relationships:**

```
roles (1) ──< role_permissions >── (many) permissions
```

**Business rules:**

- Super Admin role has all permissions assigned at seed time
- Adding or removing permissions from a role takes effect immediately — the next request by a user with that role will reflect the change since permissions are read from the JWT on every request
- Deleting a role cascades and removes all its role_permission entries

---

### Table: `fired_staff`

**Purpose:**
Permanent termination records. The Go backend checks this table on
every login attempt before validating the password. If the Staff ID
is found here and `is_overridden = false`, login is refused
immediately regardless of password correctness.

Records in this table are never deleted. An incorrect termination
is corrected by setting `is_overridden = true` with a mandatory reason.

| Column             | Type        | Constraints                            | Notes                                                                                                                                                                            |
| ------------------ | ----------- | -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`               | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                                                                              |
| `staff_id`         | TEXT        | NOT NULL                               | The terminated staff member's business ID e.g. `BLN-0042`. Loose reference to `staff.staff_id` — deliberately no FK so the block persists even if staff records are restructured |
| `termination_date` | DATE        | NOT NULL                               | Date the staff member was terminated                                                                                                                                             |
| `reason`           | TEXT        | NOT NULL                               | Mandatory termination reason e.g. Resigned, Dismissed, Contract Ended, Redundancy                                                                                                |
| `recorded_by`      | TEXT        | NOT NULL                               | Staff ID of the HR Manager or Super Admin who recorded the termination                                                                                                           |
| `recorded_at`      | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go backend, UTC                                                                                                                                                           |
| `is_overridden`    | BOOLEAN     | NOT NULL, DEFAULT false                | Set to true if the termination was recorded in error and a Super Admin has overridden it                                                                                         |
| `overridden_by`    | TEXT        | NULL                                   | Staff ID of the Super Admin who overrode the termination. NULL until overridden                                                                                                  |
| `overridden_at`    | TIMESTAMPTZ | NULL                                   | Timestamp of the override. Set by Go, UTC. NULL until overridden                                                                                                                 |
| `override_reason`  | TEXT        | NULL                                   | Mandatory when `is_overridden` is set to true. Must be provided by the Super Admin                                                                                               |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `staff_id` — this column is queried on every login attempt, must be fast

**Business rules:**

- Login check: `SELECT * FROM fired_staff WHERE staff_id = $1 AND is_overridden = false`
- If a row is returned, login is refused. No password check occurs.
- Only Super Admin can set `is_overridden = true`
- Super Admin's own Staff ID cannot be added to this table by HR Manager
- All termination events and overrides are written to the `audit_log`

---

## STAFF

**Tables:** `staff`, `user_roles`, `role_history`

The Staff Master Record is the single source of truth for every person
in the business — from the owner down to tailors and assistants. Every
module that references a person references their `staff_id`. Records
are never deleted.

---

### Table: `staff`

**Purpose:**
The single source of truth for every person employed by the business.
Created when a staff member is hired. Never deleted — status is set
to `inactive` when a staff member leaves. The `staff_id` field is
the login username for login-enabled staff and the identifier used
across every module: payroll, tasks, attendance, job costing, and
audit logs.

| Column                    | Type          | Constraints                            | Notes                                                                                                                                                                     |
| ------------------------- | ------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`                      | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier — used for all FK relationships across the system                                                                                                     |
| `staff_id`                | TEXT          | NOT NULL, UNIQUE                       | Auto-generated business ID e.g. `BLN-0001`. Format: `BLN-` followed by 4-digit zero-padded number. This is the login username for `has_login = true` staff. Never changes |
| `first_name`              | TEXT          | NOT NULL                               |                                                                                                                                                                           |
| `last_name`               | TEXT          | NOT NULL                               |                                                                                                                                                                           |
| `phone`                   | TEXT          | NOT NULL                               | Primary contact number                                                                                                                                                    |
| `email`                   | TEXT          | NULL                                   | Contact email only — never used for login                                                                                                                                 |
| `address`                 | TEXT          | NULL                                   | Home address                                                                                                                                                              |
| `date_of_birth`           | DATE          | NULL                                   |                                                                                                                                                                           |
| `date_of_hire`            | DATE          | NOT NULL                               | Date the staff member joined the company                                                                                                                                  |
| `emergency_contact_name`  | TEXT          | NULL                                   |                                                                                                                                                                           |
| `emergency_contact_phone` | TEXT          | NULL                                   |                                                                                                                                                                           |
| `bank_name`               | TEXT          | NULL                                   | For payroll disbursement                                                                                                                                                  |
| `bank_account_number`     | TEXT          | NULL                                   |                                                                                                                                                                           |
| `bank_account_name`       | TEXT          | NULL                                   | Account name as it appears on the bank record                                                                                                                             |
| `department`              | TEXT          | NOT NULL                               | One of: `factory`, `admin`, `sales`, `management`                                                                                                                         |
| `job_title`               | TEXT          | NOT NULL                               | Descriptive title only e.g. Tailor, Sales Manager, Production Supervisor. Not the system role                                                                             |
| `pay_type`                | TEXT          | NOT NULL, DEFAULT 'monthly'            | v1 supports `monthly` only                                                                                                                                                |
| `base_salary`             | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | Monthly gross base salary in NGN. Use NUMERIC not FLOAT — exact decimal required for payroll                                                                              |
| `status`                  | TEXT          | NOT NULL, DEFAULT 'active'             | One of: `active`, `inactive`. Inactive staff are excluded from payroll runs and task assignment                                                                           |
| `has_login`               | BOOLEAN       | NOT NULL, DEFAULT false                | `true` for the 7 named roles. `false` for no-login staff categories (Tailors, Production Assistants, Social Media Managers, etc.)                                         |
| `supabase_uid`            | UUID          | NULL, UNIQUE                           | The Supabase Auth user ID returned when a login account is created. NULL for `has_login = false` staff. Used to link incoming JWTs back to the correct staff record       |
| `created_by`              | TEXT          | NOT NULL                               | Staff ID of the HR Manager or Super Admin who created this record                                                                                                         |
| `created_at`              | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go backend, UTC                                                                                                                                                    |
| `updated_at`              | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Updated by Go on every record change, UTC                                                                                                                                 |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `staff_id`
- `UNIQUE` on `supabase_uid` (partial — where not null)
- Index on `status` — frequently filtered in payroll and task queries
- Index on `department` — frequently filtered in task and dashboard queries
- `UNIQUE PARTIAL INDEX` on `staff_id` where role = `super_admin` and `status = active` — enforces one active Super Admin at DB level

**Business rules:**

- `staff_id` is generated by Go at record creation. Format: `BLN-` + zero-padded 4-digit sequential number
- For `has_login = true` staff, Go creates a Supabase Auth account after saving the staff record and stores the returned `supabase_uid`
- For `has_login = false` staff, no Supabase Auth account is created and `supabase_uid` remains NULL
- Setting `status = inactive` does not delete the Supabase Auth account — it is disabled via the Supabase Auth admin API
- `base_salary` of 0 is valid for staff whose compensation is managed outside the system

---

### Table: `user_roles`

**Purpose:**
Stores the current active role for each login-enabled staff member.
One row per staff member — updated when a role changes. Before
updating, Go writes the previous role to `role_history`.

Only exists for `has_login = true` staff. No-login staff have no
system role and no entry in this table.

| Column        | Type        | Constraints                             | Notes                                                            |
| ------------- | ----------- | --------------------------------------- | ---------------------------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid()  | Internal identifier                                              |
| `staff_id`    | TEXT        | NOT NULL, UNIQUE, FK → `staff.staff_id` | UNIQUE enforces one active role per staff member at the DB level |
| `role_id`     | UUID        | NOT NULL, FK → `roles.id`               | The currently active role                                        |
| `assigned_by` | TEXT        | NOT NULL                                | Staff ID of the HR Manager or Super Admin who assigned this role |
| `assigned_at` | TIMESTAMPTZ | NOT NULL, DEFAULT now()                 | Set by Go, UTC                                                   |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `staff_id` — one active role per staff member enforced at DB level

**Relationships:**

```
staff (1) ──── user_roles (1) ──── roles (1)
```

**Business rules:**

- When a role changes, Go must:
  1. Read the current `user_roles` row for the staff member
  2. Write the old role to `role_history` with `unassigned_by` and `unassigned_at`
  3. Update `user_roles` with the new role
  4. All three steps run inside a single DB transaction
- HR Manager can assign any role except `super_admin`
- Only Super Admin can assign the `super_admin` role (ownership transfer)
- The `super_admin` role can only be held by one active staff member at a time

---

### Table: `role_history`

**Purpose:**
Append-only log of every role change for every staff member.
A new row is written each time a role is changed — never updated,
never deleted. When a staff member is fetched, their role history
is joined and returned alongside their current role.

| Column          | Type        | Constraints                            | Notes                                                                                   |
| --------------- | ----------- | -------------------------------------- | --------------------------------------------------------------------------------------- |
| `id`            | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                     |
| `staff_id`      | TEXT        | NOT NULL, FK → `staff.staff_id`        | The staff member whose role changed                                                     |
| `role_id`       | UUID        | NOT NULL, FK → `roles.id`              | The role that was held during this period                                               |
| `assigned_by`   | TEXT        | NOT NULL                               | Staff ID of who assigned this role                                                      |
| `assigned_at`   | TIMESTAMPTZ | NOT NULL                               | When this role was assigned. Copied from `user_roles.assigned_at` at the time of change |
| `unassigned_by` | TEXT        | NOT NULL                               | Staff ID of who changed the role away from this one                                     |
| `unassigned_at` | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | When this role was removed. Set by Go, UTC                                              |
| `reason`        | TEXT        | NULL                                   | Optional reason for the role change e.g. "Promoted to Production Supervisor"            |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `staff_id` — queried every time a staff member is fetched

**Relationships:**

```
staff (1) ──< role_history >── roles (1)
```

**Business rules:**

- Never updated after insert — append only
- Written by Go inside the same transaction as the `user_roles` update
- Returned as an array when a staff member record is fetched

---

## Full Relationship Map

```
roles (1) ──< role_permissions >── (many) permissions

staff (1) ──── user_roles (1) ──── roles (1)
staff (1) ──< role_history >── roles (1)

fired_staff.staff_id ........> staff.staff_id
(loose reference — no FK enforced, intentional)

staff.supabase_uid ──── Supabase Auth (external system)
```

---

## Notes

- `....>` denotes a loose reference (no FK constraint)
- `────` denotes a direct FK relationship
- `──<` denotes a one-to-many relationship
- `>──<` denotes a many-to-many junction

---

_BloansBooks ERD · Internal Technical Documentation · Confidential_
