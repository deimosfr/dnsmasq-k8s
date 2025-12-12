# Design

## Navbar Restructuring
Current structure:
```html
<div class="collapse navbar-collapse">
  <ul class="navbar-nav">...</ul>
  <button id="theme-toggle" class="ms-auto ...">...</button>
  <a class="nav-link ..." ...>...</a>
</div>
```

Proposed structure:
```html
<div class="collapse navbar-collapse">
  <ul class="navbar-nav me-auto">
    <!-- Main Links -->
    ...
  </ul>
  <ul class="navbar-nav ms-auto align-items-lg-center">
    <!-- Right-aligned items -->
    <li class="nav-item">
      <button id="theme-toggle" class="nav-link btn btn-link ...">...</button>
    </li>
    <li class="nav-item">
      <a class="nav-link" ...>...</a>
    </li>
  </ul>
</div>
```
This ensures:
1.  **Desktop**: `ms-auto` pushes the second list to the right. `align-items-lg-center` ensures vertical centering.
2.  **Mobile**: Both lists stack vertically as standard nav items, respecting Bootstrap's mobile navbar styles (padding, alignment).
