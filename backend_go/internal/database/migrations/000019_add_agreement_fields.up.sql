-- Add testing agreement acceptance fields to users table
ALTER TABLE users ADD COLUMN agreement_accepted_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN agreement_version VARCHAR(20);
ALTER TABLE users ADD COLUMN agreement_signature VARCHAR(128);
ALTER TABLE users ADD COLUMN agreement_ip VARCHAR(45);
ALTER TABLE users ADD COLUMN agreement_user_agent TEXT;
