# Suppliers API

**Base path:** `/api/v1/suppliers`

The suppliers module manages the businesses and individuals the company buys from — raw-material vendors, printers, logistics providers, artisans, utilities, and others. Records are created by staff and always carry the identity of their creator.

All endpoints require authentication. Permission requirements are noted per endpoint.

### A note on validation

Required string fields are trimmed of surrounding whitespace before use, and empty or whitespace-only values are rejected:

- **On create**, a missing or blank required field (`name`, `phone`, `category`) is rejected with the message `<field> is required`.
- **On update**, a field that is present but blank is rejected with `<field> cannot be empty`.

_These validation failures currently return `500`. Mapping them to `400` is a planned cross-module change; the exact messages are listed per endpoint below._

---

## Data Models

### Supplier Object (full)

Returned by `GET /suppliers/:id` and as each item in the list response.

| Field | Type | Description |
|---|---|---|
| `id` | UUID | Internal identifier |
| `supplierId` | string | Business ID e.g. `SPL-0001` |
| `name` | string | Supplier or business name |
| `phone` | string | Primary phone number |
| `email` | string \| null | Contact email |
| `address` | string \| null | Address |
| `notes` | string \| null | Free-form notes |
| `category` | string | One of: `raw_materials`, `printing`, `logistics`, `artisans`, `utilities`, `other` |
| `currency` | string | ISO currency code — `NGN` in v1 |
| `status` | string | One of: `active`, `inactive` |
| `createdBy` | StaffSummary | Compact profile of the staff member who created the record |
| `createdAt` | string | ISO 8601 timestamp |
| `updatedAt` | string | ISO 8601 timestamp |

### Supplier Summary Object

The compact shape returned inside the create and update responses.

| Field | Type | Description |
|---|---|---|
| `supplierId` | string | Business ID |
| `name` | string | Supplier name |
| `phone` | string | Phone number |
| `email` | string \| null | Contact email — omitted when null |
| `category` | string | One of the six category values |
| `status` | string | `active` or `inactive` |

---

## Endpoints

### List All Suppliers

Returns a filtered, sorted, paginated list of supplier records. Each item is a full Supplier Object.

```
GET /api/v1/suppliers
```

**Authentication:** Required.
**Permission:** `suppliers.view`

**Query Parameters**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `limit` | integer | `20` | Number of records to return |
| `offset` | integer | `0` | Number of records to skip |
| `search` | string | — | ILIKE match against `name`, `supplierId`, `phone` |
| `category` | string | — | Filter by category: `raw_materials`, `printing`, `logistics`, `artisans`, `utilities`, `other` |
| `status` | string | — | Filter by status: `active` or `inactive` |
| `sortBy` | string | `createdAt` | Column to sort by. Allowed values: `createdAt`, `name`, `supplierId`, `category`, `status` |
| `sortOrder` | string | `desc` | `asc` or `desc` |

**Example requests**

```
GET /api/v1/suppliers?category=raw_materials&status=active&limit=20&offset=0
GET /api/v1/suppliers?search=leather&sortBy=name&sortOrder=asc
GET /api/v1/suppliers?status=inactive&limit=10
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Suppliers fetched successfully",
  "statusCode": 200,
  "data": {
    "data_items": [
      {
        "id": "018e1c2d-...",
        "supplierId": "SPL-0001",
        "name": "Lagos Leather & Hides Co",
        "phone": "+2348012345678",
        "email": "supply@lagosleather.example",
        "address": "45 Balogun Street, Lagos",
        "notes": null,
        "category": "raw_materials",
        "currency": "NGN",
        "status": "active",
        "createdBy": {
          "staffId": "BLN-0002",
          "firstName": "Funmi",
          "lastName": "Adeyemi",
          "email": "funmi@example.com",
          "phone": "+2348055001122",
          "department": "procurement",
          "jobTitle": "Procurement Officer",
          "status": "active"
        },
        "createdAt": "2026-08-01T09:00:00Z",
        "updatedAt": "2026-08-01T09:00:00Z"
      }
    ],
    "count": 1,
    "total_count": 42,
    "limit": 20,
    "offset": 0
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `401` | `Authentication required` | No valid token |
| `403` | `You do not have permission to perform this action` | Missing `suppliers.view` |

---

### Get Supplier by ID

Retrieve the full profile of a single supplier.

```
GET /api/v1/suppliers/:id
```

**Authentication:** Required.
**Permission:** `suppliers.view`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The supplier's internal `id` |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Supplier fetched successfully",
  "statusCode": 200,
  "data": {
    "id": "018e1c2d-...",
    "supplierId": "SPL-0001",
    "name": "Lagos Leather & Hides Co",
    "phone": "+2348012345678",
    "email": "supply@lagosleather.example",
    "address": "45 Balogun Street, Lagos",
    "notes": "Delivers hides every second Tuesday.",
    "category": "raw_materials",
    "currency": "NGN",
    "status": "active",
    "createdBy": {
      "staffId": "BLN-0002",
      "firstName": "Funmi",
      "lastName": "Adeyemi",
      "email": "funmi@example.com",
      "phone": "+2348055001122",
      "department": "procurement",
      "jobTitle": "Procurement Officer",
      "status": "active"
    },
    "createdAt": "2026-08-01T09:00:00Z",
    "updatedAt": "2026-08-01T09:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | `:id` is not a valid UUID |
| `404` | `Supplier not found` | No record found for the given ID |

---

### Create Supplier

Create a new supplier. The system auto-generates a `supplierId`. The creator is taken from the caller's authenticated identity.

```
POST /api/v1/suppliers
```

**Authentication:** Required.
**Permission:** `suppliers.create`

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | Yes | Supplier or business name |
| `phone` | string | Yes | Primary phone number |
| `email` | string | No | Contact email |
| `address` | string | No | Address |
| `notes` | string | No | Free-form notes |
| `category` | string | Yes | One of: `raw_materials`, `printing`, `logistics`, `artisans`, `utilities`, `other` |

> `currency` and `status` are not accepted here. `currency` is fixed to `NGN` in v1, and `status` defaults to `active` — both are set by the database.

```json
{
  "name": "PrintWorks Nigeria",
  "phone": "+2348055667788",
  "email": "jobs@printworks.example",
  "address": "12 Aba Road, Port Harcourt",
  "notes": "Handles box and label printing.",
  "category": "printing"
}
```

**Success Response — `201 Created`**

```json
{
  "status": "success",
  "message": "Supplier created successfully",
  "statusCode": 201,
  "data": {
    "id": "018e1c2d-...",
    "supplier": {
      "supplierId": "SPL-0002",
      "name": "PrintWorks Nigeria",
      "phone": "+2348055667788",
      "email": "jobs@printworks.example",
      "category": "printing",
      "status": "active"
    },
    "createdAt": "2026-08-25T11:30:00Z",
    "updatedAt": "2026-08-25T11:30:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | Malformed JSON body |
| `401` | `Authentication required` | No valid token, or caller identity is missing |
| `403` | `You do not have permission to perform this action` | Missing `suppliers.create` |
| `500` | `name is required` / `phone is required` / `category is required` | A required field was empty or whitespace-only |
| `500` | `failed to create supplier: ...` | Database error or constraint violation (e.g. invalid `category`) |

---

### Update Supplier

Partially update a supplier record. Only fields included in the request body are changed.

```
PATCH /api/v1/suppliers/:id
```

**Authentication:** Required.
**Permission:** `suppliers.update`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The supplier's internal `id` |

**Request Body**

All fields are optional. Include only what you want to change.

| Field | Type | Description |
|---|---|---|
| `name` | string | Supplier name |
| `phone` | string | Phone number |
| `email` | string | Contact email |
| `address` | string | Address |
| `notes` | string | Free-form notes |
| `category` | string | One of the six category values |
| `status` | string | `active` or `inactive` |

> `supplierId`, `currency`, and `createdBy` are immutable and cannot be changed here.

```json
{
  "phone": "+2348090011223",
  "status": "inactive"
}
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Supplier updated successfully",
  "statusCode": 200,
  "data": {
    "id": "018e1c2d-...",
    "supplier": {
      "supplierId": "SPL-0002",
      "name": "PrintWorks Nigeria",
      "phone": "+2348090011223",
      "email": "jobs@printworks.example",
      "category": "printing",
      "status": "inactive"
    },
    "createdAt": "2026-08-25T11:30:00Z",
    "updatedAt": "2026-08-25T14:00:00Z"
  }
}
```

**Error Responses**

| Status | Message | Cause |
|---|---|---|
| `400` | `Invalid request data` | `:id` is not a valid UUID, or body is malformed |
| `401` | `Authentication required` | No valid token |
| `403` | `You do not have permission to perform this action` | Missing `suppliers.update` |
| `404` | `Supplier not found` | No record found for the given ID |
| `500` | `name cannot be empty` / `phone cannot be empty` / `category cannot be empty` | A provided field was whitespace-only |
| `500` | `failed to update supplier: ...` | Database error |

---

## Business Rules

- `supplierId` is auto-generated at creation in the format `{PREFIX}-NNNN` (prefix from the `SUPPLIER_ID_PREFIX` setting) e.g. `SPL-0001`. It never changes.
- Suppliers are never deleted — there is no delete endpoint and no `suppliers.delete` permission. Use `status: inactive` to deactivate a supplier.
- `category` must be one of `raw_materials`, `printing`, `logistics`, `artisans`, `utilities`, or `other`, enforced by a database CHECK constraint.
- `currency` is fixed to `NGN` in v1 (database default + CHECK) and is not accepted in create or update payloads.
- `status` defaults to `active` on creation.
- `createdBy` is taken from the authenticated caller's JWT identity — never from the request body — and is immutable.
- Every supplier response embeds the creator as a compact staff summary, resolved via a `LEFT JOIN` on the staff table.
