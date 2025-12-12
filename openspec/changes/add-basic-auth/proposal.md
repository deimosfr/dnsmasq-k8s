# Add Basic Authentication

## Goal
Secure the Dnsmasq UI with Basic Authentication.

## Summary
Add a mechanism to restrict access to the UI using a login/password list.
The list is stored in a file managed by a Kubernetes Secret.
Users can configure default credentials or disable authentication via Helm values.
Accessing the UI without credentials redirects to a custom Login Page.

## Key Changes
- **Backend**: Middleware for manual Basic Auth check (suppressing browser prompt).
- **Frontend**: Custom Login Page (`login.html`) and client-side auth management (`auth.js`).
- **Helm**: Secret management and volume mounting.
- **Docs**: Updated README.
