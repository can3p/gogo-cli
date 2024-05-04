
-- +migrate Up
CREATE TABLE user_signup_requests (
  id uuid primary key,
  email varchar not null,
  reason varchar,
  signup_attribution varchar,
  created_user_id uuid references users(id),
  verification_sent_at timestamp,
  email_confirmed_at timestamp,
  created_at timestamp not null,
  updated_at timestamp not null,
);

create unique index on user_signup_requests(email);

-- +migrate Down
