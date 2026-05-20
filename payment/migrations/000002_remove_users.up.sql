ALTER TABLE devices      DROP CONSTRAINT devices_user_id_fkey;
ALTER TABLE subscriptions DROP CONSTRAINT subscriptions_user_id_fkey;
ALTER TABLE payments     DROP CONSTRAINT payments_user_id_fkey;
ALTER TABLE documents    DROP CONSTRAINT documents_user_id_fkey;

DROP TABLE users;