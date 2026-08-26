-- Add form_id column to group_messages to support sending forms as messages
ALTER TABLE group_messages ADD COLUMN form_id UUID REFERENCES forms(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_group_messages_form_id ON group_messages(form_id);
