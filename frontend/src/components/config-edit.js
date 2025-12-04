const RECOMMENDED_CONFIG = `# dnsmasq configuration file

## Global config
# Interface to serve DHCP on the host
interface=eth0
# Listen on specified address
listen-address=0.0.0.0
# Upstream DNS servers
server=8.8.8.8
# Don't run as daemon (required for container)
keep-in-foreground
# Log to stdout (for container logging)
log-facility=-
# Custom configuration
conf-dir=/etc/dnsmasq.d

## DNS
# Set domain for local network
domain=mydomain.local
# Cache size
cache-size=1000
# Don't read /etc/resolv.conf
no-resolv
# Don't read /etc/hosts
no-hosts
# Disable negative cache
no-negcache

## DHCP Configuration
dhcp-range=192.168.0.100,192.168.94.200,24h
# Set gateway
dhcp-option=3,192.168.0.1
# Set DNS servers for DHCP clients
dhcp-option=6,192.168.0.1,8.8.8.8
# Enable DHCP authoritative mode
dhcp-authoritative
# Push search domain to DHCP
dhcp-option=option:domain-search,mydomain.local`;

// Display config in view-only mode
async function displayConfigView() {
    const config = await getConfig();
    const configDisplay = document.getElementById('config-display');
    configDisplay.value = config;
}

// Render edit mode
function renderConfigEdit(config) {
    const configDiv = document.getElementById('config-edit');
    configDiv.innerHTML = `
        <textarea id="config-textarea" class="form-control font-monospace" style="height: 70vh;">${config}</textarea>
        <div class="mt-3">
            <button id="save-button" class="btn btn-success">Save</button>
            <button id="cancel-button" class="btn btn-secondary ms-2">Cancel</button>
            <button id="recommended-button" class="btn btn-purple ms-2">Load Recommended Config</button>
        </div>
    `;

    // Save button handler
    document.getElementById('save-button').addEventListener('click', async () => {
        const newConfig = document.getElementById('config-textarea').value;
        try {
            await updateConfig(newConfig);
            if (window.showRestartBanner) {
                window.showRestartBanner();
            }
            alert('Configuration saved successfully!');
            switchToViewMode();
        } catch (error) {
            alert('Failed to save configuration: ' + error.message);
        }
    });

    // Cancel button handler
    document.getElementById('cancel-button').addEventListener('click', () => {
        if (confirm('Discard changes and return to view mode?')) {
            switchToViewMode();
        }
    });

    // Recommended config button handler
    document.getElementById('recommended-button').addEventListener('click', () => {
        if (confirm('This will replace the current configuration in the editor. Are you sure?')) {
            document.getElementById('config-textarea').value = RECOMMENDED_CONFIG;
        }
    });
}

// Switch to edit mode
async function switchToEditMode() {
    const config = await getConfig();
    renderConfigEdit(config);
    document.getElementById('config-view').style.display = 'none';
    document.getElementById('config-edit').style.display = 'block';
}

// Switch to view mode
async function switchToViewMode() {
    await displayConfigView();
    document.getElementById('config-edit').style.display = 'none';
    document.getElementById('config-view').style.display = 'block';
}

// Get config via API
async function getConfig() {
    const response = await fetch(`${window.env.API_URL}/api/v1/config`);
    if (!response.ok) {
        throw new Error('Failed to fetch configuration');
    }
    return await response.text();
}

// Update config via API
async function updateConfig(config) {
    const response = await fetch(`${window.env.API_URL}/api/v1/config`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ config }),
    });
    
    if (!response.ok) {
        throw new Error('Failed to update configuration');
    }
}

// Initialize: Show view mode by default
displayConfigView();

// Edit button handler
document.getElementById('edit-button').addEventListener('click', switchToEditMode);
