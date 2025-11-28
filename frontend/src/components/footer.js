document.addEventListener('DOMContentLoaded', async () => {
    const footer = document.createElement('footer');
    footer.className = 'footer mt-auto py-3 text-center';
    
    try {
        const response = await fetch(`${window.env.API_URL}/api/v1/version`);
        const data = await response.json();
        const version = data.version || 'unknown';
        
        footer.innerHTML = `
            <div class="container">
                <span class="text-muted">
                    Dnsmasq K8s v${version} | 
                    <a href="https://github.com/deimosfr/dnsmasq-k8s-ui" target="_blank" class="text-decoration-none text-secondary">
                        <i class="bi bi-github"></i> GitHub
                    </a>
                </span>
            </div>
        `;
    } catch (error) {
        console.error('Failed to fetch version:', error);
        footer.innerHTML = `
            <div class="container">
                <span class="text-muted">
                    Dnsmasq K8s | 
                    <a href="https://github.com/deimosfr/dnsmasq-k8s-ui" target="_blank" class="text-decoration-none text-secondary">
                        <i class="bi bi-github"></i> GitHub
                    </a>
                </span>
            </div>
        `;
    }
    
    document.body.appendChild(footer);
});
