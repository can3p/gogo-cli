
-- +migrate Up
create table users (
  id uuid not null primary key,
  email varchar not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  timezone varchar not null,
  email_confirmed_at timestamp without time zone,
  email_confirm_seed uuid,
  signup_attribution varchar,
  pwdhash varchar
);

create unique index on users(email);
create index on users(pwdhash);

-- +migrate Down
