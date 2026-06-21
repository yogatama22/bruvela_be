-- Create shipping_type table
CREATE TABLE IF NOT EXISTS shipping_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipping_code VARCHAR(20) UNIQUE NOT NULL,
    shipping_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on shipping_code
CREATE INDEX IF NOT EXISTS idx_shipping_type_code ON shipping_type(shipping_code);
