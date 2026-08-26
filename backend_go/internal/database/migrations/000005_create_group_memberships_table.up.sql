-- Group memberships table
CREATE TABLE IF NOT EXISTS group_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_group_memberships_user_group UNIQUE (user_id, group_id)
);

CREATE INDEX IF NOT EXISTS idx_group_memberships_user_id ON group_memberships (user_id);
CREATE INDEX IF NOT EXISTS idx_group_memberships_group_id ON group_memberships (group_id);
