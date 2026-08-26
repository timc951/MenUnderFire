-- Add MODERATOR to the group membership role constraint
ALTER TABLE group_memberships DROP CONSTRAINT group_memberships_role_check;
ALTER TABLE group_memberships ADD CONSTRAINT group_memberships_role_check
    CHECK (role IN ('OWNER', 'LEADER', 'MODERATOR', 'MEMBER'));
