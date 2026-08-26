ALTER TABLE users DROP COLUMN IF EXISTS agreement_accepted_at;
ALTER TABLE users DROP COLUMN IF EXISTS agreement_version;
ALTER TABLE users DROP COLUMN IF EXISTS agreement_signature;
ALTER TABLE users DROP COLUMN IF EXISTS agreement_ip;
ALTER TABLE users DROP COLUMN IF EXISTS agreement_user_agent;
