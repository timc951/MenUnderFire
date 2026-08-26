-- Repeatable migration: Test groups and memberships
-- Depends on: R__test_users.sql

DELETE FROM group_memberships WHERE group_id IN (
    '00000000-0000-0000-0000-000000000101',
    '00000000-0000-0000-0000-000000000102'
);

DELETE FROM groups WHERE id IN (
    '00000000-0000-0000-0000-000000000101',
    '00000000-0000-0000-0000-000000000102'
);

-- group_fitness: A fitness accountability group (led by user_leader)
INSERT INTO groups (id, name, description, created_by, invite_code, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000101',
    'Fitness Accountability',
    'A group for fitness goals and accountability',
    '00000000-0000-0000-0000-000000000001',
    'FITNESS2024',
    '2024-01-10 00:00:00+00',
    '2024-01-10 00:00:00+00'
);

-- group_study: A study accountability group (led by user_leader)
INSERT INTO groups (id, name, description, created_by, invite_code, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000102',
    'Study Group',
    'A group for study accountability',
    '00000000-0000-0000-0000-000000000001',
    'STUDY2024',
    '2024-01-11 00:00:00+00',
    '2024-01-11 00:00:00+00'
);

-- Memberships for group_fitness
INSERT INTO group_memberships (id, user_id, group_id, role, joined_at)
VALUES
    ('00000000-0000-0000-0000-000000000201', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000101', 'leader', '2024-01-10 00:00:00+00'),
    ('00000000-0000-0000-0000-000000000202', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000101', 'member', '2024-01-12 00:00:00+00'),
    ('00000000-0000-0000-0000-000000000203', '00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000101', 'member', '2024-01-13 00:00:00+00');

-- Memberships for group_study
INSERT INTO group_memberships (id, user_id, group_id, role, joined_at)
VALUES
    ('00000000-0000-0000-0000-000000000204', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000102', 'leader', '2024-01-11 00:00:00+00'),
    ('00000000-0000-0000-0000-000000000205', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000102', 'member', '2024-01-14 00:00:00+00');
