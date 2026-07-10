-- Rollback: remove all seeded role_permissions and permissions
DELETE FROM role_permissions
WHERE role_id = (SELECT id FROM roles WHERE name = 'super_admin');

DELETE FROM permissions
WHERE code IN (
  'auth.manage_accounts', 'auth.reset_password', 'auth.manage_roles',
  'staff.create', 'staff.view', 'staff.update', 'staff.deactivate', 'staff.terminate',
  'customers.create', 'customers.view', 'customers.update',
  'suppliers.create', 'suppliers.view', 'suppliers.update',
  'products.create', 'products.view', 'products.update',
  'materials.create', 'materials.view', 'materials.update',
  'inventory.view', 'inventory.adjust', 'inventory.approve_adjustment', 'inventory.issue_to_job',
  'purchase_orders.create', 'purchase_orders.view', 'purchase_orders.approve', 'purchase_orders.deliver', 'purchase_orders.cancel',
  'bills.create', 'bills.view', 'bills.approve', 'bills.reverse',
  'invoices.create', 'invoices.view', 'invoices.post', 'invoices.reverse', 'invoices.cancel', 'invoices.apply_discount',
  'payments.create', 'payments.view',
  'job_costing.view', 'job_costing.add_labor', 'job_costing.add_overhead',
  'payroll.create_run', 'payroll.view', 'payroll.approve', 'payroll.mark_paid',
  'attendance.import', 'attendance.view', 'attendance.approve_exception', 'attendance.approve_overtime',
  'tasks.create', 'tasks.view', 'tasks.update', 'tasks.assign',
  'errands.create', 'errands.view', 'errands.approve',
  'reports.view', 'reports.export',
  'audit.view',
  'settings.view', 'settings.update',
  'periods.close', 'periods.reopen'
);
