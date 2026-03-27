# BloansBooks — Entity Relationship Document

**Version:** 1.0  
**Last Updated:** March 2026  
**Status:** Complete

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

## GROUP 1 — Auth & Security

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

| Column               | Type        | Constraints                            | Notes                                                                                                                                                                            |
| -------------------- | ----------- | -------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`                 | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                                                                                              |
| `staff_id`           | TEXT        | NOT NULL                               | The terminated staff member's business ID e.g. `BLN-0042`. Loose reference to `staff.staff_id` — deliberately no FK so the block persists even if staff records are restructured |
| `termination_date`   | DATE        | NOT NULL                               | Date the staff member was terminated                                                                                                                                             |
| `termination_reason` | TEXT        | NOT NULL                               | Mandatory termination reason e.g. Resigned, Dismissed, Contract Ended, Redundancy                                                                                                |
| `recorded_by`        | TEXT        | NOT NULL                               | Staff ID of the HR Manager or Super Admin who recorded the termination                                                                                                           |
| `recorded_at`        | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go backend, UTC                                                                                                                                                           |
| `is_overridden`      | BOOLEAN     | NOT NULL, DEFAULT false                | Set to true if the termination was recorded in error and a Super Admin has overridden it                                                                                         |
| `overridden_by`      | TEXT        | NULL                                   | Staff ID of the Super Admin who overrode the termination. NULL until overridden                                                                                                  |
| `overridden_at`      | TIMESTAMPTZ | NULL                                   | Timestamp of the override. Set by Go, UTC. NULL until overridden                                                                                                                 |
| `override_reason`    | TEXT        | NULL                                   | Mandatory when `is_overridden` is set to true. Must be provided by the Super Admin                                                                                               |

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

## GROUP 2 — Staff

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
| `status`                  | TEXT          | NOT NULL, DEFAULT 'active'             | One of: `active`, `fired`, `inactive`. Inactive staff are excluded from payroll runs and task assignment. Fired staff cannot login                                        |
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

## Full Relationship Map — Groups 1 & 2

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

## GROUP 3 — Master Data

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
| `currency`    | TEXT        | NOT NULL, DEFAULT 'NGN'                | Billing currency. NGN only for v1                                                  |
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

## GROUP 4 — Inventory

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

## GROUP 5 — Procurement

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

## GROUP 6 — Payables

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

---

## GROUP 7 — Sales

**Tables:** `orders`, `order_lines`, `invoices`, `customer_payments`

The sales workflow begins with a Sales Order and ends when the invoice
is fully paid and the order is closed. One order has exactly one invoice.
Partial deliveries are tracked at the order line level — the invoice
always covers the full order value regardless of delivery progress.

---

### Table: `orders`

**Purpose:**
A Sales Order represents a customer request to purchase one or more
product variants. Starting point of the sales workflow. Must be approved
before an invoice can be raised. Closes automatically when the linked
invoice is fully paid.

| Column                   | Type        | Constraints                            | Notes                                                               |
| ------------------------ | ----------- | -------------------------------------- | ------------------------------------------------------------------- |
| `id`                     | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                 |
| `order_id`               | TEXT        | NOT NULL, UNIQUE                       | Format: `BH/YYYY/MM/DD/001`. Generated by Go                        |
| `customer_id`            | UUID        | NOT NULL, FK → `customers.id`          | The customer placing the order                                      |
| `date`                   | DATE        | NOT NULL                               | Order creation date                                                 |
| `expected_delivery_date` | DATE        | NULL                                   | Target delivery date                                                |
| `status`                 | TEXT        | NOT NULL, DEFAULT 'draft'              | One of: `draft`, `approved`, `in_production`, `delivered`, `closed` |
| `customisation_notes`    | TEXT        | NULL                                   | Free-text job-specific instructions                                 |
| `created_by`             | TEXT        | NOT NULL                               | Staff ID of creator                                                 |
| `approved_by`            | TEXT        | NULL                                   | Staff ID of approver                                                |
| `approved_at`            | TIMESTAMPTZ | NULL                                   | Set by Go when approved, UTC                                        |
| `created_at`             | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                      |
| `updated_at`             | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                      |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `order_id`
- Index on `customer_id`
- Index on `status`
- Index on `date`

**Business rules:**

- Draft orders are editable by Admin and Sales Manager
- Once approved, order is locked for editing
- `in_production` is set when Production Supervisor begins issuing materials
- `delivered` is set automatically when all order lines have `delivery_status = delivered`
- `closed` is set automatically when the linked invoice is fully paid
- Orders are never deleted

---

### Table: `order_lines`

**Purpose:**
Each row is a specific product variant being ordered. Partial delivery
is tracked per line via `quantity_delivered` and `delivery_status`.
The invoice always covers the full order value regardless of delivery progress.
Production Supervisor can view lines when issuing materials to a job.

| Column               | Type          | Constraints                                  | Notes                                                                                    |
| -------------------- | ------------- | -------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `id`                 | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid()       | Internal identifier                                                                      |
| `order_id`           | UUID          | NOT NULL, FK → `orders.id` ON DELETE CASCADE | Parent order                                                                             |
| `variant_id`         | UUID          | NOT NULL, FK → `product_variants.id`         | The specific product variant being ordered                                               |
| `quantity`           | NUMERIC(15,2) | NOT NULL                                     | Total units ordered                                                                      |
| `quantity_delivered` | NUMERIC(15,2) | NOT NULL, DEFAULT 0                          | Updated by Go when Production Supervisor marks units delivered. Cannot exceed `quantity` |
| `delivery_status`    | TEXT          | NOT NULL, DEFAULT 'pending'                  | One of: `pending`, `partial`, `delivered`. Updated automatically by Go                   |
| `unit_price`         | NUMERIC(15,2) | NOT NULL                                     | Selling price per unit. Defaults to `product_variants.selling_price`. Can be overridden  |
| `total_price`        | NUMERIC(15,2) | NOT NULL                                     | Computed by Go: `quantity x unit_price`. Never accepted from client                      |
| `notes`              | TEXT          | NULL                                         | Line-specific notes                                                                      |
| `created_at`         | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                      | Set by Go, UTC                                                                           |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `order_id`
- Index on `variant_id`
- Index on `delivery_status`

**Business rules:**

- At least one line must exist before an order can be approved
- Delivery status logic:
  quantity_delivered = 0 → pending
  0 < quantity_delivered < quantity → partial
  quantity_delivered >= quantity → delivered
- When all lines are delivered, Go advances order status to `delivered`
- Lines are locked once the parent order is approved

---

### Table: `invoices`

**Purpose:**
Financial document issued to the customer for a specific order.
Exactly one invoice per order enforced at DB level. Tracks amount owed,
discounts, and payment status. Drives outstanding receivables on the dashboard.
Supports reversals via self-referencing pattern identical to bills.

| Column                     | Type          | Constraints                            | Notes                                                                  |
| -------------------------- | ------------- | -------------------------------------- | ---------------------------------------------------------------------- |
| `id`                       | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                    |
| `invoice_id`               | TEXT          | NOT NULL, UNIQUE                       | Format: `INV/YYYY/MM/DD/001`. Generated by Go                          |
| `order_id`                 | UUID          | NOT NULL, UNIQUE, FK → `orders.id`     | Linked sales order. UNIQUE enforces one invoice per order at DB level  |
| `customer_id`              | UUID          | NOT NULL, FK → `customers.id`          | Denormalised for faster customer statement queries                     |
| `date`                     | DATE          | NOT NULL                               | Invoice creation date                                                  |
| `due_date`                 | DATE          | NULL                                   | Payment due date                                                       |
| `subtotal`                 | NUMERIC(15,2) | NOT NULL                               | Sum of all order line totals before discount                           |
| `discount_type`            | TEXT          | NULL                                   | One of: `percentage`, `fixed`. NULL if no discount                     |
| `discount_value`           | NUMERIC(15,2) | NULL                                   | The discount amount or percentage entered                              |
| `discount_amount`          | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | Computed by Go: actual NGN amount deducted                             |
| `total`                    | NUMERIC(15,2) | NOT NULL                               | Computed by Go: `subtotal - discount_amount`                           |
| `amount_paid`              | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | Cache. Running total of all payments. Updated after every payment      |
| `status`                   | TEXT          | NOT NULL, DEFAULT 'unpaid'             | One of: `unpaid`, `part_paid`, `paid`, `cancelled`, `reversed`         |
| `discount_approval_status` | TEXT          | NULL                                   | One of: `pending`, `approved`, `rejected`. NULL if no threshold breach |
| `discount_approved_by`     | TEXT          | NULL                                   | Staff ID of Super Admin who approved the discount                      |
| `discount_approved_at`     | TIMESTAMPTZ   | NULL                                   | Set by Go when discount approved, UTC                                  |
| `reversed_by_invoice_id`   | UUID          | NULL, FK → `invoices.id`               | Points to the reversal invoice if this invoice has been reversed       |
| `reverses_invoice_id`      | UUID          | NULL, FK → `invoices.id`               | Points to the original invoice if this is a reversal entry             |
| `reversal_reason`          | TEXT          | NULL                                   | Mandatory when this invoice is a reversal entry                        |
| `created_by`               | TEXT          | NOT NULL                               | Staff ID of creator                                                    |
| `posted_at`                | TIMESTAMPTZ   | NULL                                   | Set by Go when finalised and locked. Not client-supplied               |
| `created_at`               | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                         |
| `updated_at`               | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                         |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `invoice_id`
- `UNIQUE` on `order_id`
- Index on `customer_id`
- Index on `status`
- Index on `due_date`
- Index on `discount_approval_status`

**Business rules:**

- Cannot be posted if `discount_approval_status = pending`
- If discount exceeds configured threshold, Go sets `discount_approval_status = pending`
- Status update logic after every payment:
  amount_paid = 0 → unpaid
  0 < amount_paid < total → part_paid
  amount_paid >= total → paid → order.status = closed
- Default fetch excludes reversed:
  WHERE status != 'reversed' AND reverses_invoice_id IS NULL
- Full fetch: GET /api/v1/invoices?include_reversed=true

---

### Table: `customer_payments`

**Purpose:**
Records every payment received from a customer against an invoice.
Supports partial payments. Each payment automatically creates a
cash_transactions inflow entry and updates parent invoice status.
When invoice is fully paid, the linked order closes automatically.

| Column        | Type          | Constraints                            | Notes                                                 |
| ------------- | ------------- | -------------------------------------- | ----------------------------------------------------- |
| `id`          | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                   |
| `payment_id`  | TEXT          | NOT NULL, UNIQUE                       | Format: `PAY/YYYY/MM/DD/001`. Generated by Go         |
| `invoice_id`  | UUID          | NOT NULL, FK → `invoices.id`           | The invoice this payment is applied to                |
| `customer_id` | UUID          | NOT NULL, FK → `customers.id`          | Denormalised for faster customer statement queries    |
| `order_id`    | UUID          | NOT NULL, FK → `orders.id`             | Denormalised for faster order-level payment queries   |
| `date`        | DATE          | NOT NULL                               | Payment received date                                 |
| `amount`      | NUMERIC(15,2) | NOT NULL                               | Amount received in this payment                       |
| `method`      | TEXT          | NOT NULL                               | One of: `cash`, `bank_transfer`, `pos`                |
| `reference`   | TEXT          | NULL                                   | Bank reference, receipt number, or POS transaction ID |
| `recorded_by` | TEXT          | NOT NULL                               | Staff ID of recorder                                  |
| `created_at`  | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                        |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `payment_id`
- Index on `invoice_id`
- Index on `customer_id`
- Index on `order_id`
- Index on `date`

**Business rules:**

- Payment cannot exceed outstanding balance:
  amount <= (invoice.total - invoice.amount_paid). Go enforces before insert
- After every insert, Go runs inside a single DB transaction:
  1. Insert customer_payments row
  2. Update invoices.amount_paid += payment.amount
  3. Recompute and update invoices.status
  4. If invoices.status = paid, set orders.status = closed
  5. Create cash_transactions inflow entry
  6. Write audit_log entry
     All six steps succeed together or roll back together
- Payments are never deleted — corrections require an invoice reversal

---

## Relationships in Group 7

```
-- orders
customers (1) ──< orders (many)
  One customer can have many orders.
  Each order belongs to exactly one customer.

orders (1) ──< order_lines (many)
  One order can have many lines.
  Each line belongs to exactly one order.
  Cascade delete on order_lines if order deleted (orders never deleted in practice).

product_variants (1) ──< order_lines (many)
  One variant can appear on many order lines across different orders.
  Each order line references exactly one variant.

-- invoices
orders (1) ──── invoices (1)
  One order has exactly one invoice.
  One invoice belongs to exactly one order.
  Enforced by UNIQUE constraint on invoices.order_id.

customers (1) ──< invoices (many)
  One customer can have many invoices across all their orders.
  Each invoice belongs to exactly one customer. Denormalised.

invoices (1) ──── invoices (self, 1)
  A reversal invoice points back to the original via reverses_invoice_id.
  The original points forward to the reversal via reversed_by_invoice_id.
  Self-referencing FK. Each invoice can have at most one reversal.

-- customer_payments
invoices (1) ──< customer_payments (many)
  One invoice can have many payments. Partial payment support.
  Each payment is applied to exactly one invoice.

customers (1) ──< customer_payments (many)
  Denormalised. One customer can have many payments across all invoices.
  Each payment belongs to exactly one customer.

orders (1) ──< customer_payments (many)
  Denormalised. One order can have many payments.
  Each payment belongs to exactly one order.

customer_payments (many) ──> cash_transactions (Group 10) (1 per payment)
  Every customer payment creates exactly one cash_transactions inflow entry.
  Created by Go inside the same DB transaction as the payment insert.
```

---

---

## GROUP 8 — Job Costing

**Tables:** `job_labor`, `job_overhead`

Job costing gives the business a precise profit figure per order.
Materials cost is derived automatically from inventory transactions —
no separate table needed. Labor and overhead are captured here and
combined with materials to produce the full job cost summary which
is computed on demand by Go, never stored.

---

### Table: `job_labor`

**Purpose:**
Records direct labor applied to a specific order. Each row captures
one staff member's time on one job. Labor cost feeds into the job
cost summary alongside materials and overhead.

| Column        | Type          | Constraints                            | Notes                                                                                                         |
| ------------- | ------------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| `id`          | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                           |
| `order_id`    | UUID          | NOT NULL, FK → `orders.id`             | The order this labor is charged to                                                                            |
| `staff_id`    | TEXT          | NOT NULL, FK → `staff.staff_id`        | The staff member whose time is being recorded                                                                 |
| `hours`       | NUMERIC(8,2)  | NOT NULL                               | Hours worked on this job                                                                                      |
| `rate`        | NUMERIC(15,2) | NOT NULL                               | Hourly rate in NGN. Defaults to staff base salary converted to hourly equivalent. Can be overridden per entry |
| `total`       | NUMERIC(15,2) | NOT NULL                               | Computed by Go: `hours x rate`. Never accepted from client                                                    |
| `notes`       | TEXT          | NULL                                   | Optional task description for this labor entry                                                                |
| `recorded_by` | TEXT          | NOT NULL                               | Staff ID of recorder                                                                                          |
| `created_at`  | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                                                |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `order_id` — queried every time job cost summary is computed
- Index on `staff_id` — queried for staff labor reports

**Business rules:**

- Multiple labor entries can exist per order and per staff member
- `rate` defaults to staff base_salary converted to hourly equivalent.
  Can be overridden per entry by Admin or Production Supervisor
- `total` is always computed by Go — never accepted from client
- Labor entries are locked once the parent order is closed
- Materials cost is derived from inventory_transactions where
  type = ISSUE_TO_JOB and order_id matches — no duplication needed

---

### Table: `job_overhead`

**Purpose:**
Records overhead costs allocated to a specific order. Covers generator
usage, machine time, and miscellaneous costs that cannot be attributed
to a specific material or staff member but are genuinely incurred in
fulfilling the order.

| Column        | Type          | Constraints                            | Notes                                                                              |
| ------------- | ------------- | -------------------------------------- | ---------------------------------------------------------------------------------- |
| `id`          | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                |
| `order_id`    | UUID          | NOT NULL, FK → `orders.id`             | The order this overhead is charged to                                              |
| `type`        | TEXT          | NOT NULL                               | One of: `generator`, `machine`, `miscellaneous`                                    |
| `amount`      | NUMERIC(15,2) | NOT NULL                               | Overhead cost in NGN allocated to this order                                       |
| `basis`       | TEXT          | NULL                                   | Optional allocation basis e.g. `per_hour`, `per_order`. For reporting context only |
| `notes`       | TEXT          | NULL                                   | Description of the overhead being allocated                                        |
| `recorded_by` | TEXT          | NOT NULL                               | Staff ID of recorder                                                               |
| `created_at`  | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                     |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `order_id` — queried every time job cost summary is computed
- Index on `type` — filtered in overhead reports

**Business rules:**

- Multiple overhead entries can exist per order
- Overhead entries are locked once the parent order is closed
- `amount` is entered directly — no formula involved

---

## Job Cost Summary — Computed by Go on Demand

The job cost summary is not a stored table. It is computed by Go
by aggregating from four sources:

```
Revenue
  = invoices.total WHERE order_id = target order

Materials Cost
  = SUM(inventory_transactions.total_cost)
    WHERE type = ISSUE_TO_JOB AND order_id = target order
    Uses wac_at_transaction for historical accuracy

Labor Cost
  = SUM(job_labor.total) WHERE order_id = target order

Overhead Cost
  = SUM(job_overhead.amount) WHERE order_id = target order

COGS          = Materials Cost + Labor Cost + Overhead Cost
Gross Profit  = Revenue - COGS
Gross Margin% = (Gross Profit / Revenue) x 100
```

---

## Relationships in Group 8

```
-- job_labor
orders (1) ──< job_labor (many)
  One order can have many labor entries.
  Each labor entry belongs to exactly one order.

staff (1) ──< job_labor (many)
  One staff member can have labor entries across many orders.
  Each labor entry is attributed to exactly one staff member.

-- job_overhead
orders (1) ──< job_overhead (many)
  One order can have many overhead entries.
  Each overhead entry belongs to exactly one order.

-- job cost summary inputs (no stored table)
orders (1) ──── invoices (1)               → revenue source
orders (1) ──< inventory_transactions      → materials cost (Group 4)
orders (1) ──< job_labor                   → labor cost
orders (1) ──< job_overhead                → overhead cost
```

---

---

## GROUP 9 — Cash Transactions

**Tables:** `cash_transactions`

A unified ledger of every financial movement in the system. Never
created manually — always created by Go as a side effect of a customer
payment, bill payment, or payroll run. Drives cash and bank balance
tiles on the dashboard.

---

### Table: `cash_transactions`

**Purpose:**
Unified inflow and outflow ledger across cash and bank accounts.
Every entry is system-generated by Go inside the same DB transaction
as the triggering event. Never created, updated, or deleted manually.
Corrections create opposing adjustment entries authorised by Super Admin.

| Column               | Type          | Constraints                            | Notes                                                                            |
| -------------------- | ------------- | -------------------------------------- | -------------------------------------------------------------------------------- |
| `id`                 | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                              |
| `txn_id`             | TEXT          | NOT NULL, UNIQUE                       | Format: `CTX/YYYY/MM/DD/001`. Generated by Go                                    |
| `date`               | DATE          | NOT NULL                               | Transaction date                                                                 |
| `type`               | TEXT          | NOT NULL                               | One of: `inflow`, `outflow`                                                      |
| `amount`             | NUMERIC(15,2) | NOT NULL                               | Always positive. Direction determined by `type`                                  |
| `account`            | TEXT          | NOT NULL                               | One of: `cash`, `bank`                                                           |
| `linked_entity_type` | TEXT          | NOT NULL                               | One of: `customer_payment`, `bill_payment`, `payroll_run`, `adjustment`          |
| `linked_entity_id`   | UUID          | NOT NULL                               | UUID of the source record. Loose reference — no FK, resolved by Go at query time |
| `reference`          | TEXT          | NULL                                   | Transaction reference or description                                             |
| `created_by`         | TEXT          | NOT NULL                               | Staff ID or `SYSTEM` for Go-generated entries                                    |
| `created_at`         | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                   |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `txn_id`
- Index on `date`
- Index on `type`
- Index on `account`
- Index on `linked_entity_type`
- Index on `linked_entity_id`

**Business rules:**

- Never created manually — Go only
- Every customer payment → one inflow entry
- Every bill payment → one outflow entry
- Every payroll run marked paid → one outflow entry for total net pay
- All created inside the same DB transaction as the triggering event
- Never updated or deleted — corrections via Super Admin authorised adjustment entries
- linked_entity_id is a loose reference — no FK constraint since it
  points to different tables depending on linked_entity_type

**Dashboard queries:**

```
Cash balance
  = SUM(amount) WHERE account = cash AND type = inflow
  - SUM(amount) WHERE account = cash AND type = outflow

Bank balance
  = SUM(amount) WHERE account = bank AND type = inflow
  - SUM(amount) WHERE account = bank AND type = outflow

This month revenue
  = SUM(amount) WHERE linked_entity_type = customer_payment
    AND date >= first day of current month

This month expenses
  = SUM(amount) WHERE linked_entity_type IN (bill_payment, payroll_run)
    AND date >= first day of current month
```

---

## Relationships in Group 9

```
cash_transactions ....> customer_payments (Group 7)
  Loose reference via linked_entity_id where linked_entity_type = customer_payment.
  Every customer payment creates exactly one inflow cash_transaction.
  No FK — resolved by Go at query time.

cash_transactions ....> bill_payments (Group 6)
  Loose reference via linked_entity_id where linked_entity_type = bill_payment.
  Every bill payment creates exactly one outflow cash_transaction.
  No FK — resolved by Go at query time.

cash_transactions ....> payroll_runs (Group 10)
  Loose reference via linked_entity_id where linked_entity_type = payroll_run.
  Every payroll run marked paid creates exactly one outflow cash_transaction.
  No FK — resolved by Go at query time.
```

---

---

## GROUP 10 — Payroll

**Tables:** `payroll_runs`, `payroll_lines`

One monthly payroll cycle per run. Cannot be approved if unresolved
attendance exceptions exist. Marking paid automatically creates a
cash_transactions outflow. All payroll line values are computed by Go
from attendance data and system settings.

---

### Table: `payroll_runs`

**Purpose:**
Represents one monthly payroll cycle. Created by Admin, approved by
Super Admin. Go blocks approval if unresolved attendance exceptions
exist for any active staff member for that month.

| Column          | Type          | Constraints                            | Notes                                                                         |
| --------------- | ------------- | -------------------------------------- | ----------------------------------------------------------------------------- |
| `id`            | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                           |
| `run_id`        | TEXT          | NOT NULL, UNIQUE                       | Format: `PAY-RUN/YYYY/MM`. Generated by Go                                    |
| `month`         | TEXT          | NOT NULL, UNIQUE                       | Format: `YYYY-MM`. UNIQUE enforces one payroll run per month                  |
| `pay_date`      | DATE          | NULL                                   | Actual or expected payment date                                               |
| `status`        | TEXT          | NOT NULL, DEFAULT 'draft'              | One of: `draft`, `approved`, `paid`                                           |
| `total_net_pay` | NUMERIC(15,2) | NOT NULL, DEFAULT 0                    | Cache. Sum of all payroll_lines.net_pay. Computed by Go when run is finalised |
| `created_by`    | TEXT          | NOT NULL                               | Staff ID of Admin who initiated the run                                       |
| `approved_by`   | TEXT          | NULL                                   | Staff ID of Super Admin who approved                                          |
| `approved_at`   | TIMESTAMPTZ   | NULL                                   | Set by Go when approved, UTC                                                  |
| `paid_at`       | TIMESTAMPTZ   | NULL                                   | Set by Go when marked paid, UTC                                               |
| `created_at`    | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                |
| `updated_at`    | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                                |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `run_id`
- `UNIQUE` on `month`
- Index on `status`

**Business rules:**

- Go blocks approval if any attendance_daily.is_exception = true
  exists without supervisor approval for that month
- Marking paid triggers one cash_transactions outflow for total_net_pay
- Only active staff are included
- Both login and no-login staff are included in payroll

---

### Table: `payroll_lines`

**Purpose:**
One row per staff member per payroll run. Full breakdown of earnings
and deductions. Auto-populated by Go from attendance data. Can be
manually adjusted before approval with a mandatory logged reason.

| Column                    | Type          | Constraints                                        | Notes                                                                             |
| ------------------------- | ------------- | -------------------------------------------------- | --------------------------------------------------------------------------------- |
| `id`                      | UUID          | PRIMARY KEY, DEFAULT gen_random_uuid()             | Internal identifier                                                               |
| `run_id`                  | UUID          | NOT NULL, FK → `payroll_runs.id` ON DELETE CASCADE | Parent payroll run                                                                |
| `staff_id`                | TEXT          | NOT NULL, FK → `staff.staff_id`                    | The staff member this line covers                                                 |
| `base_salary`             | NUMERIC(15,2) | NOT NULL                                           | Monthly gross base pulled from staff record at run creation time                  |
| `overtime_hours_weekday`  | NUMERIC(8,2)  | NOT NULL, DEFAULT 0                                | Total approved weekday overtime hours for the month                               |
| `overtime_hours_saturday` | NUMERIC(8,2)  | NOT NULL, DEFAULT 0                                | Total approved Saturday overtime hours for the month                              |
| `overtime_rate_weekday`   | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Rate per hour for weekday overtime from system_settings at run time               |
| `overtime_rate_saturday`  | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Rate per hour for Saturday overtime from system_settings at run time              |
| `overtime_earnings`       | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Computed by Go: (weekday hours x weekday rate) + (saturday hours x saturday rate) |
| `night_shift_count`       | INTEGER       | NOT NULL, DEFAULT 0                                | Number of approved night shifts for the month                                     |
| `sunday_shift_count`      | INTEGER       | NOT NULL, DEFAULT 0                                | Number of approved Sunday shifts for the month                                    |
| `night_shift_earnings`    | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Computed: night_shift_count x night_shift_rate from system_settings               |
| `sunday_shift_earnings`   | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Computed: sunday_shift_count x sunday_shift_rate from system_settings             |
| `allowances`              | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Additional payments e.g. bonuses                                                  |
| `allowances_notes`        | TEXT          | NULL                                               | Description of allowances                                                         |
| `lateness_deduction`      | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Computed by Go from late_minutes in attendance_daily                              |
| `absence_deduction`       | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Computed by Go from absent days in attendance_daily                               |
| `other_deductions`        | NUMERIC(15,2) | NOT NULL, DEFAULT 0                                | Loan repayments or other manual deductions                                        |
| `other_deductions_notes`  | TEXT          | NULL                                               | Description of other deductions                                                   |
| `net_pay`                 | NUMERIC(15,2) | NOT NULL                                           | Computed by Go: base + overtime + night + sunday + allowances - deductions        |
| `is_manually_adjusted`    | BOOLEAN       | NOT NULL, DEFAULT false                            | True if any value was manually edited after auto-population                       |
| `adjustment_reason`       | TEXT          | NULL                                               | Mandatory when is_manually_adjusted = true                                        |
| `created_at`              | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                            | Set by Go, UTC                                                                    |
| `updated_at`              | TIMESTAMPTZ   | NOT NULL, DEFAULT now()                            | Set by Go on every update, UTC                                                    |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `run_id`
- Index on `staff_id`
- `UNIQUE` on `(run_id, staff_id)` — one line per staff per run

**Business rules:**

- Go auto-populates all lines from attendance_daily and system_settings
- Only approved overtime minutes are used — observed overtime ignored
- net_pay always recomputed by Go — never accepted from client
- Manual adjustments require is_manually_adjusted = true and adjustment_reason
- Lines are locked once the parent run is approved

---

## Relationships in Group 10

```
payroll_runs (1) ──< payroll_lines (many)
  One payroll run has many lines — one per active staff member.
  Each line belongs to exactly one run.
  UNIQUE on (run_id, staff_id) enforces one line per staff per run.

staff (1) ──< payroll_lines (many)
  One staff member has one payroll line per run.
  Each line belongs to exactly one staff member.

payroll_runs ──> cash_transactions (Group 9)
  When run status = paid, Go creates exactly one outflow
  cash_transaction for total_net_pay.
  Loose reference via linked_entity_id in cash_transactions.

payroll_lines ....> attendance_daily (Group 11)
  Loose reference — payroll lines are populated from attendance data
  but do not FK into attendance_daily. Go reads attendance at run
  creation time and bakes values into the payroll line.
```

---

---

## GROUP 11 — Attendance

**Tables:** `attendance_daily`, `overtime_requests`

Attendance data comes exclusively from the T500 biometric device via
Excel import. Go parses and populates attendance_daily automatically.
Observed overtime creates overtime_requests which must be approved
before minutes flow into payroll. Unresolved exceptions block payroll approval.

---

### Table: `attendance_daily`

**Purpose:**
One row per staff member per working day. Populated by Go from T500
Excel import. Idempotent — re-importing the same file updates existing
records. Any exception flags block payroll approval until resolved.

| Column                       | Type        | Constraints                            | Notes                                                                            |
| ---------------------------- | ----------- | -------------------------------------- | -------------------------------------------------------------------------------- |
| `id`                         | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                              |
| `staff_id`                   | TEXT        | NOT NULL, FK → `staff.staff_id`        | The staff member this record belongs to                                          |
| `date`                       | DATE        | NOT NULL                               | The working day this record covers                                               |
| `clock_in`                   | TIMESTAMPTZ | NULL                                   | Earliest clock-in event for the day from T500 export                             |
| `clock_out`                  | TIMESTAMPTZ | NULL                                   | Latest clock-out event for the day from T500 export                              |
| `minutes_worked`             | INTEGER     | NOT NULL, DEFAULT 0                    | Computed by Go: clock_out - clock_in in minutes                                  |
| `late_minutes`               | INTEGER     | NOT NULL, DEFAULT 0                    | Computed: max(0, clock_in - (scheduled_start + grace_period))                    |
| `observed_overtime_weekday`  | INTEGER     | NOT NULL, DEFAULT 0                    | max(0, clock_out - scheduled_end) on Mon-Fri. In minutes                         |
| `observed_overtime_saturday` | INTEGER     | NOT NULL, DEFAULT 0                    | max(0, clock_out - saturday_end) on Saturday. In minutes                         |
| `is_absent`                  | BOOLEAN     | NOT NULL, DEFAULT false                | True if no clock-in exists for this day in T500 export                           |
| `is_exception`               | BOOLEAN     | NOT NULL, DEFAULT false                | True if any exception condition is met. Blocks payroll approval                  |
| `exception_reason`           | TEXT        | NULL                                   | Description e.g. missing clock-out, unusually long shift                         |
| `exception_resolved`         | BOOLEAN     | NOT NULL, DEFAULT false                | Set to true by supervisor when exception is reviewed                             |
| `exception_resolved_by`      | TEXT        | NULL                                   | Staff ID of supervisor who resolved the exception                                |
| `exception_resolved_at`      | TIMESTAMPTZ | NULL                                   | Set by Go when resolved, UTC                                                     |
| `is_manual_edit`             | BOOLEAN     | NOT NULL, DEFAULT false                | True if any value manually corrected after import. Auto-sets is_exception = true |
| `manual_edit_reason`         | TEXT        | NULL                                   | Mandatory when is_manual_edit = true                                             |
| `created_at`                 | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                   |
| `updated_at`                 | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                                   |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `(staff_id, date)` — idempotent import key
- Index on `staff_id`
- Index on `date`
- Index on `is_exception`
- Index on `is_absent`

**Exception conditions — Go flags is_exception = true when:**

```
1. clock_out is NULL (missing clock-out)
2. minutes_worked > configured max shift duration
3. Multiple conflicting clock events on the same day
4. Staff ID in T500 does not match any active staff record
5. is_manual_edit = true
```

---

### Table: `overtime_requests`

**Purpose:**
Created automatically by Go when T500 import produces observed overtime

> 0. Supervisor approves full, partial, or zero minutes. Only approved
>    minutes flow into payroll calculations — observed minutes are for
>    reference only.

| Column                      | Type        | Constraints                            | Notes                                                    |
| --------------------------- | ----------- | -------------------------------------- | -------------------------------------------------------- |
| `id`                        | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                      |
| `staff_id`                  | TEXT        | NOT NULL, FK → `staff.staff_id`        | The staff member with observed overtime                  |
| `attendance_daily_id`       | UUID        | NOT NULL, FK → `attendance_daily.id`   | The attendance record this request is linked to          |
| `date`                      | DATE        | NOT NULL                               | The day overtime was observed                            |
| `observed_minutes_weekday`  | INTEGER     | NOT NULL, DEFAULT 0                    | Observed weekday overtime minutes from attendance_daily  |
| `observed_minutes_saturday` | INTEGER     | NOT NULL, DEFAULT 0                    | Observed Saturday overtime minutes from attendance_daily |
| `approved_minutes_weekday`  | INTEGER     | NOT NULL, DEFAULT 0                    | Approved weekday minutes. Set by reviewer                |
| `approved_minutes_saturday` | INTEGER     | NOT NULL, DEFAULT 0                    | Approved Saturday minutes. Set by reviewer               |
| `status`                    | TEXT        | NOT NULL, DEFAULT 'pending'            | One of: `pending`, `approved`, `partial`, `rejected`     |
| `reviewed_by`               | TEXT        | NULL                                   | Staff ID of reviewer                                     |
| `reviewed_at`               | TIMESTAMPTZ | NULL                                   | Set by Go when reviewed, UTC                             |
| `review_notes`              | TEXT        | NULL                                   | Mandatory when status = partial or rejected              |
| `created_at`                | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                           |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `(staff_id, date)` — one request per staff per day
- Index on `staff_id`
- Index on `date`
- Index on `status`
- Index on `attendance_daily_id`

---

## Relationships in Group 11

```
staff (1) ──< attendance_daily (many)
  One staff member has one attendance record per working day.
  UNIQUE on (staff_id, date) enforces one record per staff per day.

attendance_daily (1) ──── overtime_requests (1)
  One attendance record can have at most one overtime request.
  UNIQUE on (staff_id, date) in overtime_requests enforces this.

staff (1) ──< overtime_requests (many)
  One staff member can have overtime requests across many days.
  Each request belongs to exactly one staff member.

overtime_requests ....> payroll_lines (Group 10)
  Loose reference — Go reads approved overtime minutes when creating
  payroll lines. No FK stored on payroll_lines.
```

---

---

## GROUP 12 — Tasks & Errands

**Tables:** `tasks`, `checklist_items`, `errands`, `performance_flags`

Task and duty management covering daily assignments, completion tracking,
offsite errand approvals, and automated performance flagging. Does not
affect payroll, inventory, or financial records directly.

---

### Table: `tasks`

**Purpose:**
A duty assigned to a staff member for a specific date. Supervisors
assign tasks to their team. No-login staff have tasks managed by
their supervisor on their behalf.

| Column        | Type        | Constraints                            | Notes                                                 |
| ------------- | ----------- | -------------------------------------- | ----------------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                   |
| `task_id`     | TEXT        | NOT NULL, UNIQUE                       | Format: `TSK/YYYY/MM/DD/001`. Generated by Go         |
| `task_date`   | DATE        | NOT NULL                               | The date this task is assigned for                    |
| `assigned_to` | TEXT        | NOT NULL, FK → `staff.staff_id`        | The staff member responsible                          |
| `assigned_by` | TEXT        | NOT NULL, FK → `staff.staff_id`        | The login-enabled user who assigned the task          |
| `department`  | TEXT        | NOT NULL                               | One of: `factory`, `admin`                            |
| `title`       | TEXT        | NOT NULL                               | Short task description                                |
| `description` | TEXT        | NULL                                   | Optional longer instructions                          |
| `status`      | TEXT        | NOT NULL, DEFAULT 'assigned'           | One of: `assigned`, `in_progress`, `done`, `not_done` |
| `is_offsite`  | BOOLEAN     | NOT NULL, DEFAULT false                | When true, Go creates a linked errand record          |
| `due_date`    | DATE        | NOT NULL                               | Deadline. Defaults to task_date for daily duties      |
| `notes`       | TEXT        | NULL                                   | Blockers or updates added during execution            |
| `created_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                        |
| `updated_at`  | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                        |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `task_id`
- Index on `assigned_to`
- Index on `assigned_by`
- Index on `task_date`
- Index on `status`
- Index on `department`

---

### Table: `checklist_items`

**Purpose:**
Sub-steps of a task. Completing all items automatically advances
parent task status to done.

| Column         | Type        | Constraints                                 | Notes                                   |
| -------------- | ----------- | ------------------------------------------- | --------------------------------------- |
| `id`           | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid()      | Internal identifier                     |
| `task_id`      | UUID        | NOT NULL, FK → `tasks.id` ON DELETE CASCADE | Parent task                             |
| `description`  | TEXT        | NOT NULL                                    | The checklist step                      |
| `is_completed` | BOOLEAN     | NOT NULL, DEFAULT false                     | Updated by supervisor or staff member   |
| `completed_by` | TEXT        | NULL                                        | Staff ID of who marked complete         |
| `completed_at` | TIMESTAMPTZ | NULL                                        | Set by Go when is_completed = true, UTC |
| `created_at`   | TIMESTAMPTZ | NOT NULL, DEFAULT now()                     | Set by Go, UTC                          |

**Indexes:**

- `PRIMARY KEY` on `id`
- Index on `task_id`

---

### Table: `errands`

**Purpose:**
Tracks offsite movements. Created automatically when task.is_offsite = true.
Factory errands approved by Production Supervisor, office errands by Admin.

| Column           | Type        | Constraints                            | Notes                                                                  |
| ---------------- | ----------- | -------------------------------------- | ---------------------------------------------------------------------- |
| `id`             | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                    |
| `errand_id`      | TEXT        | NOT NULL, UNIQUE                       | Format: `ERR/YYYY/MM/DD/001`. Generated by Go                          |
| `task_id`        | UUID        | NULL, FK → `tasks.id`                  | Linked task. NULL if errand created independently                      |
| `staff_id`       | TEXT        | NOT NULL, FK → `staff.staff_id`        | The staff member going offsite                                         |
| `department`     | TEXT        | NOT NULL                               | One of: `factory`, `admin`. Determines approver                        |
| `purpose`        | TEXT        | NOT NULL                               | Reason for leaving                                                     |
| `destination`    | TEXT        | NOT NULL                               | Where the staff member is going                                        |
| `status`         | TEXT        | NOT NULL, DEFAULT 'requested'          | One of: `requested`, `approved`, `rejected`, `in_transit`, `completed` |
| `requested_by`   | TEXT        | NOT NULL                               | Staff ID of supervisor who created the request                         |
| `approved_by`    | TEXT        | NULL                                   | Staff ID of approver                                                   |
| `approved_at`    | TIMESTAMPTZ | NULL                                   | Set by Go when approved, UTC                                           |
| `approval_notes` | TEXT        | NULL                                   | Mandatory when rejected. Optional when approved                        |
| `time_out`       | TIMESTAMPTZ | NULL                                   | When the staff member left the premises                                |
| `time_returned`  | TIMESTAMPTZ | NULL                                   | When the staff member returned                                         |
| `created_at`     | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                         |
| `updated_at`     | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go on every update, UTC                                         |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `errand_id`
- Index on `task_id`
- Index on `staff_id`
- Index on `status`
- Index on `department`

---

### Table: `performance_flags`

**Purpose:**
Created by nightly Go job when a staff member accumulates 3 or more
not_done tasks in a calendar month. Informational only — does not
automatically affect payroll.

| Column                  | Type        | Constraints                            | Notes                                   |
| ----------------------- | ----------- | -------------------------------------- | --------------------------------------- |
| `id`                    | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                     |
| `staff_id`              | TEXT        | NOT NULL, FK → `staff.staff_id`        | The staff member who triggered the flag |
| `month`                 | TEXT        | NOT NULL                               | Format: `YYYY-MM`                       |
| `incomplete_task_count` | INTEGER     | NOT NULL                               | Number of not_done tasks this month     |
| `status`                | TEXT        | NOT NULL, DEFAULT 'open'               | One of: `open`, `reviewed`              |
| `reviewed_by`           | TEXT        | NULL                                   | Staff ID of reviewer                    |
| `reviewed_at`           | TIMESTAMPTZ | NULL                                   | Set by Go when reviewed, UTC            |
| `created_at`            | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go nightly job, UTC              |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `(staff_id, month)` — one flag per staff per month
- Index on `staff_id`
- Index on `status`
- Index on `month`

---

## Relationships in Group 12

```
staff (1) ──< tasks (many) [assigned_to]
  One staff member can have many tasks assigned to them.
  Each task is assigned to exactly one staff member.

staff (1) ──< tasks (many) [assigned_by]
  One staff member can assign many tasks.
  Each task is assigned by exactly one staff member.

tasks (1) ──< checklist_items (many)
  One task can have many checklist items.
  Each checklist item belongs to exactly one task.
  Cascade delete on checklist_items if task deleted.

tasks (1) ──── errands (1)
  One task can have at most one errand when is_offsite = true.
  Each errand optionally links to one task.

staff (1) ──< errands (many)
  One staff member can have many errands.
  Each errand belongs to exactly one staff member.

staff (1) ──< performance_flags (many)
  One staff member can have one performance flag per month.
  UNIQUE on (staff_id, month) enforces one flag per staff per month.
```

---

---

## GROUP 13 — Audit Log

**Tables:** `audit_log`

Permanent, append-only record of every significant action in the system.
Written exclusively by Go. No user including Super Admin can delete or
modify entries. RLS policy has no DELETE permission for any role.

---

### Table: `audit_log`

| Column        | Type        | Constraints                            | Notes                                                                                                            |
| ------------- | ----------- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `id`          | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                                              |
| `log_id`      | TEXT        | NOT NULL, UNIQUE                       | Format: `LOG/YYYY/MM/DD/001`. Generated by Go                                                                    |
| `timestamp`   | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Exact UTC time. Set by Go — never client-supplied                                                                |
| `user_id`     | TEXT        | NOT NULL                               | Staff ID of the authenticated user who performed the action                                                      |
| `action`      | TEXT        | NOT NULL                               | One of: `CREATE`, `UPDATE`, `APPROVE`, `REVERSE`, `LOGIN`, `LOGOUT`, `EXPORT`, `REJECT`, `TERMINATE`, `OVERRIDE` |
| `entity_type` | TEXT        | NOT NULL                               | Type of record affected e.g. `invoice`, `bill`, `inventory_transaction`, `payroll_run`                           |
| `entity_id`   | TEXT        | NOT NULL                               | Business ID of the affected record e.g. `INV/2026/03/24/001`. TEXT to accommodate all ID formats                 |
| `before_json` | JSONB       | NULL                                   | Complete snapshot before the change. NULL for CREATE actions                                                     |
| `after_json`  | JSONB       | NULL                                   | Complete snapshot after the change                                                                               |
| `reason`      | TEXT        | NULL                                   | Mandatory for UPDATE, REVERSE, APPROVE, REJECT, TERMINATE, OVERRIDE                                              |
| `ip_address`  | TEXT        | NULL                                   | IP address from the HTTP request                                                                                 |
| `user_agent`  | TEXT        | NULL                                   | Browser or client user agent                                                                                     |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `log_id`
- Index on `user_id`
- Index on `entity_type`
- Index on `entity_id`
- Index on `timestamp`
- Index on `action`

**Business rules:**

- Append only — no UPDATE or DELETE endpoints exist
- RLS policy has no DELETE permission for any role including Super Admin
- audit.Service.Log() called inside the same DB transaction as the
  action being logged — entry and action succeed or fail together
- reason is mandatory for sensitive actions — Go enforces at API level
- Auditor role has read-only access to the full log
- HR Manager has partial access — HR-related entries only

---

## Relationships in Group 13

```
audit_log ....> all entities (loose reference)
  audit_log.entity_id is a TEXT reference to the business ID of any
  record. No FK — audit log must remain intact even if source records
  are restructured. Go resolves the source at query time when needed.

audit_log.user_id ....> staff.staff_id (loose reference)
  References the staff member who performed the action.
  No FK — audit entries must persist even if staff records change.
```

---

---

## GROUP 14 — System Settings

**Tables:** `system_settings`

Stores all configurable business rules and thresholds. Managed
exclusively by Super Admin. Go reads these at runtime — no hardcoded
business rules in application code.

---

### Table: `system_settings`

| Column          | Type        | Constraints                            | Notes                                                                                           |
| --------------- | ----------- | -------------------------------------- | ----------------------------------------------------------------------------------------------- |
| `id`            | UUID        | PRIMARY KEY, DEFAULT gen_random_uuid() | Internal identifier                                                                             |
| `setting_key`   | TEXT        | NOT NULL, UNIQUE                       | Machine-readable key e.g. `overtime_rate_weekday`, `discount_threshold`, `grace_period_minutes` |
| `setting_value` | TEXT        | NOT NULL                               | Value stored as TEXT. Go casts to correct type at read time                                     |
| `data_type`     | TEXT        | NOT NULL                               | One of: `integer`, `decimal`, `text`, `boolean`, `time`. Tells Go how to cast                   |
| `description`   | TEXT        | NOT NULL                               | Human-readable description of what this setting controls                                        |
| `module`        | TEXT        | NOT NULL                               | e.g. `payroll`, `invoicing`, `inventory`, `attendance`, `tasks`, `general`                      |
| `updated_by`    | TEXT        | NULL                                   | Staff ID of Super Admin who last updated this setting                                           |
| `updated_at`    | TIMESTAMPTZ | NULL                                   | Set by Go when updated, UTC                                                                     |
| `created_at`    | TIMESTAMPTZ | NOT NULL, DEFAULT now()                | Set by Go, UTC                                                                                  |

**Indexes:**

- `PRIMARY KEY` on `id`
- `UNIQUE` on `setting_key`
- Index on `module`

**Seeded settings:**

| Setting Key                  | Default      | Module       | Description                                          |
| ---------------------------- | ------------ | ------------ | ---------------------------------------------------- |
| `scheduled_start_time`       | `08:30`      | `attendance` | Official start time                                  |
| `scheduled_end_weekday`      | `18:00`      | `attendance` | Official end time Mon-Fri                            |
| `scheduled_end_saturday`     | `15:00`      | `attendance` | Official end time Saturday                           |
| `grace_period_minutes`       | `10`         | `attendance` | Minutes before lateness is recorded                  |
| `overtime_rounding_minutes`  | `15`         | `attendance` | Round overtime to nearest N minutes                  |
| `max_shift_hours`            | `14`         | `attendance` | Shifts beyond this trigger exception flag            |
| `overtime_rate_weekday`      | `0`          | `payroll`    | NGN per hour weekday overtime                        |
| `overtime_rate_saturday`     | `0`          | `payroll`    | NGN per hour Saturday overtime                       |
| `night_shift_rate`           | `0`          | `payroll`    | NGN per night shift                                  |
| `sunday_shift_rate`          | `0`          | `payroll`    | NGN per Sunday shift                                 |
| `discount_threshold_type`    | `percentage` | `invoicing`  | One of: percentage, fixed                            |
| `discount_threshold_value`   | `10`         | `invoicing`  | Discount above this triggers approval workflow       |
| `backdating_limit_days`      | `7`          | `general`    | Entries older than this require Super Admin approval |
| `po_default_tolerance`       | `5`          | `inventory`  | Default PO quantity variance tolerance percentage    |
| `performance_flag_threshold` | `3`          | `tasks`      | not_done tasks per month before flag is raised       |
| `low_stock_alert_enabled`    | `true`       | `inventory`  | Whether low stock alerts are active                  |

**Business rules:**

- Only Super Admin can update settings
- Every change written to audit_log with before and after values
- Go reads settings at request time — always latest values used
- New settings added via migration without code changes

---

## Relationships in Group 14

```
system_settings has no FK relationships to other tables.
It is a standalone configuration store read by Go at runtime.

system_settings ....> all modules (loose runtime dependency)
  Go reads setting values when processing requests in each module.
  The dependency lives in application code, not the schema.
```

---

## COMPLETE ERD SUMMARY — All Groups

| Group                 | Tables                                                                | Phase |
| --------------------- | --------------------------------------------------------------------- | ----- |
| 1 — Auth & Security   | `roles`, `permissions`, `role_permissions`, `fired_staff`             | 1     |
| 2 — Staff             | `staff`, `user_roles`, `role_history`                                 | 1     |
| 3 — Master Data       | `customers`, `suppliers`, `products`, `product_variants`, `materials` | 1     |
| 4 — Inventory         | `inventory_transactions`, `wac_snapshots`                             | 2     |
| 5 — Procurement       | `purchase_orders`, `purchase_order_lines`, `po_delivery_logs`         | 2     |
| 6 — Payables          | `bills`, `bill_payments`                                              | 2     |
| 7 — Sales             | `orders`, `order_lines`, `invoices`, `customer_payments`              | 3     |
| 8 — Job Costing       | `job_labor`, `job_overhead`                                           | 4     |
| 9 — Cash Transactions | `cash_transactions`                                                   | 4     |
| 10 — Payroll          | `payroll_runs`, `payroll_lines`                                       | 5     |
| 11 — Attendance       | `attendance_daily`, `overtime_requests`                               | 5     |
| 12 — Tasks & Errands  | `tasks`, `checklist_items`, `errands`, `performance_flags`            | 7     |
| 13 — Audit Log        | `audit_log`                                                           | 6     |
| 14 — System Settings  | `system_settings`                                                     | 1     |

**Total tables: 34**

---

_BloansBooks ERD · Complete · Internal Technical Documentation · Confidential_
