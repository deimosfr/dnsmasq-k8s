class Navbar extends HTMLElement {
    constructor() {
        super();
    }

    async connectedCallback() {
        const activePage = this.getAttribute('active-page');
        
        // Try to load from cache first
        const cachedItems = localStorage.getItem('navbar-cache');
        if (cachedItems) {
            try {
                const items = JSON.parse(cachedItems);
                this.render(items, activePage);
            } catch (e) {
                console.warn('Failed to parse cached navbar items', e);
            }
        }

        try {
            const response = await fetch(`${window.env.API_URL}/api/v1/navbar`);
            const items = await response.json();
            
            // Render again with fresh data
            this.render(items, activePage);
            
            // Update cache
            localStorage.setItem('navbar-cache', JSON.stringify(items));
            
            // Check redirections only after fresh data to avoid incorrect redirects from stale cache
            const itemIds = items.map(i => i.activePageId);
            if (['dns', 'dhcp'].includes(activePage) && !itemIds.includes(activePage)) {
                 window.location.href = '/static/index.html';
            }

        } catch (error) {
            console.error("Failed to load navbar:", error);
            // If we have no cache and failed to load, show error
            if (!this.innerHTML) {
                this.innerHTML = `<nav class="navbar navbar-expand-lg navbar-dark bg-danger"><div class="container-fluid"><span class="navbar-brand">Error loading menu</span></div></nav>`;
            }
        }
    }

    render(items, activePage) {
        const navItemsHtml = items.map(item => `
            <li class="nav-item" ${item.id ? `id="${item.id}"` : ''}>
              <a class="nav-link ${activePage === item.activePageId ? 'active" aria-current="page"' : '"'} href="${item.link}">${item.label}</a>
            </li>
            `).join('');

        this.innerHTML = `
        <nav class="navbar navbar-expand-lg navbar-dark bg-success">
          <div class="container-fluid">
            <a class="navbar-brand" href="/static/index.html">Dnsmasq K8s</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
              <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarNav">
              <ul class="navbar-nav">
                ${navItemsHtml}
              </ul>
              <button id="theme-toggle" class="ms-auto" title="Toggle theme">
                <i class="bi bi-moon-fill"></i>
              </button>
            </div>
          </div>
        </nav>
        `;

        // Initialize theme toggle functionality
        const themeToggle = this.querySelector('#theme-toggle');
        if (themeToggle && window.initializeThemeToggle) {
            window.initializeThemeToggle(themeToggle);
        }
    }
}

customElements.define('app-navbar', Navbar);
