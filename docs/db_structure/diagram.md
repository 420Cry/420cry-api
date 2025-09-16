https://dbdiagram.io/d
```bash
// Users
Table users {
  id integer [primary key]
  uuid varchar [unique, not null]
  username varchar [unique, not null]
  email varchar [unique, not null]
  fullname varchar
  password varchar [not null]
  is_verified boolean [default: false, not null]
  two_fa_secret varchar
  two_fa_enabled boolean [default: false, not null]
  created_at timestamp [default: `CURRENT_TIMESTAMP`, not null]
  updated_at timestamp
}

// Generic tokens (OTP, verification, password reset, etc.)
Table user_tokens {
  id integer [primary key]
  user_id integer [not null]
  token varchar [not null] // could be random string or 6-digit code
  purpose varchar [not null] // e.g. 'account_verification', 'password_reset', 'login_otp', 'transaction'
  created_at timestamp [default: `CURRENT_TIMESTAMP`, not null]
  expires_at timestamp [not null]
  consumed boolean [default: false, not null]
  used_at timestamp
}

// Relationships
Ref: user_tokens.user_id > users.id
```