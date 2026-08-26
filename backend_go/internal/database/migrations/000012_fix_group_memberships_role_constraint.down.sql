-- Revert the role constraint back to lowercase
ALTER TABLE group_memberships DROP CONSTRAINT group_memberships_role_check;

-- Update any uppercase roles to lowercase
UPDATE group_memberships SET role = LOWER(role);

-- Re-add the original constraint (without LEADER)
ALTER TABLE group_memberships ADD CONSTRAINT group_memberships_role_check
    CHECK (role IN ('owner', 'member'));

-- Update the default value
ALTER TABLE group_memberships ALTER COLUMN role SET DEFAULT 'member';
