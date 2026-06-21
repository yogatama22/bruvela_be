-- ================================================================
-- BRUVELA BAKEHOUSE - PostgreSQL Database Schema
-- ================================================================
-- This schema is designed for PostgreSQL 12+
-- Uses UUID as primary keys with gen_random_uuid()
-- ================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ================================================================
-- TABLE: users
-- ================================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'staff',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster email lookup
CREATE INDEX idx_users_email ON users(email);

-- ================================================================
-- TABLE: products (master menu)
-- ================================================================
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    price INTEGER NOT NULL,
    pcs_per_box INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster code lookup
CREATE INDEX idx_products_code ON products(code);
CREATE INDEX idx_products_status ON products(status);

-- ================================================================
-- TABLE: ingredients (master bahan baku)
-- ================================================================
CREATE TABLE ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    pack_unit VARCHAR(20) NOT NULL,
    qty_per_pack DECIMAL(10,3) NOT NULL,
    use_unit VARCHAR(20) NOT NULL,
    price_per_pack INTEGER NOT NULL,
    price_per_use DECIMAL(10,4),
    min_stock DECIMAL(10,3) DEFAULT 0,
    current_stock DECIMAL(10,3) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for low stock alerts
CREATE INDEX idx_ingredients_stock ON ingredients(current_stock, min_stock);

-- ================================================================
-- TABLE: recipes (komposisi resep per produk)
-- ================================================================
CREATE TABLE recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    qty_per_box DECIMAL(10,3) NOT NULL,
    use_unit VARCHAR(20) NOT NULL,
    cost_per_box DECIMAL(10,4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, ingredient_id)
);

-- Indexes for faster lookups
CREATE INDEX idx_recipes_product ON recipes(product_id);
CREATE INDEX idx_recipes_ingredient ON recipes(ingredient_id);

-- ================================================================
-- TABLE: customers
-- ================================================================
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    location VARCHAR(200),
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster name search
CREATE INDEX idx_customers_name ON customers(name);

-- ================================================================
-- TABLE: batches (periode/batch produksi)
-- ================================================================
CREATE TABLE batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_number INTEGER UNIQUE NOT NULL,
    name VARCHAR(100),
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(20) DEFAULT 'open',
    total_modal INTEGER DEFAULT 0,
    total_revenue INTEGER DEFAULT 0,
    total_hpp INTEGER DEFAULT 0,
    gross_profit INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for batch queries
CREATE INDEX idx_batches_status ON batches(status);
CREATE INDEX idx_batches_number ON batches(batch_number);

-- ================================================================
-- TABLE: orders
-- ================================================================
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES batches(id) ON DELETE RESTRICT,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    customer_name VARCHAR(100),
    order_date DATE NOT NULL,
    channel VARCHAR(30) DEFAULT 'whatsapp',
    shipping_type VARCHAR(20),
    shipping_dest VARCHAR(200),
    shipping_cost INTEGER DEFAULT 0,
    discount INTEGER DEFAULT 0,
    total_product INTEGER DEFAULT 0,
    total_bill INTEGER DEFAULT 0,
    pay_status VARCHAR(20) DEFAULT 'belum_bayar',
    prod_status VARCHAR(20) DEFAULT 'baru',
    note TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for filtering and sorting
CREATE INDEX idx_orders_batch ON orders(batch_id);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_date ON orders(order_date);
CREATE INDEX idx_orders_pay_status ON orders(pay_status);
CREATE INDEX idx_orders_prod_status ON orders(prod_status);
CREATE INDEX idx_orders_channel ON orders(channel);

-- ================================================================
-- TABLE: order_items
-- ================================================================
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_code VARCHAR(20),
    product_name VARCHAR(100),
    qty_box INTEGER NOT NULL,
    price_per_box INTEGER NOT NULL,
    subtotal INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for order items lookup
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_product ON order_items(product_id);

-- ================================================================
-- TABLE: ingredient_purchases (pembelian bahan)
-- ================================================================
CREATE TABLE ingredient_purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES batches(id) ON DELETE RESTRICT,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    purchase_date DATE NOT NULL,
    supplier VARCHAR(100),
    qty_pack DECIMAL(10,3) NOT NULL,
    price_per_pack INTEGER NOT NULL,
    total_price INTEGER NOT NULL,
    note TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for purchase queries
CREATE INDEX idx_purchases_batch ON ingredient_purchases(batch_id);
CREATE INDEX idx_purchases_ingredient ON ingredient_purchases(ingredient_id);
CREATE INDEX idx_purchases_date ON ingredient_purchases(purchase_date);

-- ================================================================
-- TABLE: stock_logs (log pergerakan stok)
-- ================================================================
CREATE TABLE stock_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    batch_id UUID NOT NULL REFERENCES batches(id) ON DELETE RESTRICT,
    log_type VARCHAR(20) NOT NULL,
    qty DECIMAL(10,3) NOT NULL,
    stock_before DECIMAL(10,3),
    stock_after DECIMAL(10,3),
    reference_id UUID,
    reference_type VARCHAR(20),
    note TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for stock log queries
CREATE INDEX idx_stock_logs_ingredient ON stock_logs(ingredient_id);
CREATE INDEX idx_stock_logs_batch ON stock_logs(batch_id);
CREATE INDEX idx_stock_logs_type ON stock_logs(log_type);
CREATE INDEX idx_stock_logs_date ON stock_logs(created_at);

-- ================================================================
-- TABLE: journal_entries (jurnal keuangan)
-- ================================================================
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID NOT NULL REFERENCES batches(id) ON DELETE RESTRICT,
    entry_date DATE NOT NULL,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    amount INTEGER NOT NULL,
    balance INTEGER DEFAULT 0,
    partner VARCHAR(100),
    reference_id UUID,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for journal queries
CREATE INDEX idx_journal_batch ON journal_entries(batch_id);
CREATE INDEX idx_journal_type ON journal_entries(type);
CREATE INDEX idx_journal_date ON journal_entries(entry_date);
CREATE INDEX idx_journal_partner ON journal_entries(partner);

-- ================================================================
-- TRIGGERS FOR UPDATED_AT
-- ================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ingredients_updated_at BEFORE UPDATE ON ingredients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_recipes_updated_at BEFORE UPDATE ON recipes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_batches_updated_at BEFORE UPDATE ON batches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================
-- SEED DATA - Default Admin User
-- ================================================================
-- Password: admin123 (hashed with bcrypt)
-- You should change this password after first login
INSERT INTO users (name, email, password, role) VALUES
('Admin', 'admin@bruvela.com', '$2a$10$YourHashedPasswordHere', 'admin');

-- ================================================================
-- SEED DATA - Sample Batch
-- ================================================================
INSERT INTO batches (batch_number, name, start_date, status) VALUES
(1, 'Batch 1', CURRENT_DATE, 'open');

-- ================================================================
-- VIEWS FOR ANALYTICS
-- ================================================================

-- View: Low Stock Ingredients
CREATE OR REPLACE VIEW v_low_stock_ingredients AS
SELECT 
    id,
    name,
    current_stock,
    min_stock,
    use_unit,
    (min_stock - current_stock) as shortage
FROM ingredients
WHERE current_stock < min_stock
ORDER BY (min_stock - current_stock) DESC;

-- View: Order Summary by Batch
CREATE OR REPLACE VIEW v_batch_order_summary AS
SELECT 
    b.id as batch_id,
    b.batch_number,
    b.name as batch_name,
    COUNT(o.id) as total_orders,
    SUM(o.total_bill) as total_revenue,
    SUM(CASE WHEN o.pay_status = 'lunas' THEN o.total_bill ELSE 0 END) as paid_amount,
    SUM(CASE WHEN o.pay_status != 'lunas' THEN o.total_bill ELSE 0 END) as pending_amount
FROM batches b
LEFT JOIN orders o ON b.id = o.batch_id
GROUP BY b.id, b.batch_number, b.name
ORDER BY b.batch_number DESC;

-- View: Product Sales Summary
CREATE OR REPLACE VIEW v_product_sales AS
SELECT 
    p.id as product_id,
    p.code,
    p.name,
    p.price,
    COUNT(oi.id) as times_ordered,
    SUM(oi.qty_box) as total_boxes_sold,
    SUM(oi.subtotal) as total_revenue
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id
GROUP BY p.id, p.code, p.name, p.price
ORDER BY total_boxes_sold DESC NULLS LAST;

-- ================================================================
-- COMMENTS FOR DOCUMENTATION
-- ================================================================

COMMENT ON TABLE users IS 'User accounts for system access';
COMMENT ON TABLE products IS 'Product/menu master data';
COMMENT ON TABLE ingredients IS 'Ingredient/raw material master data';
COMMENT ON TABLE recipes IS 'Recipe composition linking products and ingredients';
COMMENT ON TABLE customers IS 'Customer master data';
COMMENT ON TABLE batches IS 'Production batch/period tracking';
COMMENT ON TABLE orders IS 'Customer orders';
COMMENT ON TABLE order_items IS 'Line items for each order';
COMMENT ON TABLE ingredient_purchases IS 'Ingredient purchase transactions';
COMMENT ON TABLE stock_logs IS 'Audit trail for stock movements';
COMMENT ON TABLE journal_entries IS 'Financial journal entries';

-- ================================================================
-- END OF SCHEMA
-- ================================================================
