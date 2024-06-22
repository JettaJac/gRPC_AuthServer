-- Active: 1714057451806@@127.0.0.1@5432@sso
INSERT INTO apps (id,name,secret)
VALUES(1,'test','test-secret')
ON CONFLICT DO NOTHING;