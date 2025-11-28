async function getStatus() {
    const response = await fetch(`${window.env.API_URL}/api/v1/status`);
    const data = await response.json();
    return data;
}

function renderStatus(status) {
    const apiStatusDiv = document.getElementById('api-status');
    const dnsStatusDiv = document.getElementById('dns-status');
    const dhcpStatusDiv = document.getElementById('dhcp-status');

    const pidDiv = document.getElementById('dnsmasq-pid');

    if (apiStatusDiv) {
        if (status.api) {
            apiStatusDiv.innerHTML = '<span class="badge bg-success">Enabled</span>';
        } else {
            apiStatusDiv.innerHTML = '<span class="badge bg-secondary">Disabled</span>';
        }
    }

    if (dnsStatusDiv) {
        if (status.dns) {
            dnsStatusDiv.innerHTML = '<span class="badge bg-success">Enabled</span>';
        } else {
            dnsStatusDiv.innerHTML = '<span class="badge bg-secondary">Disabled</span>';
        }
    }

    if (dhcpStatusDiv) {
        if (status.dhcp) {
            dhcpStatusDiv.innerHTML = '<span class="badge bg-success">Enabled</span>';
        } else {
            dhcpStatusDiv.innerHTML = '<span class="badge bg-secondary">Disabled</span>';
        }
    }
    

    
    if (pidDiv && status.dnsmasq_pid) {
        pidDiv.textContent = status.dnsmasq_pid;
    }

    const supervisorDiv = document.getElementById('supervisor-status');
    if (supervisorDiv && status.supervisor_services) {
        let html = '<ul class="list-group list-group-flush">';
        status.supervisor_services.forEach(service => {
            const parts = service.split(/\s+/);
            const name = parts[0];
            const state = parts[1];
            let uptime = '';
            
            // Parse uptime if available (usually after "uptime")
            const uptimeIndex = parts.indexOf('uptime');
            if (uptimeIndex !== -1 && uptimeIndex + 1 < parts.length) {
                uptime = parts[uptimeIndex + 1].replace(',', '');
            }

            let badgeClass = 'bg-secondary';
            if (state === 'RUNNING') {
                badgeClass = 'bg-success';
            } else if (state === 'STOPPED' || state === 'EXITED') {
                badgeClass = 'bg-warning text-dark';
            } else if (state === 'FATAL' || state === 'BACKOFF') {
                badgeClass = 'bg-danger';
            }
            
            html += `<li class="list-group-item d-flex justify-content-between align-items-center">
                <div>
                    ${name}
                    <span class="badge ${badgeClass} ms-2">${state}</span>
                    ${uptime ? `<span class="badge bg-purple ms-1" data-bs-toggle="tooltip" title="Uptime">${uptime}</span>` : ''}
                </div>
                <div class="btn-group btn-group-sm gap-1" role="group">
                    <button class="btn btn-sm btn-outline-success" onclick="controlSupervisor('${name}', 'start', this)" data-bs-toggle="tooltip" title="Start service">
                        <i class="bi bi-play-fill"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-danger" onclick="controlSupervisor('${name}', 'stop', this)" data-bs-toggle="tooltip" title="Stop service">
                        <i class="bi bi-stop-fill"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-primary" onclick="controlSupervisor('${name}', 'restart', this)" data-bs-toggle="tooltip" title="Restart service">
                        <i class="bi bi-arrow-clockwise"></i>
                    </button>
                </div>
            </li>`;
        });
        html += '</ul>';
        supervisorDiv.innerHTML = html;
        
        // Initialize Bootstrap tooltips
        const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]');
        [...tooltipTriggerList].map(tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl, { trigger: 'hover' }));
    }

    updateNavbar(status);
}

function updateNavbar(status) {
    const dnsNavItem = document.getElementById('nav-item-dns');
    const dhcpNavItem = document.getElementById('nav-item-dhcp');

    if (dnsNavItem) {
        if (status.dns) {
            dnsNavItem.style.display = 'block';
        } else {
            dnsNavItem.style.display = 'none';
            // Redirect if on DNS page and it's disabled
            if (window.location.pathname.includes('/dns.html')) {
                window.location.href = '/static/index.html';
            }
        }
    }

    if (dhcpNavItem) {
        if (status.dhcp) {
            dhcpNavItem.style.display = 'block';
        } else {
            dhcpNavItem.style.display = 'none';
            // Redirect if on DHCP page and it's disabled
            if (window.location.pathname.includes('/dhcp.html')) {
                window.location.href = '/static/index.html';
            }
        }
    }
}

async function displayStatus() {
    try {
        const status = await getStatus();
        renderStatus(status);
    } catch (error) {
        console.error("Failed to fetch status:", error);
    }
}

// Control supervisor service
window.controlSupervisor = async function(serviceName, action, btn) {
    // Hide tooltip immediately
    if (btn) {
        const tooltip = bootstrap.Tooltip.getInstance(btn);
        if (tooltip) {
            tooltip.hide();
        }
        btn.blur(); // Remove focus
    }

    try {
        const response = await fetch(`${window.env.API_URL}/api/v1/supervisor/${serviceName}/${action}`, {
            method: 'POST',
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            alert(`Error: ${errorData.error || 'Failed to ' + action + ' service'}`);
            return;
        }
        
        // Hide restart banner if dnsmasq was restarted
        if (serviceName === 'dnsmasq' && action === 'restart') {
            if (typeof hideRestartBanner === 'function') {
                hideRestartBanner();
            }
        }
        
        // Wait a moment for the service to change state
        setTimeout(() => {
            displayStatus();
        }, 500);
    } catch (error) {
        console.error(`Failed to ${action} service ${serviceName}:`, error);
        alert(`Failed to ${action} service: ${error.message}`);
    }
};

document.addEventListener('DOMContentLoaded', () => {
    displayStatus();
    // Refresh status every 3 seconds
    setInterval(displayStatus, 3000);
});
