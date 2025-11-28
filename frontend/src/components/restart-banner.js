// Restart banner state management
console.log('Restart banner script loaded');
const RESTART_BANNER_KEY = 'dnsmasq-restart-needed';

// Show the restart banner
window.showRestartBanner = function() {
    localStorage.setItem(RESTART_BANNER_KEY, 'true');
    displayBanner();
};

// Hide the restart banner
window.hideRestartBanner = function() {
    localStorage.removeItem(RESTART_BANNER_KEY);
    const banner = document.getElementById('restart-banner');
    if (banner) {
        banner.style.display = 'none';
    }
};

// Display the banner if needed
function displayBanner() {
    const banner = document.getElementById('restart-banner');
    if (banner && localStorage.getItem(RESTART_BANNER_KEY) === 'true') {
        banner.style.display = 'flex';
    }
}

// Restart dnsmasq and hide the banner
async function restartDnsmasq() {
    try {
        const response = await fetch(window.env.API_URL + '/api/v1/supervisor/dnsmasq/restart', {
            method: 'POST',
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            alert(`Error restarting dnsmasq: ${errorData.error || 'Unknown error'}`);
            return;
        }
        
        // Hide the banner on successful restart
        hideRestartBanner();
        
        // Optionally show a success message
        console.log('Dnsmasq restarted successfully');
    } catch (error) {
        console.error('Failed to restart dnsmasq:', error);
        alert(`Failed to restart dnsmasq: ${error.message}`);
    }
}

// Create and inject the banner HTML
function createBanner() {
    // Check if banner already exists
    if (document.getElementById('restart-banner')) {
        return;
    }
    
    const banner = document.createElement('div');
    banner.id = 'restart-banner';
    banner.className = 'restart-banner';
    banner.innerHTML = `
        <div class="restart-banner-content">
            <div class="restart-banner-message">
                <i class="bi bi-exclamation-triangle-fill me-2"></i>
                <span>Configuration changes require a dnsmasq restart to take effect.</span>
            </div>
            <div class="restart-banner-actions">
                <button class="btn btn-sm btn-dark me-2" onclick="restartDnsmasq()">
                    <i class="bi bi-arrow-clockwise me-1"></i>Restart Now
                </button>
                <button class="btn-close btn-close-white" onclick="hideRestartBanner()" aria-label="Dismiss"></button>
            </div>
        </div>
    `;
    
    // Insert banner after navbar
    const navbar = document.querySelector('nav.navbar');
    if (navbar && navbar.nextSibling) {
        navbar.parentNode.insertBefore(banner, navbar.nextSibling);
    } else if (navbar) {
        navbar.parentNode.appendChild(banner);
    } else {
        // Fallback: insert at the beginning of body
        document.body.insertBefore(banner, document.body.firstChild);
    }
    
    // Check if banner should be visible
    displayBanner();
}

// Listen for storage events to sync across tabs
window.addEventListener('storage', (event) => {
    if (event.key === RESTART_BANNER_KEY) {
        if (event.newValue === 'true') {
            displayBanner();
        } else {
            const banner = document.getElementById('restart-banner');
            if (banner) {
                banner.style.display = 'none';
            }
        }
    }
});

// Make restartDnsmasq globally accessible
window.restartDnsmasq = restartDnsmasq;

// Initialize banner on page load
document.addEventListener('DOMContentLoaded', () => {
    createBanner();
});
