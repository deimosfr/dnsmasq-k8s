# Design

## Theme Logic
- **Defaulting**: Change `theme-toggle.js` logic: `const savedTheme = localStorage.getItem('theme') || 'dark';` (previously 'light').
- **Login Page**: Inject the toggle button. Rename "Welcome Back" title to "Dnsmasq K8s".

## Navbar Logout
- Add a new button to the navbar.
- **Placement**: MUST be on the right of the theme toggle button.
- On click, call `Auth.logout()`.

## Option Key Issue
- **Investigation**: Check for `keydown/keyup` listeners preventing default.
- **Fix**: Remove or refine any specific key blocks. (Current analysis shows no obvious blocking code, so this might be a browser interaction or `bootstrap` oddity to be verified).
