CREATE TABLE IF NOT EXISTS customers (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    customer_id TEXT UNIQUE NOT NULL,

    name TEXT NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email TEXT,
    address TEXT,

    notes TEXT,
    type TEXT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'NGN',

    created_by UUID NOT NULL REFERENCES staff(id),

    status TEXT NOT NULL DEFAULT 'active',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT customer_type_check
    CHECK (type IN ('retail', 'corporate')),

    CONSTRAINT customer_currency_check 
    CHECK (currency = 'NGN')
);

CREATE INDEX idx_customers_customer_id ON customers(customer_id);
CREATE INDEX idx_customers_name ON customers(name);

CREATE TABLE IF NOT EXISTS suppliers (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),

    supplier_id TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT,
    address TEXT,

    currency TEXT NOT NULL DEFAULT 'NGN',
    category TEXT NOT NULL,

    notes TEXT,

    created_by UUID NOT NULL REFERENCES staff(id),

    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT supplier_currency_check
    CHECK (currency = 'NGN'),

    CONSTRAINT supplier_category_check
    CHECK (category IN (
        'raw_materials',
        'printing',
        'logistics',
        'artisans',
        'utilities',
        'other'
    ))
);

CREATE INDEX idx_suppliers_name ON suppliers(name);
CREATE INDEX idx_suppliers_category ON suppliers(category);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_id TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,
    description TEXT,

    status TEXT NOT NULL DEFAULT 'active',

    created_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT product_status_check
    CHECK (status IN ('active', 'inactive'))
);

CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_name ON products(name);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    variant_id TEXT NOT NULL UNIQUE,

    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    sku TEXT NOT NULL UNIQUE,

    size TEXT,
    color TEXT,

    attributes JSONB,

    selling_price NUMERIC(15,2) NOT NULL DEFAULT 0,

    current_quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    current_wac NUMERIC(15,2) NOT NULL DEFAULT 0,
    current_stock_value NUMERIC(15,2) NOT NULL DEFAULT 0,

    status TEXT NOT NULL DEFAULT 'active',

    created_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT variant_status_check
    CHECK (status IN ('active', 'inactive')),

    CONSTRAINT non_negative_values
    CHECK (
        selling_price >= 0 AND
        current_quantity >= 0 AND
        current_wac >= 0 AND
        current_stock_value >= 0
    )
);

CREATE INDEX idx_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_variants_status ON product_variants(status);

CREATE TABLE materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    material_id TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,

    unit_of_measure TEXT NOT NULL,

    reorder_level NUMERIC(15,2) NOT NULL DEFAULT 0,

    current_quantity NUMERIC(15,2) NOT NULL DEFAULT 0,
    current_wac NUMERIC(15,2) NOT NULL DEFAULT 0,
    current_stock_value NUMERIC(15,2) NOT NULL DEFAULT 0,

    status TEXT NOT NULL DEFAULT 'active',

    created_by UUID NOT NULL REFERENCES staff(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT material_status_check
    CHECK (status IN ('active', 'inactive')),

    CONSTRAINT non_negative_materials
    CHECK (
        reorder_level >= 0 AND
        current_quantity >= 0 AND
        current_wac >= 0 AND
        current_stock_value >= 0
    )
);

CREATE INDEX idx_materials_status ON materials(status);
CREATE INDEX idx_materials_quantity ON materials(current_quantity);