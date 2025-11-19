// Theme toggle functionality
(function() {
    const THEME_KEY = 'dnsmasq-theme';
    const THEME_DARK = 'dark';
    const THEME_LIGHT = 'light';
    
    // Get saved theme or default to dark
    function getSavedTheme() {
        const saved = localStorage.getItem(THEME_KEY);
        return saved || THEME_DARK;
    }
    
    // Apply theme to document
    function applyTheme(theme) {
        if (theme === THEME_DARK) {
            document.documentElement.setAttribute('data-theme', 'dark');
        } else {
            document.documentElement.removeAttribute('data-theme');
        }
        updateToggleIcon(theme);
    }
    
    // Update toggle button icon
    function updateToggleIcon(theme) {
        const toggleBtn = document.getElementById('theme-toggle');
        if (toggleBtn) {
            if (theme === THEME_DARK) {
                // Show sun icon in dark mode (click to go light)
                toggleBtn.innerHTML = '<i class="bi bi-sun-fill"></i>';
            } else {
                // Show moon icon in light mode (click to go dark)
                toggleBtn.innerHTML = '<i class="bi bi-moon-fill"></i>';
            }
        }
    }
    
    // Toggle theme
    function toggleTheme() {
        const currentTheme = getSavedTheme();
        const newTheme = currentTheme === THEME_DARK ? THEME_LIGHT : THEME_DARK;
        localStorage.setItem(THEME_KEY, newTheme);
        applyTheme(newTheme);
    }
    
    // Apply theme immediately (before DOM loads to prevent flash)
    applyTheme(getSavedTheme());
    
    // Set up toggle button when DOM is ready
    document.addEventListener('DOMContentLoaded', function() {
        const toggleBtn = document.getElementById('theme-toggle');
        if (toggleBtn) {
            toggleBtn.addEventListener('click', toggleTheme);
            updateToggleIcon(getSavedTheme());
        }
    });
})();
