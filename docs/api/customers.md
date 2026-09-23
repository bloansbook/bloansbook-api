# Customers API

**Base path:** `/api/v1/customers`

The customers module manages the business's customer records — the people and organizations that loans and transactions are attached to. Customers are either individuals (`retail`) or businesses (`corporate`). Records are created by staff and always carry the identity of their creator.

All endpoints require authentication. Permission requirements are noted per endpoint.

### A note on validation

Required string fields are trimmed of surrounding whitespace before use, and empty or whitespace-only values are rejected:

- **On create**, a missing or blank required field (`name`, `phone`, `type`) is rejected with the message `<field> is required`.
- **On update**, a field that is present but blank is rejected with `<field> cannot be empty`.

_These validation failures currently return `500`. Mapping them to `400` is a planned cross-module change; the exact messages are listed per endpoint below._

---

## Data Models

### Customer Object (full)

Returned by `GET /customers/:id` and as each item in the list response.

| Field | Type | Description |
|---|---|---|
| `id` | UUID | Internal identifier |
| `customerId` | string | Business ID e.g. `CUS-0001` |
| `name` | string | Customer full name or business name |
| `phone` | string | Primary phone number |
| `email` | string \| null | Contact email |
| `address` | string \| null | Address |
| `notes` | string \| null | Free-form notes |
| `type` | string | One of: `retail`, `corporate` |
| `currency` | string | ISO currency code — `NGN` in v1 |
| `status` | string | One of: `active`, `inactive` |
| `createdBy` | StaffSummary | Compact profile of the staff member who created the record |
| `createdAt` | string | ISO 8601 timestamp |
| `updatedAt` | string | ISO 8601 timestamp |

### Customer Summary Object

The compact shape returned inside the create and update responses.

| Field | Type | Description |
|---|---|---|
| `customerId` | string | Business ID |
| `name` | string | Customer name |
| `phone` | string | Phone number |
| `email` | string \| null | Contact email — omitted when null |
| `type` | string | `retail` or `corporate` |
| `status` | string | `active` or `inactive` |

---

## Endpoints

### List All Customers

Returns a filtered, sorted, paginated list of customer records. Each item is a full Customer Object.

```
GET /api/v1/customers
```

**Authentication:** Required.
**Permission:** `customers.view`

**Query Parameters**

| Parameter | Type | Default | Description |
|---|---|---|---|
| `limit` | integer | `20` | Number of records to return |
| `offset` | integer | `0` | Number of records to skip |
| `search` | string | — | ILIKE match against `name`, `customerId`, `phone` |
| `type` | string | — | Filter by type: `retail` or `corporate` |
| `status` | string | — | Filter by status: `active` or `inactive` |
| `sortBy` | string | `createdAt` | Column to sort by. Allowed values: `createdAt`, `name`, `customerId`, `type`, `status` |
| `sortOrder` | string | `desc` | `asc` or `desc` |

**Example requests**

```
GET /api/v1/customers?type=corporate&status=active&limit=20&offset=0
GET /api/v1/customers?search=adewale&sortBy=name&sortOrder=asc
GET /api/v1/customers?status=inactive&limit=10
```

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Customers fetched successfully",
  "statusCode": 200,
  "data": {
    "data_items": [
      {
        "id": "018e1c2d-...",
        "customerId": "CUS-0001",
        "name": "Adewale Ventures Ltd",
        "phone": "+2348012345678",
        "email": "accounts@adewale.example",
        "address": "45 Balogun Street, Lagos",
        "notes": null,
        "type": "corporate",
        "currency": "NGN",
        "status": "active",
        "createdBy": {
          "staffId": "BLN-0002",
          "firstName": "Funmi",
          "lastName": "Adeyemi",
          "email": "funmi@example.com",
          "phone": "+2348055001122",
          "department": "sales",
          "jobTitle": "Sales Executive",
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
| `403` | `You do not have permission to perform this action` | Missing `customers.view` |

---

### Get Customer by ID

Retrieve the full profile of a single customer.

```
GET /api/v1/customers/:id
```

**Authentication:** Required.
**Permission:** `customers.view`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The customer's internal `id` |

**Success Response — `200 OK`**

```json
{
  "status": "success",
  "message": "Customer fetched successfully",
  "statusCode": 200,
  "data": {
    "id": "018e1c2d-...",
    "customerId": "CUS-0001",
    "name": "Adewale Ventures Ltd",
    "phone": "+2348012345678",
    "email": "accounts@adewale.example",
    "address": "45 Balogun Street, Lagos",
    "notes": "Prefers invoices by email.",
    "type": "corporate",
    "currency": "NGN",
    "status": "active",
    "createdBy": {
      "staffId": "BLN-0002",
      "firstName": "Funmi",
      "lastName": "Adeyemi",
      "email": "funmi@example.com",
      "phone": "+2348055001122",
      "department": "sales",
      "jobTitle": "Sales Executive",
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
| `404` | `Customer not found` | No record found for the given ID |

---

### Create Customer

Create a new customer. The system auto-generates a `customerId`. The creator is taken from the caller's authenticated identity.

```
POST /api/v1/customers
```

**Authentication:** Required.
**Permission:** `customers.create`

**Request Body**

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | Yes | Customer full name or business name |
| `phone` | string | Yes | Primary phone number |
| `email` | string | No | Contact email |
| `address` | string | No | Address |
| `notes` | string | No | Free-form notes |
| `type` | string | Yes | One of: `retail`, `corporate` |

> `currency` and `status` are not accepted here. `currency` is fixed to `NGN` in v1, and `status` defaults to `active` — both are set by the database.

```json
{
  "name": "Chinelo Retail Stores",
  "phone": "+2348055667788",
  "email": "chinelo@example.com",
  "address": "12 Aba Road, Port Harcourt",
  "notes": "Walk-in customer, referred by BLN-0002.",
  "type": "retail"
}
```

**Success Response — `201 Created`**

```json
{
  "status": "success",
  "message": "Customer created successfully",
  "statusCode": 201,
  "data": {
    "id": "018e1c2d-...",
    "customer": {
      "customerId": "CUS-0002",
      "name": "Chinelo Retail Stores",
      "phone": "+2348055667788",
      "email": "chinelo@example.com",
      "type": "retail",
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
| `403` | `You do not have permission to perform this action` | Missing `customers.create` |
| `500` | `name is required` / `phone is required` / `type is required` | A required field was empty or whitespace-only |
| `500` | `failed to create customer: ...` | Database error or constraint violation (e.g. invalid `type`) |

---

### Update Customer

Partially update a customer record. Only fields included in the request body are changed.

```
PATCH /api/v1/customers/:id
```

**Authentication:** Required.
**Permission:** `customers.update`

**Path Parameters**

| Parameter | Type | Description |
|---|---|---|
| `id` | UUID | The customer's internal `id` |

**Request Body**

All fields are optional. Include only what you want to change.

| Field | Type | Description |
|---|---|---|
| `name` | string | Customer name |
| `phone` | string | Phone number |
| `email` | string | Contact email |
| `address` | string | Address |
| `notes` | string | Free-form notes |
| `type` | string | `retail` or `corporate` |
| `status` | string | `active` or `inactive` |

> `customerId`, `currency`, and `createdBy` are immutable and cannot be changed here.

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
  "message": "Customer updated successfully",
  "statusCode": 200,
  "data": {
    "id": "018e1c2d-...",
    "customer": {
      "customerId": "CUS-0002",
      "name": "Chinelo Retail Stores",
      "phone": "+2348090011223",
      "email": "chinelo@example.com",
      "type": "retail",
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
| `403` | `You do not have permission to perform this action` | Missing `customers.update` |
| `404` | `Customer not found` | No record found for the given ID |
| `500` | `name cannot be empty` / `phone cannot be empty` / `type cannot be empty` | A provided field was whitespace-only |
| `500` | `failed to update customer: ...` | Database error |

---

## Business Rules

- `customerId` is auto-generated at creation in the format `{PREFIX}-NNNN` (prefix from the `CUSTOMER_ID_PREFIX` setting) e.g. `CUS-0001`. It never changes.
- Customers are never deleted — there is no delete endpoint and no `customers.delete` permission. Use `status: inactive` to deactivate a customer.
- `type` must be one of `retail` or `corporate`, enforced by a database CHECK constraint.
- `currency` is fixed to `NGN` in v1 (database default + CHECK) and is not accepted in create or update payloads.
- `status` defaults to `active` on creation.
- `createdBy` is taken from the authenticated caller's JWT identity — never from the request body — and is immutable.
- Every customer response embeds the creator as a compact staff summary, resolved via a `LEFT JOIN` on the staff table.
