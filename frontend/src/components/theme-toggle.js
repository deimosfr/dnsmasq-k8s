// Theme toggle functionality
window.initializeThemeToggle = function(themeToggleBtn) {
    if (themeToggleBtn.dataset.initialized === 'true') return;
    themeToggleBtn.dataset.initialized = 'true';

    const icon = themeToggleBtn.querySelector('i');

    function updateIcon(theme) {
        if (theme === 'dark') {
            icon.classList.remove('bi-moon-fill');
            icon.classList.add('bi-sun-fill');
        } else {
            icon.classList.remove('bi-sun-fill');
            icon.classList.add('bi-moon-fill');
        }
    }

    // Initialize icon based on current theme
    const currentTheme = document.documentElement.getAttribute('data-bs-theme');
    updateIcon(currentTheme);

    themeToggleBtn.addEventListener('click', () => {
        const currentTheme = document.documentElement.getAttribute('data-bs-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        document.documentElement.setAttribute('data-bs-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        updateIcon(newTheme);
    });
};

// Apply theme immediately
(function() {
    const savedTheme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-bs-theme', savedTheme);
})();

document.addEventListener('DOMContentLoaded', () => {
    // If there is a theme toggle on the page (legacy), initialize it
    // The navbar component will call window.initializeThemeToggle itself
    const themeToggleBtn = document.getElementById('theme-toggle');
    if (themeToggleBtn) {
        window.initializeThemeToggle(themeToggleBtn);
    }
});
