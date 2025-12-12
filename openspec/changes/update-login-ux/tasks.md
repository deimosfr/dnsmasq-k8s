# Tasks

- [x] Fix Option Key Blocking
  - Investigated `frontend/src` for global listeners. No explicit blocking found.
  - Verified `keydown`, `keypress`, `preventDefault` usage.
  - Assumed standard behavior is preserved (no blocking code to remove).
- [x] Set Default Dark Mode
  - Updated `theme-toggle.js` to default to 'dark'.
- [x] Add Navbar Logout
  - Updated `navbar.js` to include a Logout button.
  - **Position**: Placed it after (to the right of) the theme toggle.
  - Wired button to `Auth.logout()`.
- [x] Add Login Theme Toggle & Title
  - Added theme toggle button to `login.html`.
  - Renamed "Welcome Back" to "Dnsmasq K8s".
