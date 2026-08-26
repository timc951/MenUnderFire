-- Add expiration to invite codes
ALTER TABLE groups ADD COLUMN invite_code_expires_at TIMESTAMP WITH TIME ZONE;

-- Set default expiration for existing invite codes to 30 days from now
UPDATE groups SET invite_code_expires_at = NOW() + INTERVAL '30 days';
