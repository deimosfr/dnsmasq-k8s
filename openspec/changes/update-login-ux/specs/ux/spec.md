# UX Improvements

## ADDED Requirements

### Requirement: Default Dark Mode
The system MUST default to Dark Mode if no user preference is saved.

#### Scenario: First Visit
Given the user visits the site for the first time
Then the UI MUST be rendered in Dark Mode.

### Requirement: Logout Navigation
The navigation bar MUST provide a user-accessible way to log out.
The Logout button MUST be positioned to the right of the theme toggle button.

#### Scenario: Logout Action
Given the user is logged in
When they click the Logout button in the navbar
Then their session credentials MUST be cleared
And they MUST be redirected to the Login page.

### Requirement: Login Theme Toggle
The Login page MUST allow the user to toggle between Light and Dark themes.

#### Scenario: Toggling Theme on Login
Given the user is on the Login page
When they click the theme toggle button
Then the UI theme MUST switch.

### Requirement: Input Accessibility
Text inputs MUST support standard operating system modifier keys.

#### Scenario: Option Key on macOS
Given the user is typing in a text input
When they press the Option key combinations (e.g., Option+e for accents)
Then the corresponding character MUST be input
And the key press MUST NOT be blocked.
