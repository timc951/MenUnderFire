-- Update the role constraint to use uppercase role names and add LEADER role
ALTER TABLE group_memberships DROP CONSTRAINT group_memberships_role_check;
ALTER TABLE group_memberships ADD CONSTRAINT group_memberships_role_check
    CHECK (role IN ('OWNER', 'LEADER', 'MEMBER'));

-- Update any existing lowercase roles to uppercase
UPDATE group_memberships SET role = UPPER(role);

-- Update the default value
ALTER TABLE group_memberships ALTER COLUMN role SET DEFAULT 'MEMBER';
