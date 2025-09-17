# Cry API Documentation

This document describes the available REST API endpoints for authentication, users, two-factor authentication (2FA), coin market data, and wallet exploration.
All requests and responses use **JSON** unless otherwise noted.

---

## Authentication

Some routes require a valid **JWT token** in the `Authorization` header:

```
Authorization: Bearer <token>
```

---

## Users

### `POST /users/signup`

Create a new user account.

### `POST /users/verify-email-token`

Verify a user using an email OTP token.

### `POST /users/verify-account-token`

Verify a user using an account verification token (URL-based).

### `POST /users/signin`

Authenticate a user and return a JWT.

### `POST /users/reset-password`

Request a password reset. An OTP or token will be sent to the user’s email.

### `POST /users/verify-reset-password-token`

Verify a reset password token and set a new password.

---

## Two-Factor Authentication (2FA)

### `POST /2fa/setup`

Generate a new 2FA secret and QR code for the authenticated user.

### `POST /2fa/setup/verify-otp`

Verify the OTP entered during initial 2FA setup.

### `POST /2fa/auth/verify-otp`

Verify the OTP during login/authentication.

### `POST /2fa/alternative/send-email-otp`

Send a one-time login OTP via email as an alternative to app-based 2FA.

---

## Coin MarketCap

> **Authentication Required** (JWT)

### `GET /coin-marketcap/fear-and-greed-lastest`

Get the latest **Fear & Greed Index**.

### `GET /coin-marketcap/fear-and-greed-historical`

Get historical data for the **Fear & Greed Index**.

---

## Wallet Explorer

> **Authentication Required** (JWT)

### `GET /wallet-explorer/tx`

Retrieve transaction information for a given transaction hash.

### `GET /wallet-explorer/xpub`

Retrieve transactions associated with an **XPUB** key.

---

## Notes

* All timestamps are returned in **UTC**.
* Failed requests include a `message` field describing the error.
* Future versions may add rate limiting and pagination for list-based responses.
