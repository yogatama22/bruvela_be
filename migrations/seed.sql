-- ================================================================
-- BRUVELA BAKEHOUSE - Seed Data
-- ================================================================
-- Sample data untuk testing dan development
-- ================================================================

-- ================================================================
-- PAYMENT STATUSES - Status Pembayaran
-- ================================================================
INSERT INTO payment_status (status_code, status_name) VALUES
('belum_bayar', 'Belum Bayar'),
('dp', 'DP'),
('lunas', 'Lunas')
ON CONFLICT (status_code) DO NOTHING;

-- ================================================================
-- SHIPPING TYPES - Jenis Pengiriman
-- ================================================================
INSERT INTO shipping_type (shipping_code, shipping_name) VALUES
('dalam_kota', 'Dalam Kota'),
('luar_kota', 'Luar Kota'),
('jnt', 'JNE/J&T'),
('gojek', 'Gojek/Grab')
ON CONFLICT (shipping_code) DO NOTHING;

-- ================================================================
-- USERS - Admin dan Staff
-- ================================================================
-- Password untuk semua user: admin123
-- Hash: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy

INSERT INTO users (name, email, password, role) VALUES
('Admin Bruvela', 'admin@bruvela.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin'),
('Staff Aul', 'aul@bruvela.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'staff'),
('Staff Dhavinna', 'dhavinna@bruvela.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'staff')
ON CONFLICT (email) DO NOTHING;

-- ================================================================
-- PRODUCTS - Menu Bruvela
-- ================================================================
INSERT INTO products (code, name, price, pcs_per_box, status) VALUES
('BRV-CLASSIC', 'Bruv Classic', 20000, 6, 'active'),
('BRV-ALMOND', 'Bruv Almond', 22000, 6, 'active'),
('BRV-CHEESE', 'Bruv Cheese', 22000, 6, 'active'),
('BRV-COOKIES', 'Bruv Cookies', 23000, 6, 'active'),
('BRV-CLASSIC-LG', 'Bruv Classic Large', 80000, 24, 'active'),
('BRV-ALMOND-LG', 'Bruv Almond Large', 88000, 24, 'active'),
('BRV-CHEESE-LG', 'Bruv Cheese Large', 88000, 24, 'active'),
('BRV-COOKIES-LG', 'Bruv Cookies Large', 92000, 24, 'active'),
('BRV-MIX-LG', 'Bruv Mix Large', 87000, 24, 'active')
ON CONFLICT (code) DO NOTHING;

-- ================================================================
-- INGREDIENTS - Bahan Baku
-- ================================================================
INSERT INTO ingredients (name, pack_unit, qty_per_pack, use_unit, price_per_pack, min_stock, current_stock) VALUES
('Tepung Terigu', 'kg', 1000, 'gram', 15000, 2000, 5000),
('Gula Pasir', 'kg', 1000, 'gram', 14000, 500, 3000),
('Butter/Margarin', 'kg', 1000, 'gram', 45000, 200, 1000),
('Telur', 'kg', 20, 'butir', 30000, 10, 40),
('Baking Powder', 'pack', 100, 'gram', 8000, 50, 200),
('Vanilla Extract', 'botol', 50, 'ml', 25000, 20, 100),
('Susu Bubuk', 'pack', 500, 'gram', 35000, 200, 800),
('Cocoa Powder', 'pack', 250, 'gram', 40000, 100, 500),
('Almond Slice', 'pack', 250, 'gram', 55000, 100, 300),
('Keju Parut', 'pack', 200, 'gram', 35000, 100, 400),
('Chocochips', 'pack', 250, 'gram', 30000, 100, 500),
('Oreo Cookies', 'pack', 10, 'pcs', 15000, 5, 20),
('Garam', 'pack', 500, 'gram', 5000, 100, 500),
('Box Kecil', 'pack', 50, 'pcs', 75000, 10, 100),
('Box Besar', 'pack', 25, 'pcs', 100000, 5, 20),
('Kertas Roti', 'pack', 100, 'lembar', 20000, 10, 50),
('Plastik Wrapping', 'roll', 100, 'meter', 30000, 5, 15),
('Sticker Label', 'pack', 100, 'pcs', 25000, 20, 80)
ON CONFLICT DO NOTHING;

-- ================================================================
-- BATCHES - Production Batches
-- ================================================================
INSERT INTO batches (batch_number, name, start_date, status, total_modal) VALUES
(1, 'Batch 1', '2024-01-01', 'closed', 1578000),
(2, 'Batch 2', '2024-02-01', 'closed', 1059476),
(3, 'Batch 3', CURRENT_DATE, 'open', 978000)
ON CONFLICT (batch_number) DO NOTHING;

-- ================================================================
-- CUSTOMERS - Sample Customers
-- ================================================================
INSERT INTO customers (name, phone, location) VALUES
('Teto', '081234567890', 'Jakarta Selatan'),
('Budi Santoso', '081234567891', 'Jakarta Pusat'),
('Siti Nurhaliza', '081234567892', 'Tangerang'),
('Ahmad Rizki', '081234567893', 'Bekasi'),
('Dewi Lestari', '081234567894', 'Depok')
ON CONFLICT DO NOTHING;

-- ================================================================
-- RECIPES - Sample Recipe untuk Bruv Classic
-- ================================================================
-- Catatan: Anda perlu mengambil UUID dari products dan ingredients
-- yang sudah di-insert sebelumnya. Ini adalah contoh template.
-- Sesuaikan dengan UUID aktual dari database Anda.

-- Contoh untuk Bruv Classic (per 1 box = 6 pcs)
-- INSERT INTO recipes (product_id, ingredient_id, qty_per_box, use_unit) VALUES
-- ((SELECT id FROM products WHERE code = 'BRV-CLASSIC'), 
--  (SELECT id FROM ingredients WHERE name = 'Tepung Terigu'), 
--  150, 'gram'),
-- ((SELECT id FROM products WHERE code = 'BRV-CLASSIC'), 
--  (SELECT id FROM ingredients WHERE name = 'Gula Pasir'), 
--  80, 'gram'),
-- ((SELECT id FROM products WHERE code = 'BRV-CLASSIC'), 
--  (SELECT id FROM ingredients WHERE name = 'Butter/Margarin'), 
--  100, 'gram'),
-- ((SELECT id FROM products WHERE code = 'BRV-CLASSIC'), 
--  (SELECT id FROM ingredients WHERE name = 'Telur'), 
--  2, 'butir');

-- ================================================================
-- JOURNAL ENTRIES - Sample Financial Data
-- ================================================================
-- INSERT INTO journal_entries (batch_id, entry_date, description, type, amount, partner) VALUES
-- ((SELECT id FROM batches WHERE batch_number = 1),
--  '2024-01-01', 'Modal Awal - Aul', 'modal', 789000, 'Aul'),
-- ((SELECT id FROM batches WHERE batch_number = 1),
--  '2024-01-01', 'Modal Awal - Dhavinna', 'modal', 789000, 'Dhavinna');

-- ================================================================
-- END OF SEED DATA
-- ================================================================

-- Untuk melihat data yang sudah di-insert:
SELECT 'Payment Statuses:' as table_name, COUNT(*) as count FROM payment_status
UNION ALL
SELECT 'Shipping Types:' as table_name, COUNT(*) as count FROM shipping_type
UNION ALL
SELECT 'Users:' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Products:', COUNT(*) FROM products
UNION ALL
SELECT 'Ingredients:', COUNT(*) FROM ingredients
UNION ALL
SELECT 'Batches:', COUNT(*) FROM batches
UNION ALL
SELECT 'Customers:', COUNT(*) FROM customers;
