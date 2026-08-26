-- Repeatable migration: Test reports
-- Depends on: R__test_users.sql, R__test_groups.sql

DELETE FROM reports WHERE id IN (
    '00000000-0000-0000-0000-000000000301',
    '00000000-0000-0000-0000-000000000302',
    '00000000-0000-0000-0000-000000000303'
);

-- report_public: A non-anonymous report by user_member1
INSERT INTO reports (id, user_id, group_id, title, content, is_anonymous_to_group, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000301',
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000101',
    'Weekly Fitness Update',
    'Completed 5 workouts this week. Feeling great about progress on deadlifts.',
    FALSE,
    '2024-01-20 10:00:00+00'
);

-- report_anonymous: An anonymous report by user_member2
INSERT INTO reports (id, user_id, group_id, title, content, is_anonymous_to_group, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000302',
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0000-000000000101',
    'Struggling This Week',
    'Had a tough week. Only managed 2 sessions. Need more accountability.',
    TRUE,
    '2024-01-21 14:00:00+00'
);

-- report_leader: A report by the leader
INSERT INTO reports (id, user_id, group_id, title, content, is_anonymous_to_group, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000303',
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000101',
    'Leader Check-in',
    'Leading by example - hit all my targets this week.',
    FALSE,
    '2024-01-22 09:00:00+00'
);
