# UI Alignment

## MODIFIED Requirements

### Requirement: Responsive Alignment
Navbar items (links, buttons, icons) /MUST/ be aligned consistently in both desktop and mobile views.

#### Scenario: Mobile View Alignment
Given the browser window is resized to mobile width (collapsed navbar)
When the navbar is expanded
Then all items (links, theme toggle, logout) MUST be stacked vertically
And aligned to the start (left) with consistent padding
And icons MUST be aligned with text baselines.
