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

## GROUP 1 - AUTH & SECURITY

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

## GROUP 2 - STAFF

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

---

## GROUP 3 - MASTER DATA

**Tables:** `customers`, `suppliers`, `products`, `product_variants`, `materials`

These are the core reference records that all transactional modules point to.
They are created and maintained independently of any transaction. Nothing
financial happens here — these are purely descriptive master records that
other modules reference outward.

---

### Table: `customers`

**Purpose:**
Every person or business that places an order or receives an invoice.
Created by Admin or Sales Manager. Customer information can only be edited
by Admin or Super Admin — Sales Manager can create but not edit.

| Column        | Type        | Constraints                            | Notes                                         |
| ------------- | ----------- | -------------------------------------- | --------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                           |
| `customer_id` | TEXT        | NOT NULL, UNIQUE                       | Business ID e.g. `CUST-0001`. Generated by Go |
| `name`        | TEXT        | NOT NULL                               | Full business or individual name              |
| `phone`       | TEXT        | NOT NULL                               | Primary contact number                        |
| `email`       | TEXT        | NULL                                   | Contact email                                 |
| `address`     | TEXT        | NULL                                   | Physical or postal address                    |
| `type`        | TEXT        | NOT NULL                               | One of: `retail`, `corporate`                 |
| `currency`    | TEXT        | NOT NULL, DEFAULT 'NGN'                | Billing currency. NGN only for v1             |
| `notes`       | TEXT        | NULL                                   | Free-text account notes                       |
| `created_by`  | TEXT        | NOT NULL                               | Staff ID of creator                           |
| `created_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                |
| `updated_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `customer_id`
- Index on `name` — frequently searched

**Business rules:**

- `customer_id` format: `CUST-` + 4-digit zero-padded sequential number
- Sales Manager can create but cannot edit after creation
- Admin and Super Admin can create and edit
- Customers are never deleted — relationship notes go in the `notes` field

---

### Table: `suppliers`

**Purpose:**
Every business or individual the company procures from. Referenced by
Purchase Orders, Bills, and Inventory transactions. Created and managed
by Admin or Inventory Officer.

| Column        | Type        | Constraints                            | Notes                                                                              |
| ------------- | ----------- | -------------------------------------- | ---------------------------------------------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                |
| `supplier_id` | TEXT        | NOT NULL, UNIQUE                       | Business ID e.g. `SPL-0001`. Generated by Go                                       |
| `name`        | TEXT        | NOT NULL                               | Supplier business name                                                             |
| `phone`       | TEXT        | NOT NULL                               | Primary contact number                                                             |
| `email`       | TEXT        | NULL                                   | Contact email — used when sending PO PDFs via Resend                               |
| `address`     | TEXT        | NULL                                   | Registered or delivery address                                                     |
| `currency`    | TEXT        | NOT NULL, DEFAULT 'NGN'                | Billing currency.                                                                  |
| `category`    | TEXT        | NOT NULL                               | One of: `raw_materials`, `printing`, `logistics`, `artisans`, `utilities`, `other` |
| `notes`       | TEXT        | NULL                                   | Free-text notes                                                                    |
| `created_by`  | TEXT        | NOT NULL                               | Staff ID of creator                                                                |
| `created_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                     |
| `updated_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                                     |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `supplier_id`
- Index on `name` — frequently searched
- Index on `category` — filtered in supplier reports

**Business rules:**

- `supplier_id` format: `SPL-` + 4-digit zero-padded sequential number
- `email` should be populated for suppliers who receive POs electronically
- Suppliers are never deleted

---

### Table: `products`

**Purpose:**
The base product record. Stores attributes shared across all variants of a
product. A product on its own is not directly orderable — orders and
inventory transactions always reference a specific variant from
`product_variants`. Every product must have at least one variant.

| Column        | Type        | Constraints                            | Notes                                        |
| ------------- | ----------- | -------------------------------------- | -------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                          |
| `product_id`  | TEXT        | NOT NULL, UNIQUE                       | Business ID e.g. `PRD-0001`. Generated by Go |
| `name`        | TEXT        | NOT NULL                               | Product name e.g. Jute Bag, Tote Bag         |
| `description` | TEXT        | NULL                                   | Full product description                     |
| `status`      | TEXT        | NOT NULL, DEFAULT 'active'             | One of: `active`, `inactive`                 |
| `created_by`  | TEXT        | NOT NULL                               | Staff ID of creator                          |
| `created_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                               |
| `updated_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC               |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `product_id`
- Index on `status`

**Business rules:**

- `product_id` format: `PRD-` + 4-digit zero-padded sequential number
- When a product is created with no explicit variants, Go automatically
  creates one default variant in `product_variants`
- Products are never deleted — set `status = inactive` to retire

---

### Table: `product_variants`

**Purpose:**
Each row represents a specific orderable variant of a product. Variants
differ by size, colour, or any other attribute defined in the `attributes`
JSONB field. Sales Orders, Invoices, and Inventory transactions always
reference a variant, never the base product directly.

| Column                | Type          | Constraints                                    | Notes                                                                                                                                      |
| --------------------- | ------------- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `id`                  | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid()         | Internal identifier                                                                                                                        |
| `variant_id`          | TEXT          | NOT NULL, UNIQUE                               | Business ID e.g. `VAR-0001`. Generated by Go                                                                                               |
| `product_id`          | UUID          | NOT NULL, FK → `products.id` ON DELETE CASCADE | Parent product                                                                                                                             |
| `sku`                 | TEXT          | NOT NULL, UNIQUE                               | Auto-generated e.g. `SKU-000001`. Generated by Go                                                                                          |
| `size`                | TEXT          | NULL                                           | e.g. Small, Medium, Large, A4, A3                                                                                                          |
| `color`               | TEXT          | NULL                                           | e.g. Red, Navy Blue, Natural                                                                                                               |
| `attributes`          | JSONB         | NULL                                           | Any additional variant dimensions e.g. `{"material": "canvas", "handle": "rope"}`. Allows future variant attributes without schema changes |
| `selling_price`       | NUMERIC(15,2) | NOT NULL, DEFAULT 0                            | Default unit selling price in NGN. Can be overridden at order level                                                                        |
| `current_quantity`    | NUMERIC(15,2) | NOT NULL, DEFAULT 0                            | **Cache.** Updated after every `PRODUCTION-IN` and `SALE-OUT` transaction. Never edited directly                                           |
| `current_wac`         | NUMERIC(15,2) | NOT NULL, DEFAULT 0                            | **Cache.** Current weighted average cost per unit. Never edited directly                                                                   |
| `current_stock_value` | NUMERIC(15,2) | NOT NULL, DEFAULT 0                            | **Cache.** Computed: `current_quantity × current_wac`. Used by dashboard. Never edited directly                                            |
| `status`              | TEXT          | NOT NULL, DEFAULT 'active'                     | One of: `active`, `inactive`                                                                                                               |
| `created_by`          | TEXT          | NOT NULL                                       | Staff ID of creator                                                                                                                        |
| `created_at`          | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                        | Set by Go, UTC                                                                                                                             |
| `updated_at`          | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                        | Set by Go on every update, UTC                                                                                                             |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `variant_id`
- `UNIQUE` on `sku`
- Index on `product_id`
- Index on `status`

**Business rules:**

- `variant_id` format: `VAR-` + 4-digit zero-padded sequential number
- `sku` format: `SKU-` + 6-digit zero-padded sequential number
- All three cache columns update inside the same DB transaction as the
  inventory transaction that caused the change — they cannot get out of sync
- Cache columns are read-only from the API — no endpoint allows direct editing

---

### Table: `materials`

**Purpose:**
Raw materials catalogue. Every input material used in production.
Referenced by Inventory transactions of type `PURCHASE`, `ISSUE-TO-JOB`,
`ADJUSTMENT-IN`, and `ADJUSTMENT-OUT`. The `reorder_level` field drives
low stock alerts. Cache columns keep dashboard reads fast without scanning
the full inventory ledger on every request.

| Column                | Type          | Constraints                            | Notes                                                                                                                         |
| --------------------- | ------------- | -------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `id`                  | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                           |
| `material_id`         | TEXT          | NOT NULL, UNIQUE                       | Business ID e.g. `MAT-0001`. Generated by Go                                                                                  |
| `name`                | TEXT          | NOT NULL                               | Material name e.g. Cotton Fabric, Ink — Cyan                                                                                  |
| `unit_of_measure`     | TEXT          | NOT NULL                               | e.g. `metres`, `litres`, `kg`, `pcs`, `rolls`                                                                                 |
| `reorder_level`       | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | When `current_quantity` falls to or below this value, a low stock alert is triggered                                          |
| `current_quantity`    | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | **Cache.** Updated after every inventory transaction affecting this material. Never edited directly                           |
| `current_wac`         | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | **Cache.** Current weighted average cost per unit. Updated after every approved `PURCHASE` transaction. Never edited directly |
| `current_stock_value` | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | **Cache.** Computed: `current_quantity × current_wac`. Used by dashboard inventory value tile. Never edited directly          |
| `status`              | TEXT          | NOT NULL, DEFAULT 'active'             | One of: `active`, `inactive`                                                                                                  |
| `created_by`          | TEXT          | NOT NULL                               | Staff ID of creator                                                                                                           |
| `created_at`          | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                                                                |
| `updated_at`          | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                                                                                |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `material_id`
- Index on `status`
- Index on `current_quantity` — used in low stock alert queries

**Business rules:**

- `material_id` format: `MAT-` + 4-digit zero-padded sequential number
- All three cache columns update inside the same DB transaction as the
  inventory transaction that caused the change
- Low stock check: after every transaction decreasing `current_quantity`,
  Go checks `current_quantity <= reorder_level`. If true, alert is triggered
- Materials are never deleted — set `status = inactive` to retire

---

## Relationships in Group 3

```
products (1) ──< product_variants (many)

product_variants ──< inventory_transactions (Group 4)
product_variants ──< order_lines (Group 7)

customers ──< orders (Group 7)
customers ──< invoices (Group 7)
customers ──< customer_payments (Group 7)

suppliers ──< purchase_orders (Group 5)
suppliers ──< bills (Group 6)
suppliers ──< inventory_transactions (Group 4)

materials ──< inventory_transactions (Group 4)
materials ──< purchase_order_lines (Group 5)
```

Group 3 has no internal relationships between its own tables — all five
tables are independent master records that downstream modules reference
outward.

---

---

## GROUP 4 - INVENTORY

**Tables:** `inventory_transactions`, `wac_snapshots`

The inventory system is strictly ledger-based. Stock levels never change
by direct editing — they change only as the result of validated inventory
transactions. The ledger is the single source of truth for all stock
movements and valuations. Records are never updated or deleted —
corrections are made through reversal entries.

---

### Table: `inventory_transactions`

**Purpose:**
Every stock movement in the system creates a row here. Purchases,
issues to jobs, production completions, adjustments, sales, and
reversals are all recorded as individual rows. This table is
append-only. The current stock level of any item is always derived
by summing this ledger filtered to active transactions.

| Column               | Type          | Constraints                            | Notes                                                                                                                                                                                                                                                                                                               |
| -------------------- | ------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`                 | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                                                                                                                                                                                                                 |
| `txn_id`             | TEXT          | NOT NULL, UNIQUE                       | Business ID. Format: `TXN/YYYY/MM/DD/001`. Generated by Go. Daily format chosen because transaction volume is unbounded over time and daily tracking aids reconciliation                                                                                                                                            |
| `idempotency_key`    | TEXT          | UNIQUE, NULL                           | Caller-supplied or Go-constructed key. If a transaction with this key already exists, Go returns the existing record instead of creating a new one. Prevents duplicate transactions from retries or double-submits. For bill approvals, Go constructs this deterministically as `bill-{bill_id}-line-{line_number}` |
| `date`               | DATE          | NOT NULL                               | Transaction date                                                                                                                                                                                                                                                                                                    |
| `type`               | TEXT          | NOT NULL                               | One of: `PURCHASE`, `ISSUE_TO_JOB`, `ADJUSTMENT_IN`, `ADJUSTMENT_OUT`, `PRODUCTION_IN`, `SALE_OUT`, `REVERSAL`                                                                                                                                                                                                      |
| `item_type`          | TEXT          | NOT NULL                               | One of: `raw_material`, `finished_product`                                                                                                                                                                                                                                                                          |
| `material_id`        | UUID          | NULL, FK → `materials.id`              | Populated when `item_type = raw_material`. NULL otherwise                                                                                                                                                                                                                                                           |
| `variant_id`         | UUID          | NULL, FK → `product_variants.id`       | Populated when `item_type = finished_product`. NULL otherwise                                                                                                                                                                                                                                                       |
| `quantity`           | NUMERIC(15,2) | NOT NULL                               | Units moved. Always positive — direction is determined by transaction type                                                                                                                                                                                                                                          |
| `unit_cost`          | NUMERIC(15,2) | NULL                                   | Required for `PURCHASE`, `ADJUSTMENT_IN`, `PRODUCTION_IN`. NULL for issue and sale transactions                                                                                                                                                                                                                     |
| `total_cost`         | NUMERIC(15,2) | NULL                                   | Computed by Go: `quantity × unit_cost`. NULL where `unit_cost` is NULL                                                                                                                                                                                                                                              |
| `wac_at_transaction` | NUMERIC(15,2) | NULL                                   | The WAC in effect at the exact moment this transaction was created. Stored permanently for historical job costing accuracy. Populated for `ISSUE_TO_JOB` and `SALE_OUT`                                                                                                                                             |
| `order_id`           | UUID          | NULL, FK → `orders.id`                 | Required for `ISSUE_TO_JOB`. Optional for `PRODUCTION_IN`. NULL otherwise                                                                                                                                                                                                                                           |
| `supplier_id`        | UUID          | NULL, FK → `suppliers.id`              | Required for `PURCHASE` transactions. NULL otherwise                                                                                                                                                                                                                                                                |
| `purchase_order_id`  | UUID          | NULL, FK → `purchase_orders.id`        | Reference to originating PO for `PURCHASE` transactions                                                                                                                                                                                                                                                             |
| `bill_id`            | UUID          | NULL, FK → `bills.id`                  | Reference to the approved bill that triggered this transaction                                                                                                                                                                                                                                                      |
| `reversed_by`        | UUID          | NULL, FK → `inventory_transactions.id` | If this transaction has been reversed, points to the reversal entry                                                                                                                                                                                                                                                 |
| `reverses`           | UUID          | NULL, FK → `inventory_transactions.id` | If this is a reversal entry, points to the original transaction it reverses                                                                                                                                                                                                                                         |
| `reversal_reason`    | TEXT          | NULL                                   | Mandatory when `type = REVERSAL`. Go enforces this at the API level before insert. Stored permanently alongside the reversal entry                                                                                                                                                                                  |
| `notes`              | TEXT          | NULL                                   | Free-text context                                                                                                                                                                                                                                                                                                   |
| `status`             | TEXT          | NOT NULL, DEFAULT 'approved'           | One of: `pending`, `approved`, `reversed`. Adjustments start as `pending` until Super Admin approves. `reversed` means a reversal entry has been created against this transaction                                                                                                                                   |
| `created_by`         | TEXT          | NOT NULL                               | Staff ID of creator                                                                                                                                                                                                                                                                                                 |
| `approved_by`        | TEXT          | NULL                                   | Staff ID of Super Admin approver. Required for adjustment transactions                                                                                                                                                                                                                                              |
| `approved_at`        | TIMESTAMPTZ   | NULL                                   | Set by Go when approved, UTC                                                                                                                                                                                                                                                                                        |
| `created_at`         | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                                                                                                                                                                                                                                                      |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `txn_id`
- `UNIQUE` on `idempotency_key` (partial — where not null)
- Index on `material_id` — heavily queried for stock calculations
- Index on `variant_id` — heavily queried for stock calculations
- Index on `order_id` — queried for job costing
- Index on `type` — filtered in reports
- Index on `status` — filtered for pending approvals and active ledger
- Index on `date` — filtered in all date-range reports
- Index on `bill_id` — queried when tracing PO → Bill → Inventory chain

**Business rules:**

- Exactly one of `material_id` or `variant_id` must be populated —
  never both, never neither. Go enforces this before insert
- `PURCHASE` transactions are created automatically by Go when a
  supplier bill is approved and the linked PO status is `delivered`.
  Never created manually
- `ADJUSTMENT_IN` and `ADJUSTMENT_OUT` start with `status = pending`
  and require Super Admin approval before Go updates cache columns
- `REVERSAL` entries require a mandatory `reversal_reason`. Go creates
  the reversal row and sets the original transaction's `status = reversed`
  and `reversed_by` in a single DB transaction
- `wac_at_transaction` is captured at the moment of transaction creation
  from the current `current_wac` cache — this preserves historical
  job cost accuracy even after material prices change later
- Idempotency check runs before any insert:
  `SELECT id FROM inventory_transactions WHERE idempotency_key = $1`
  If found, return existing record. If not, proceed with insert

**Default query filter (active ledger):**

```sql
WHERE status != 'reversed'
AND type != 'REVERSAL'
```

**Full ledger query (include reversals):**

```sql
-- No filter — returns all rows including reversed and reversal entries
-- Exposed via: GET /api/v1/inventory/transactions?include_reversed=true
```

---

### Table: `wac_snapshots`

**Purpose:**
Every time a `PURCHASE` transaction is approved and triggers a WAC
recalculation, Go saves a timestamped snapshot here. This enables
point-in-time inventory valuation — answering "what was this material
worth on a specific past date" using the WAC that was actually in
effect at that time. Non-negotiable from Phase 2 day one —
retrofitting this later requires rebuilding the entire transaction
ledger.

| Column                    | Type          | Constraints                                | Notes                                                                           |
| ------------------------- | ------------- | ------------------------------------------ | ------------------------------------------------------------------------------- |
| `id`                      | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid()     | Internal identifier                                                             |
| `item_type`               | TEXT          | NOT NULL                                   | One of: `raw_material`, `finished_product`                                      |
| `material_id`             | UUID          | NULL, FK → `materials.id`                  | Populated when `item_type = raw_material`                                       |
| `variant_id`              | UUID          | NULL, FK → `product_variants.id`           | Populated when `item_type = finished_product`                                   |
| `wac`                     | NUMERIC(15,2) | NOT NULL                                   | The new weighted average cost after this recalculation                          |
| `previous_wac`            | NUMERIC(15,2) | NOT NULL                                   | The WAC before this recalculation. Useful for auditing price movement over time |
| `quantity_at_snapshot`    | NUMERIC(15,2) | NOT NULL                                   | Total stock quantity at the moment of this snapshot                             |
| `stock_value_at_snapshot` | NUMERIC(15,2) | NOT NULL                                   | Total stock value at the moment of this snapshot: `quantity × wac`              |
| `triggered_by_txn_id`     | UUID          | NOT NULL, FK → `inventory_transactions.id` | The purchase transaction that caused this WAC recalculation                     |
| `created_at`              | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                    | Set by Go, UTC                                                                  |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `material_id`
- Index on `variant_id`
- Index on `created_at` — queried for date-range valuation reports
- Index on `triggered_by_txn_id`

**Business rules:**

- Written by Go inside the same database transaction as the WAC
  recalculation — if the transaction rolls back, the snapshot rolls
  back with it
- Never updated or deleted after creation — append only
- Exactly one of `material_id` or `variant_id` must be populated
- Point-in-time valuation query:
  ```sql
  SELECT wac FROM wac_snapshots
  WHERE material_id = $1
  AND created_at <= $target_date
  ORDER BY created_at DESC
  LIMIT 1
  ```

---

## WAC Recalculation — Go Sequence

Every approved `PURCHASE` transaction triggers this sequence inside
a single database transaction. All steps succeed together or roll
back together:

```
1. Read current_quantity and current_wac from materials cache

2. Compute new WAC:
   current_stock_value = current_quantity × current_wac
   new_purchase_value  = new_quantity × unit_cost
   new_wac = (current_stock_value + new_purchase_value)
             ÷ (current_quantity + new_quantity)

3. Update materials cache:
   current_quantity    += new_quantity
   current_wac          = new_wac
   current_stock_value  = updated_quantity × new_wac

4. Write wac_snapshots row:
   previous_wac             = old current_wac
   wac                      = new_wac
   quantity_at_snapshot     = updated current_quantity
   stock_value_at_snapshot  = updated current_stock_value

5. Set inventory_transaction.status = approved

6. Check: if current_quantity <= reorder_level → trigger low stock alert

7. Commit transaction
```

---

## Relationships in Group 4

```
inventory_transactions >── materials (Group 3)
inventory_transactions >── product_variants (Group 3)
inventory_transactions >── orders (Group 7)
inventory_transactions >── suppliers (Group 3)
inventory_transactions >── purchase_orders (Group 5)
inventory_transactions >── bills (Group 6)
inventory_transactions ──< inventory_transactions (self — reversal chain)

wac_snapshots >── materials (Group 3)
wac_snapshots >── product_variants (Group 3)
wac_snapshots >── inventory_transactions
```

---

---

## GROUP 5 — PROCUREMENT

**Tables:** `purchase_orders`, `purchase_order_lines`, `po_delivery_logs`

A Purchase Order is a formal procurement commitment sent to a supplier.
It does not affect inventory or accounting at creation. Inventory updates
only after the PO is delivered AND the linked bill is approved — both
conditions must be met.

---

### Table: `purchase_orders`

| Column                   | Type         | Constraints                            | Notes                                                                                  |
| ------------------------ | ------------ | -------------------------------------- | -------------------------------------------------------------------------------------- |
| `id`                     | UUID         | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                    |
| `po_id`                  | TEXT         | NOT NULL, UNIQUE                       | Format: `PO/YYYY/MM/DD/001`. Generated by Go                                           |
| `supplier_id`            | UUID         | NOT NULL, FK → `suppliers.id`          | The supplier this PO is sent to                                                        |
| `date`                   | DATE         | NOT NULL                               | PO creation date                                                                       |
| `expected_delivery_date` | DATE         | NULL                                   | Target delivery date                                                                   |
| `status`                 | TEXT         | NOT NULL, DEFAULT 'draft'              | One of: `draft`, `submitted`, `supplier_confirmed`, `delivered`, `closed`, `cancelled` |
| `quantity_tolerance`     | NUMERIC(5,2) | NOT NULL, DEFAULT 5.00                 | Acceptable variance percentage between ordered and delivered quantity                  |
| `notes`                  | TEXT         | NULL                                   | Delivery instructions for supplier                                                     |
| `created_by`             | TEXT         | NOT NULL                               | Staff ID of creator                                                                    |
| `approved_by`            | TEXT         | NULL                                   | Staff ID of Super Admin who approved                                                   |
| `approved_at`            | TIMESTAMPTZ  | NULL                                   | Set by Go when approved, UTC                                                           |
| `cancelled_by`           | TEXT         | NULL                                   | Staff ID of canceller                                                                  |
| `cancellation_reason`    | TEXT         | NULL                                   | Mandatory when `status = cancelled`                                                    |
| `created_at`             | TIMESTAMPTZ  | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                         |
| `updated_at`             | TIMESTAMPTZ  | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                                         |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `po_id`
- Index on `supplier_id`
- Index on `status`
- Index on `date`

**Business rules:**

- PO must be approved by Super Admin before being sent to supplier
- PO does not affect inventory or accounting at creation
- Edits after approval require Super Admin authorisation with logged reason
- PO advances to `closed` only when all lines have `is_fully_delivered = true`
- Cancellation requires a mandatory reason logged to audit

---

### Table: `purchase_order_lines`

| Column               | Type          | Constraints                                           | Notes                                                                                  |
| -------------------- | ------------- | ----------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `id`                 | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid()                | Internal identifier                                                                    |
| `purchase_order_id`  | UUID          | NOT NULL, FK → `purchase_orders.id` ON DELETE CASCADE | Parent PO                                                                              |
| `item_type`          | TEXT          | NOT NULL                                              | One of: `raw_material`, `finished_product`                                             |
| `material_id`        | UUID          | NULL, FK → `materials.id`                             | Populated when `item_type = raw_material`                                              |
| `variant_id`         | UUID          | NULL, FK → `product_variants.id`                      | Populated when `item_type = finished_product`                                          |
| `description`        | TEXT          | NULL                                                  | Free-text for context or non-catalogued items                                          |
| `quantity_ordered`   | NUMERIC(15,2) | NOT NULL                                              | Units requested from supplier                                                          |
| `quantity_delivered` | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                   | Running total updated by Go as the sum of all `po_delivery_logs` entries for this line |
| `unit_of_measure`    | TEXT          | NOT NULL                                              | e.g. `metres`, `kg`, `pcs`                                                             |
| `unit_cost`          | NUMERIC(15,2) | NOT NULL                                              | Agreed unit cost with supplier                                                         |
| `total_cost`         | NUMERIC(15,2) | NOT NULL                                              | Computed by Go: `quantity_ordered × unit_cost`                                         |
| `is_fully_delivered` | BOOLEAN       | NOT NULL, DEFAULT false                               | Set to true by Go when `quantity_delivered >= quantity_ordered` within tolerance       |
| `created_at`         | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                               | Set by Go, UTC                                                                         |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `purchase_order_id`
- Index on `material_id`
- Index on `variant_id`

**Business rules:**

- Exactly one of `material_id` or `variant_id` must be populated
- `quantity_delivered` is the sum of all `po_delivery_logs.quantity_delivered`
  for this line — recomputed by Go after every delivery log insert
- PO advances to `closed` when all lines have `is_fully_delivered = true`
- Variance beyond `quantity_tolerance` is flagged in the audit log

---

### Table: `po_delivery_logs`

**Purpose:**
Records every individual delivery event against a PO line. Each
partial or full delivery creates one row per line. This gives a
complete, timestamped history of how a PO was fulfilled — who
confirmed it, when, and how much arrived each time.

| Column                   | Type          | Constraints                              | Notes                                                        |
| ------------------------ | ------------- | ---------------------------------------- | ------------------------------------------------------------ |
| `id`                     | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid()   | Internal identifier                                          |
| `purchase_order_id`      | UUID          | NOT NULL, FK → `purchase_orders.id`      | Parent PO                                                    |
| `purchase_order_line_id` | UUID          | NOT NULL, FK → `purchase_order_lines.id` | The specific line being delivered against                    |
| `delivery_date`          | DATE          | NOT NULL                                 | Date this delivery was received at the factory               |
| `quantity_delivered`     | NUMERIC(15,2) | NOT NULL                                 | Quantity received in this specific delivery event            |
| `notes`                  | TEXT          | NULL                                     | Context e.g. damage noted, short delivery, substitution      |
| `confirmed_by`           | TEXT          | NOT NULL                                 | Staff ID of Admin or Inventory Officer who confirmed receipt |
| `created_at`             | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                  | Set by Go, UTC                                               |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `purchase_order_id`
- Index on `purchase_order_line_id`
- Index on `delivery_date`

**Business rules:**

- Every delivery event — full or partial — creates one row per line
- After each insert, Go recomputes `purchase_order_lines.quantity_delivered`
  as the SUM of all `po_delivery_logs.quantity_delivered` for that line
- If the recomputed total meets or exceeds `quantity_ordered` within
  tolerance, Go sets `is_fully_delivered = true` on the line
- When all lines are fully delivered, Go advances PO status to `delivered`
- Records are permanent — never deleted or updated after creation

---

## Relationships in Group 5

```
purchase_orders >── suppliers (Group 3)
purchase_orders ──< purchase_order_lines
purchase_orders ──< po_delivery_logs
purchase_order_lines ──< po_delivery_logs
purchase_order_lines >── materials (Group 3)
purchase_order_lines >── product_variants (Group 3)
purchase_orders ──< bills (Group 6)
purchase_orders ──< inventory_transactions (Group 4)
```

---

---

## GROUP 6 — PAYABLES

**Tables:** `bills`, `bill_payments`

Supplier obligations are tracked through Bills. Inventory-related bills
must always link to a Purchase Order. Service bills — utilities,
logistics, artisan fees — can exist without a PO. Bills support partial
payments and reversals. Reversed bills create a new opposing bill row
so the full correction history is always visible in the ledger.

---

### Table: `bills`

**Purpose:**
Represents money owed to a supplier for goods or services received.
Inventory bills trigger the inventory update sequence on approval.
Service bills do not affect inventory. All bills are permanent —
corrections are made through reversal entries, never by deletion.

| Column                | Type          | Constraints                            | Notes                                                                                                                  |
| --------------------- | ------------- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `id`                  | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                    |
| `bill_id`             | TEXT          | NOT NULL, UNIQUE                       | Format: `BILL/YYYY/MM/DD/001`. Generated by Go                                                                         |
| `supplier_id`         | UUID          | NOT NULL, FK → `suppliers.id`          | The supplier this bill is from                                                                                         |
| `purchase_order_id`   | UUID          | NULL, FK → `purchase_orders.id`        | Mandatory when `is_inventory_bill = true`. Must be NULL when `is_inventory_bill = false`                               |
| `is_inventory_bill`   | BOOLEAN       | NOT NULL, DEFAULT false                | True for bills linked to inventory purchases. False for service bills. Go enforces the PO link rule based on this flag |
| `category`            | TEXT          | NOT NULL                               | One of: `raw_materials`, `printing`, `logistics`, `utilities`, `artisans`, `other`                                     |
| `description`         | TEXT          | NOT NULL                               | Description of goods or services billed                                                                                |
| `amount`              | NUMERIC(15,2) | NOT NULL                               | Total bill amount in NGN                                                                                               |
| `amount_paid`         | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | Cache. Running total of all payments. Updated by Go after every bill payment                                           |
| `status`              | TEXT          | NOT NULL, DEFAULT 'unpaid'             | One of: `unpaid`, `part_paid`, `paid`, `reversed`. Updated automatically by Go                                         |
| `due_date`            | DATE          | NULL                                   | Payment due date                                                                                                       |
| `reversed_by_bill_id` | UUID          | NULL, FK → `bills.id`                  | If this bill has been reversed, points to the reversal bill entry                                                      |
| `reverses_bill_id`    | UUID          | NULL, FK → `bills.id`                  | If this is a reversal entry, points to the original bill it reverses                                                   |
| `reversal_reason`     | TEXT          | NULL                                   | Mandatory when this bill is a reversal entry. Go enforces before insert                                                |
| `created_by`          | TEXT          | NOT NULL                               | Staff ID of creator                                                                                                    |
| `approved_by`         | TEXT          | NULL                                   | Staff ID of Super Admin who approved                                                                                   |
| `approved_at`         | TIMESTAMPTZ   | NULL                                   | Set by Go when approved. Triggers inventory update for inventory bills                                                 |
| `posted_at`           | TIMESTAMPTZ   | NULL                                   | Set by Go when finalised and locked. Not client-supplied                                                               |
| `created_at`          | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                                                         |
| `updated_at`          | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                                                                         |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `bill_id`
- Index on `supplier_id`
- Index on `purchase_order_id`
- Index on `status`
- Index on `is_inventory_bill`
- Index on `due_date`

**Business rules:**

- If `is_inventory_bill = true`, `purchase_order_id` is mandatory and
  the linked PO must have `status = delivered` before Go approves the bill
- If `is_inventory_bill = false`, `purchase_order_id` must be NULL
- Bill approval triggers inventory transaction creation for inventory bills
- Status update logic after every payment:
  ```
  amount_paid = 0                → unpaid
  0 < amount_paid < amount       → part_paid
  amount_paid >= amount          → paid
  ```
- Reversal flow — all steps in a single DB transaction:
  ```
  1. Go creates a new bill row with reverses_bill_id = original bill ID,
     reversal_reason = provided reason, status = reversed
  2. Go updates original bill: reversed_by_bill_id = new bill ID,
     status = reversed
  3. If original bill had payments, Go creates opposing
     cash_transactions entries
  4. Audit log entry written
  ```
- Default fetch excludes reversed bills:
  WHERE status != 'reversed' AND reverses_bill_id IS NULL
- Full fetch: GET /api/v1/bills?include_reversed=true

---

### Table: `bill_payments`

**Purpose:**
Records every payment made against a supplier bill. Supports partial
payments. Each payment automatically creates a cash_transactions
outflow entry and updates the parent bill status.

| Column        | Type          | Constraints                            | Notes                                                                                                              |
| ------------- | ------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `id`          | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                |
| `payment_id`  | TEXT          | NOT NULL, UNIQUE                       | Format: `BLL/YYYY/MM/DD/001`. Generated by Go. BLL prefix distinguishes bill payments from customer payments (PAY) |
| `bill_id`     | UUID          | NOT NULL, FK → `bills.id`              | The bill this payment is applied to                                                                                |
| `supplier_id` | UUID          | NOT NULL, FK → `suppliers.id`          | Denormalised for faster supplier statement queries                                                                 |
| `date`        | DATE          | NOT NULL                               | Payment date                                                                                                       |
| `amount`      | NUMERIC(15,2) | NOT NULL                               | Amount paid in this transaction                                                                                    |
| `method`      | TEXT          | NOT NULL                               | One of: `cash`, `bank_transfer`, `pos`                                                                             |
| `reference`   | TEXT          | NULL                                   | Bank reference, receipt number, or POS transaction ID                                                              |
| `recorded_by` | TEXT          | NOT NULL                               | Staff ID of staff member who recorded the payment                                                                  |
| `created_at`  | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                                                     |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `payment_id`
- Index on `bill_id`
- Index on `supplier_id`
- Index on `date`

**Business rules:**

- Payment amount cannot exceed outstanding balance:
  amount <= (bill.amount - bill.amount_paid). Go enforces before insert
- After every insert, Go runs inside a single DB transaction:
  1. Insert bill_payments row
  2. Update bills.amount_paid += payment.amount
  3. Recompute and update bills.status
  4. Create cash_transactions outflow entry
  5. Write audit_log entry
     All five steps succeed together or roll back together
- Payments are never deleted — corrections require a bill reversal

---

## Relationships in Group 6

```
bills >── suppliers (Group 3)
bills >── purchase_orders (Group 5)
bills ──< bill_payments
bills ──< inventory_transactions (Group 4)
bills ──< bills (self — reversal chain)
bill_payments >── bills
bill_payments >── suppliers (Group 3)
bill_payments ──> cash_transactions (Group 10)
```

---

_BloansBooks ERD · Internal Technical Documentation · Confidential_
