-- Add group settings columns
ALTER TABLE groups ADD COLUMN IF NOT EXISTS require_post_approval BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE groups ADD COLUMN IF NOT EXISTS allow_anonymous_posts BOOLEAN NOT NULL DEFAULT false;
