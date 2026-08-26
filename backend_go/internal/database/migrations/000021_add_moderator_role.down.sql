-- Revert moderator role: promote any MODERATORs to MEMBER, then restore old constraint
UPDATE group_memberships SET role = 'MEMBER' WHERE role = 'MODERATOR';
ALTER TABLE group_memberships DROP CONSTRAINT group_memberships_role_check;
ALTER TABLE group_memberships ADD CONSTRAINT group_memberships_role_check
    CHECK (role IN ('OWNER', 'LEADER', 'MEMBER'));
