-- Repeatable migration: Test users
-- These are idempotent - safe to re-run

DELETE FROM users WHERE id IN (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000004'
);

-- user_leader: A user who creates and leads groups
INSERT INTO users (id, auth0_id, email, display_name, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'auth0|leader001',
    'leader@test.com',
    'Test Leader',
    '2024-01-01 00:00:00+00',
    '2024-01-01 00:00:00+00'
);

-- user_member1: A regular group member
INSERT INTO users (id, auth0_id, email, display_name, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    'auth0|member001',
    'member1@test.com',
    'Test Member 1',
    '2024-01-02 00:00:00+00',
    '2024-01-02 00:00:00+00'
);

-- user_member2: Another regular group member
INSERT INTO users (id, auth0_id, email, display_name, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000003',
    'auth0|member002',
    'member2@test.com',
    'Test Member 2',
    '2024-01-03 00:00:00+00',
    '2024-01-03 00:00:00+00'
);

-- user_outsider: A user not in any test groups
INSERT INTO users (id, auth0_id, email, display_name, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000004',
    'auth0|outsider001',
    'outsider@test.com',
    'Test Outsider',
    '2024-01-04 00:00:00+00',
    '2024-01-04 00:00:00+00'
);
