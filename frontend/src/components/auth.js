const Auth = {
    KEY: 'dnsmasq_auth',
    
    getOrCreateHeaders(headers = {}) {
        const creds = localStorage.getItem(this.KEY);
        if (creds) {
            headers['Authorization'] = `Basic ${creds}`;
        }
        return headers;
    },

    isLoggedIn() {
        return !!localStorage.getItem(this.KEY);
    },

    async login(username, password) {
        const creds = btoa(`${username}:${password}`);
        const headers = { 'Authorization': `Basic ${creds}` };
        
        try {
            // Verify credentials by calling a protected endpoint (any will do, version is lightweight)
            // But main.go excludes version? No, it excludes status and static. Version should be protected?
            // "Protect all routes except status and static files"
            // Let's check main.go from memory: status, /, /env.js, /static, /swagger are excluded.
            // So /api/v1/version IS protected. Good.
            const baseUrl = window.env && window.env.API_URL ? window.env.API_URL : '';
            const response = await fetch(`${baseUrl}/api/v1/version`, { headers });
            
            if (response.ok) {
                localStorage.setItem(this.KEY, creds);
                // Redirect to next param or root
                const params = new URLSearchParams(window.location.search);
                const next = params.get('next') || '/';
                window.location.href = next;
                return true;
            } else {
                return false;
            }
        } catch (e) {
            console.error(e);
            return false;
        }
    },

    logout() {
        localStorage.removeItem(this.KEY);
        window.location.href = '/static/pages/login.html';
    },
    
    // Call this on page load to verify auth if needed, or redirect to login
    requireAuth() {
        if (!this.isLoggedIn()) {
             const current = window.location.pathname + window.location.search;
             if (!current.includes('login.html')) {
                 window.location.href = `/static/pages/login.html?next=${encodeURIComponent(current)}`;
             }
        }
    }
};

// Global Fetch Interceptor to inject headers
(function() {
    const originalFetch = window.fetch;
    window.fetch = async function(url, options = {}) {
        options.headers = options.headers || {};
        
        // Inject auth header if available
        Auth.getOrCreateHeaders(options.headers);
        
        const response = await originalFetch(url, options);
        
        if (response.status === 401) {
            // If we got a 401, it means creds are invalid or missing
            // If we are already on login page, don't redirect loop
            if (!window.location.pathname.includes('login.html')) {
                Auth.logout(); // Will redirect to login
            }
        }
        return response;
    };
})();
