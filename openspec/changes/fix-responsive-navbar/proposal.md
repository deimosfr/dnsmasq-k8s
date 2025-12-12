# Fix Responsive Navbar

## Goal
Fix the alignment of icons and text in the navigation bar, specifically in the responsive (mobile/collapsed) view.

## Summary
The current navbar implementation places the theme toggle and logout buttons as direct children of the collapse container, leading to misalignment in mobile view. This proposal restructures the navbar to use standard Bootstrap `navbar-nav` lists for consistent alignment in both desktop and mobile views.

## Key Changes
- **Navbar Structure**: Separate main links and utility interactions (theme/logout) into distinct `navbar-nav` lists.
- **Alignment**: Use `align-items-center` to ensure specific vertical alignment where needed.
