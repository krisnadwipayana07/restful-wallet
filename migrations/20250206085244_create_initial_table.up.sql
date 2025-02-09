CREATE TABLE IF NOT EXISTS "wallet_table" (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	wallet_name VARCHAR(255) NOT NULL,
	wallet_curr_balance NUMERIC(36, 18),
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "transaction_table" (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	wallet_id BIGINT NOT NULL,
	trc_type SMALLINT NOT NULL,
	trc_is_debit BOOL NOT NULL,
	trc_value NUMERIC(36, 18) NOT NULL,
	trc_remarks TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS "idx_transaction_table_wallet_id" ON "transaction_table" (wallet_id);