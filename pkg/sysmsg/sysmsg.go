package sysmsg

// General
const (
	NotFound        = "Resource not found"
	BadRequest      = "Invalid request data"
	Unauthorized    = "Authentication required"
	Forbidden       = "You do not have permission to perform this action"
	InternalError   = "An unexpected error occurred. Please try again"
	ValidationError = "Validation failed"
	Success         = "Operation successful"
)

// Environment Variables
const (
	NoEnvFile = "No .env file found"
)

// Database
const (
	CannotConnect        = "Unable to connect to database"
	CannotPing           = "Unable to ping database"
	ConnectionSuccessful = "Database connection successful"
	ConnectionClosed     = "Database connection closed"
)

// Auth
const (
	LoginSuccess       = "Login successful"
	LogoutSuccess      = "Logout successful"
	NoLoginAccess      = "Staff does not have login access"
	InvalidCredentials = "Invalid Staff ID or password"
	AccessDenied       = "Access denied"
	TokenExpired       = "Session expired. Please log in again"
	TokenInvalid       = "Invalid session token"
	SetupComplete      = "Super Admin account created successfully"
	SetupAlreadyDone   = "System setup has already been completed"
)

// Staff
const (
	StaffCreated      = "Staff record created successfully"
	StaffFetched      = "Staff record fetched successfully"
	StaffListFetched  = "Staff records fetched successfully"
	StaffUpdated      = "Staff record updated successfully"
	StaffDeactivated  = "Staff member deactivated successfully"
	StaffNotFound     = "Staff record not found"
	StaffIDExists     = "A staff record with this ID already exists"
	StaffTerminated   = "Staff member has been terminated"
	StaffAlreadyFired = "This staff member has already been terminated"
)

// Roles & Permissions
const (
	RoleCreated              = "Role created successfully"
	RoleFetched              = "Role fetched successfully"
	RoleListFetched          = "Roles fetched successfully"
	RoleUpdated              = "Role updated successfully"
	RoleNotFound             = "Role not found"
	RoleAssigned             = "Role assigned successfully"
	RoleRevoked              = "Role revoked successfully"
	PermissionAssignedToRole = "Permission assigned to role successfully"
	PermissionNotFound       = "Permission not found"
	PermissionCreated        = "Permission created successfully"
	PermissionFetched        = "Permission fetched successfully"
	CannotModifySuperAdmin   = "The Super Admin role cannot be modified"
)

// Customers
const (
	CustomerCreated     = "Customer created successfully"
	CustomerFetched     = "Customer fetched successfully"
	CustomerListFetched = "Customers fetched successfully"
	CustomerUpdated     = "Customer updated successfully"
	CustomerNotFound    = "Customer not found"
)

// Suppliers
const (
	SupplierCreated     = "Supplier created successfully"
	SupplierFetched     = "Supplier fetched successfully"
	SupplierListFetched = "Suppliers fetched successfully"
	SupplierUpdated     = "Supplier updated successfully"
	SupplierNotFound    = "Supplier not found"
)

// Products
const (
	ProductCreated     = "Product created successfully"
	ProductFetched     = "Product fetched successfully"
	ProductListFetched = "Products fetched successfully"
	ProductUpdated     = "Product updated successfully"
	ProductNotFound    = "Product not found"
)

// Materials
const (
	MaterialCreated     = "Material created successfully"
	MaterialFetched     = "Material fetched successfully"
	MaterialListFetched = "Materials fetched successfully"
	MaterialUpdated     = "Material updated successfully"
	MaterialNotFound    = "Material not found"
	LowStockAlert       = "Material stock is at or below reorder level"
)

// Inventory
const (
	InventoryTxnCreated  = "Inventory transaction created successfully"
	InventoryTxnFetched  = "Inventory transaction fetched successfully"
	InventoryTxnNotFound = "Inventory transaction not found"
	InventoryTxnReversed = "Inventory transaction reversed successfully"
	AdjustmentPending    = "Adjustment submitted and awaiting Super Admin approval"
	WACUpdated           = "Weighted average cost updated successfully"
)

// Purchase Orders
const (
	POCreated         = "Purchase order created successfully"
	POFetched         = "Purchase order fetched successfully"
	POListFetched     = "Purchase orders fetched successfully"
	POUpdated         = "Purchase order updated successfully"
	PONotFound        = "Purchase order not found"
	POApproved        = "Purchase order approved and sent to supplier"
	PODelivered       = "Purchase order marked as delivered"
	POClosed          = "Purchase order closed successfully"
	POCancelled       = "Purchase order cancelled successfully"
	POAlreadyApproved = "This purchase order has already been approved"
	POEditNotAllowed  = "Approved purchase orders cannot be edited without Super Admin authorisation"
)

// Bills
const (
	BillCreated        = "Supplier bill created successfully"
	BillFetched        = "Supplier bill fetched successfully"
	BillListFetched    = "Supplier bills fetched successfully"
	BillNotFound       = "Supplier bill not found"
	BillApproved       = "Supplier bill approved and inventory updated"
	BillReversed       = "Supplier bill reversed successfully"
	BillNoPO           = "Inventory bills must be linked to a purchase order"
	BillPONotDelivered = "Inventory cannot be updated until the purchase order is marked as delivered"
)

// Sales Orders
const (
	OrderCreated     = "Sales order created successfully"
	OrderFetched     = "Sales order fetched successfully"
	OrderListFetched = "Sales orders fetched successfully"
	OrderUpdated     = "Sales order updated successfully"
	OrderNotFound    = "Sales order not found"
	OrderApproved    = "Sales order approved successfully"
	OrderClosed      = "Sales order closed successfully"
)

// Invoices
const (
	InvoiceCreated         = "Invoice created successfully"
	InvoiceFetched         = "Invoice fetched successfully"
	InvoiceListFetched     = "Invoices fetched successfully"
	InvoiceNotFound        = "Invoice not found"
	InvoicePosted          = "Invoice posted successfully"
	InvoiceReversed        = "Invoice reversed successfully"
	InvoiceCancelled       = "Invoice cancelled successfully"
	InvoiceAlreadyPaid     = "This invoice has already been fully paid"
	InvoiceDiscountPending = "Discount exceeds the configured limit. Awaiting Super Admin approval"
)

// Payments
const (
	PaymentCreated        = "Payment recorded successfully"
	PaymentFetched        = "Payment fetched successfully"
	PaymentListFetched    = "Payments fetched successfully"
	PaymentNotFound       = "Payment not found"
	PaymentExceedsBalance = "Payment amount exceeds the outstanding invoice balance"
)

// Job Costing
const (
	JobCostFetched      = "Job cost summary fetched successfully"
	JobLaborAdded       = "Labour entry added successfully"
	JobLaborNotFound    = "Labour entry not found"
	JobOverheadAdded    = "Overhead entry added successfully"
	JobOverheadNotFound = "Overhead entry not found"
)

// Payroll
const (
	PayrollRunCreated          = "Payroll run created successfully"
	PayrollRunFetched          = "Payroll run fetched successfully"
	PayrollRunApproved         = "Payroll run approved successfully"
	PayrollRunPaid             = "Payroll run marked as paid"
	PayrollRunNotFound         = "Payroll run not found"
	PayrollBlockedByExceptions = "Payroll run cannot be approved. Unresolved attendance exceptions exist"
	PayslipGenerated           = "Payslip generated successfully"
)

// Attendance
const (
	AttendanceImported          = "Attendance data imported successfully"
	AttendanceFetched           = "Attendance records fetched successfully"
	AttendanceNotFound          = "Attendance record not found"
	AttendanceExceptionResolved = "Attendance exception resolved successfully"
	AttendanceFileInvalid       = "Uploaded file does not match the expected T500 export format"
	OvertimeRequestCreated      = "Overtime request created successfully"
	OvertimeApproved            = "Overtime approved successfully"
	OvertimeRejected            = "Overtime request rejected"
)

// Tasks
const (
	TaskCreated     = "Task created successfully"
	TaskFetched     = "Task fetched successfully"
	TaskListFetched = "Tasks fetched successfully"
	TaskUpdated     = "Task updated successfully"
	TaskNotFound    = "Task not found"
	TaskCompleted   = "Task marked as complete"
	ErrandCreated   = "Errand request created successfully"
	ErrandApproved  = "Errand approved successfully"
	ErrandRejected  = "Errand request rejected"
	ErrandNotFound  = "Errand not found"
	ErrandCompleted = "Errand marked as complete"
)

// Audit
const (
	AuditLogFetched = "Audit log fetched successfully"
)

// Approvals
const (
	ApprovalPending  = "Request submitted and awaiting Super Admin approval"
	ApprovalApproved = "Request approved successfully"
	ApprovalRejected = "Request rejected"
	ReasonRequired   = "A reason is required for this action"
)

// Period
const (
	PeriodLocked        = "Accounting period locked successfully"
	PeriodUnlocked      = "Accounting period unlocked successfully"
	PeriodAlreadyLocked = "This accounting period is locked. Super Admin authorisation required to make changes"
)

// Dashboard & Reports
const (
	DashboardFetched = "Dashboard data fetched successfully"
	ReportGenerated  = "Report generated successfully"
)
