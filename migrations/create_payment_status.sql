-- Create payment_status table
CREATE TABLE IF NOT EXISTS payment_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status_code VARCHAR(20) UNIQUE NOT NULL,
    status_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on status_code
CREATE INDEX IF NOT EXISTS idx_payment_status_code ON payment_status(status_code);
