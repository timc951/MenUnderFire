DROP INDEX IF EXISTS idx_group_messages_form_id;
ALTER TABLE group_messages DROP COLUMN IF EXISTS form_id;
