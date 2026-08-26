-- Add region/state column for Cloudflare geo headers
ALTER TABLE page_hits ADD COLUMN region VARCHAR(200);

-- Add user_email column for tracking who visited (denormalized for anonymous display)
ALTER TABLE page_hits ADD COLUMN user_email VARCHAR(255);
