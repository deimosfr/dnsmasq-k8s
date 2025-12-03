class Navbar extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        const activePage = this.getAttribute('active-page');
        
        this.innerHTML = `
    <nav class="navbar navbar-expand-lg navbar-dark bg-success">
      <div class="container-fluid">
        <a class="navbar-brand" href="/static/index.html">Dnsmasq K8s</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav">
            <li class="nav-item">
              <a class="nav-link ${activePage === 'home' ? 'active" aria-current="page"' : '"'} href="/static/index.html">Home</a>
            </li>
            <li class="nav-item">
              <a class="nav-link ${activePage === 'config' ? 'active" aria-current="page"' : '"'} href="/static/pages/config.html">Config</a>
            </li>
            <li class="nav-item" id="nav-item-dns">
              <a class="nav-link ${activePage === 'dns' ? 'active" aria-current="page"' : '"'} href="/static/pages/dns.html">DNS</a>
            </li>
            <li class="nav-item" id="nav-item-dhcp">
              <a class="nav-link ${activePage === 'dhcp' ? 'active" aria-current="page"' : '"'} href="/static/pages/dhcp.html">DHCP</a>
            </li>
            <li class="nav-item">
              <a class="nav-link ${activePage === 'api' ? 'active" aria-current="page"' : '"'} href="/static/pages/api.html">API</a>
            </li>
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
