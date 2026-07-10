-- =============================================================================
-- Seed: permissions + role_permissions
-- Every permission code defined in the ERD is inserted here.
-- role_permissions assigns ALL permissions to super_admin only.
-- Other roles start with no permissions — assigned by Super Admin via the UI.
-- =============================================================================

-- -----------------------------------------------------------------------------
-- PERMISSIONS
-- -----------------------------------------------------------------------------
INSERT INTO permissions (code, module, description) VALUES

  -- auth
  ('auth.manage_accounts',  'auth', 'Create, disable, and manage staff login accounts'),
  ('auth.reset_password',   'auth', 'Reset a staff member''s login password'),
  ('auth.manage_roles',     'auth', 'Create roles, assign and revoke permissions'),

  -- staff
  ('staff.create',          'staff', 'Create a new staff record'),
  ('staff.view',            'staff', 'View staff records and profiles'),
  ('staff.update',          'staff', 'Edit staff record details'),
  ('staff.deactivate',      'staff', 'Set a staff member''s status to inactive'),
  ('staff.terminate',       'staff', 'Terminate a staff member and record a fired_staff entry'),

  -- customers
  ('customers.create',      'customers', 'Create a new customer record'),
  ('customers.view',        'customers', 'View customer records'),
  ('customers.update',      'customers', 'Edit customer record details'),

  -- suppliers
  ('suppliers.create',      'suppliers', 'Create a new supplier record'),
  ('suppliers.view',        'suppliers', 'View supplier records'),
  ('suppliers.update',      'suppliers', 'Edit supplier record details'),

  -- products
  ('products.create',       'products', 'Create a new product and its variants'),
  ('products.view',         'products', 'View products and variants'),
  ('products.update',       'products', 'Edit product and variant details'),

  -- materials
  ('materials.create',      'materials', 'Create a new raw material record'),
  ('materials.view',        'materials', 'View raw material records'),
  ('materials.update',      'materials', 'Edit raw material details'),

  -- inventory
  ('inventory.view',              'inventory', 'View inventory transactions and stock levels'),
  ('inventory.adjust',            'inventory', 'Submit inventory adjustment requests'),
  ('inventory.approve_adjustment','inventory', 'Approve or reject pending inventory adjustments'),
  ('inventory.issue_to_job',      'inventory', 'Issue materials to a production job'),

  -- purchase_orders
  ('purchase_orders.create',  'purchase_orders', 'Create a new purchase order'),
  ('purchase_orders.view',    'purchase_orders', 'View purchase orders'),
  ('purchase_orders.approve', 'purchase_orders', 'Approve and send a purchase order to a supplier'),
  ('purchase_orders.deliver', 'purchase_orders', 'Record delivery of a purchase order'),
  ('purchase_orders.cancel',  'purchase_orders', 'Cancel a purchase order with a mandatory reason'),

  -- bills
  ('bills.create',   'bills', 'Create a new supplier bill'),
  ('bills.view',     'bills', 'View supplier bills'),
  ('bills.approve',  'bills', 'Approve a supplier bill and trigger inventory update'),
  ('bills.reverse',  'bills', 'Reverse a posted supplier bill'),

  -- invoices
  ('invoices.create',          'invoices', 'Create a new customer invoice'),
  ('invoices.view',            'invoices', 'View customer invoices'),
  ('invoices.post',            'invoices', 'Post and lock a finalised invoice'),
  ('invoices.reverse',         'invoices', 'Reverse a posted invoice'),
  ('invoices.cancel',          'invoices', 'Cancel a draft or unpaid invoice'),
  ('invoices.apply_discount',  'invoices', 'Apply a discount to an invoice'),

  -- payments
  ('payments.create', 'payments', 'Record a customer or supplier payment'),
  ('payments.view',   'payments', 'View payment records'),

  -- job_costing
  ('job_costing.view',         'job_costing', 'View job cost summaries and reports'),
  ('job_costing.add_labor',    'job_costing', 'Add a labour entry to a job'),
  ('job_costing.add_overhead', 'job_costing', 'Add an overhead entry to a job'),

  -- payroll
  ('payroll.create_run', 'payroll', 'Initiate a new monthly payroll run'),
  ('payroll.view',       'payroll', 'View payroll runs and payslips'),
  ('payroll.approve',    'payroll', 'Approve a payroll run'),
  ('payroll.mark_paid',  'payroll', 'Mark a payroll run as paid and trigger cash transaction'),

  -- attendance
  ('attendance.import',            'attendance', 'Import T500 attendance data from Excel'),
  ('attendance.view',              'attendance', 'View attendance records'),
  ('attendance.approve_exception', 'attendance', 'Resolve attendance exception flags'),
  ('attendance.approve_overtime',  'attendance', 'Approve or reject overtime requests'),

  -- tasks
  ('tasks.create', 'tasks', 'Create a new task assignment'),
  ('tasks.view',   'tasks', 'View tasks and their status'),
  ('tasks.update', 'tasks', 'Update task status or details'),
  ('tasks.assign', 'tasks', 'Assign tasks to staff members'),

  -- errands
  ('errands.create',  'errands', 'Create an offsite errand request'),
  ('errands.view',    'errands', 'View errand records'),
  ('errands.approve', 'errands', 'Approve or reject an errand request'),

  -- reports
  ('reports.view',   'reports', 'View system reports'),
  ('reports.export', 'reports', 'Export reports to PDF or Excel'),

  -- audit
  ('audit.view', 'audit', 'View the system audit log'),

  -- settings
  ('settings.view',   'settings', 'View system settings'),
  ('settings.update', 'settings', 'Update system settings'),

  -- periods
  ('periods.close',  'periods', 'Close an accounting period'),
  ('periods.reopen', 'periods', 'Reopen a closed accounting period with Super Admin authorisation')

ON CONFLICT (code) DO NOTHING;

-- -----------------------------------------------------------------------------
-- ROLE PERMISSIONS — super_admin gets every permission
-- All other roles start empty; Super Admin assigns them via the UI.
-- -----------------------------------------------------------------------------
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    r.id,
    p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'super_admin'
ON CONFLICT DO NOTHING;
