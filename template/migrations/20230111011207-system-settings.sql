
-- +migrate Up
CREATE TABLE system_settings (
  id uuid not null primary key,
  registration_open boolean not null
);

INSERT INTO system_settings (id, registration_open) values ('aee228a8-657e-45af-bca6-1f6efef24ad2', true);

-- +migrate Down
DROP TABLE system_settings;
