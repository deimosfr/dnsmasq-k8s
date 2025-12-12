# Design

## Backend
Implement a manual Basic Authentication middleware.
- Read credentials from a file.
- Intercepts all requests except static assets and status.
- If `Authorization` header is missing or invalid:
  - Return `401 Unauthorized`.
  - **Crucial**: Do NOT send `WWW-Authenticate` header to prevent browser's native prompt.

## Frontend
Implement a client-side authentication flow.
- **Login Page**: `login.html` with theme consistency.
- **Auth Manager**: `auth.js` to manage credentials in `localStorage` and inject headers into `fetch` calls.
- **Redirects**: Intercept 401 responses and redirect to `login.html`.

## Helm Chart
### Values
Add `auth` section:
```yaml
auth:
  enabled: true
  existingSecret: ""
  users:
    admin: "password"
```

### Secret
Create a `Secret` resource containing the `users` list.
Mount this secret to the backend container at a known path.
