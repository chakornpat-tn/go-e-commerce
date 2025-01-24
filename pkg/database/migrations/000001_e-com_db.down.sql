BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_users_table ON "users";
DROP TRIGGER IF EXISTS set_updated_at_oauth_table ON "oauth";
DROP TRIGGER IF EXISTS set_updated_at_products_table ON "products";
DROP TRIGGER IF EXISTS set_updated_at_orders_table ON "orders";
DROP TRIGGER IF EXISTS set_updated_at_images_table ON "images";

DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "oauth" CASCADE;
DROP TABLE IF EXISTS "roles" CASCADE;
DROP TABLE IF EXISTS "products" CASCADE;
DROP TABLE IF EXISTS "products_categories" CASCADE;
DROP TABLE IF EXISTS "categories" CASCADE;
DROP TABLE IF EXISTS "images" CASCADE;
DROP TABLE IF EXISTS "orders" CASCADE;
DROP TABLE IF EXISTS "products_orders" CASCADE;
 
DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS products_id_seq;
DROP SEQUENCE IF EXISTS order_id_seq;

DROP TYPE IF EXITS "order_status";


COMMIT;