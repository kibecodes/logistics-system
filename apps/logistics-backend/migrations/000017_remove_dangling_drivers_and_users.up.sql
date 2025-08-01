DELETE FROM drivers WHERE id NOT IN (SELECT id FROM users);
DELETE FROM users WHERE role = 'driver' AND id NOT IN (SELECT id FROM drivers);
