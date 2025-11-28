// Sort state
let leasesSortState = { column: 'expiry_time', direction: 'asc' };

// Generic sort function
function sortData(data, column, direction) {
    return [...data].sort((a, b) => {
        let aVal = a[column] || '';
        let bVal = b[column] || '';
        
        if (typeof aVal === 'string') {
            aVal = aVal.toLowerCase();
            bVal = bVal.toLowerCase();
        }
        
        if (aVal < bVal) return direction === 'asc' ? -1 : 1;
        if (aVal > bVal) return direction === 'asc' ? 1 : -1;
        return 0;
    });
}

// Update sort icons
function updateLeasesSortIcons(column, direction) {
    ['mac_address', 'ip_address', 'hostname', 'expiry_time'].forEach(col => {
        const icon = document.getElementById(`sort-icon-lease-${col}`);
        if (icon) {
            if (col === column) {
                icon.className = direction === 'asc' ? 'bi bi-chevron-up' : 'bi bi-chevron-down';
            } else {
                icon.className = 'bi bi-chevron-expand';
            }
        }
    });
}

// Sort leases
window.sortLeases = function(column) {
    if (leasesSortState.column === column) {
        leasesSortState.direction = leasesSortState.direction === 'asc' ? 'desc' : 'asc';
    } else {
        leasesSortState.column = column;
        leasesSortState.direction = 'asc';
    }
    displayLeases();
};

async function getLeases() {
    const response = await fetch(`${window.env.API_URL}/api/v1/dhcp/leases`);
    const data = await response.json();
    return data.leases;
}

async function displayLeases() {
    const leases = await getLeases();
    const tbody = document.getElementById('leases-table-body');
    tbody.innerHTML = '';

    const header = document.getElementById('dhcp-leases-header');
    if (header) {
        header.innerHTML = `Current DHCP Leases <span class="badge rounded-pill bg-success ms-2">${leases ? leases.length : 0}</span>`;
    }

    if (!leases || leases.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="text-center">No leases found</td></tr>';
        return;
    }

    // Sort leases
    const sortedLeases = sortData(leases, leasesSortState.column, leasesSortState.direction);

    sortedLeases.forEach((lease, index) => {
        const row = document.createElement('tr');
        
        // Calculate timing
        const expiryDate = new Date(lease.expiry_time * 1000);
        const now = new Date();
        const remainingMs = expiryDate - now;
        
        let remainingStr = "Expired";
        if (remainingMs > 0) {
            const hours = Math.floor(remainingMs / (1000 * 60 * 60));
            const minutes = Math.floor((remainingMs % (1000 * 60 * 60)) / (1000 * 60));
            remainingStr = `${hours}h ${minutes}m`;
        }

        const expiryStr = expiryDate.toLocaleString('en-GB', { 
            year: 'numeric', 
            month: '2-digit', 
            day: '2-digit', 
            hour: '2-digit', 
            minute: '2-digit', 
            second: '2-digit',
            hour12: false 
        });

        row.innerHTML = `
            <td>${lease.mac_address}</td>
            <td>${lease.ip_address}</td>
            <td>${lease.hostname}</td>
            <td><small>${expiryStr}</small></td>
            <td><span class="badge ${remainingMs > 0 ? 'bg-success' : 'bg-secondary'}">${remainingStr}</span></td>
            <td>
                <i class="bi bi-pencil text-success me-3" style="cursor: pointer;" onclick="editLease(${index}, '${lease.mac_address}', '${lease.ip_address}', '${lease.hostname}')"></i>
                <i class="bi bi-x-lg text-danger" style="cursor: pointer;" onclick="deleteLease('${lease.mac_address}', '${lease.ip_address}', '${lease.hostname}')"></i>
            </td>
        `;
        tbody.appendChild(row);
    });

    // Update sort icons
    updateLeasesSortIcons(leasesSortState.column, leasesSortState.direction);
}

// Track if a row is currently being edited
let currentlyEditingLease = null;

window.editLease = function(index, mac, ip, hostname) {
    // If another row is being edited, cancel that edit first
    if (currentlyEditingLease !== null && currentlyEditingLease !== index) {
        displayLeases();
    }
    
    currentlyEditingLease = index;
    const tbody = document.getElementById('leases-table-body');
    const row = tbody.children[index];
    
    row.innerHTML = `
        <td><input type="text" class="form-control form-control-sm" id="edit-lease-mac-${index}" value="${mac}"></td>
        <td><input type="text" class="form-control form-control-sm" id="edit-lease-ip-${index}" value="${ip}"></td>
        <td><input type="text" class="form-control form-control-sm" id="edit-lease-host-${index}" value="${hostname}"></td>
        <td colspan="2"></td>
        <td>
            <i class="bi bi-check-lg text-success me-3" style="cursor: pointer;" onclick="saveLease(${index}, '${mac}', '${ip}', '${hostname}')"></i>
            <i class="bi bi-x-circle text-secondary" style="cursor: pointer;" onclick="cancelLeaseEdit()"></i>
        </td>
    `;
}

window.cancelLeaseEdit = function() {
    currentlyEditingLease = null;
    displayLeases();
}

window.saveLease = async function(index, oldMac, oldIp, oldHostname) {
    const newMac = document.getElementById(`edit-lease-mac-${index}`).value;
    const newIp = document.getElementById(`edit-lease-ip-${index}`).value;
    const newHostname = document.getElementById(`edit-lease-host-${index}`).value;

    // Validation
    const macRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
    const ipRegex = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

    if (!macRegex.test(newMac)) {
        alert('Invalid MAC address format. Use XX:XX:XX:XX:XX:XX');
        return;
    }
    if (!ipRegex.test(newIp)) {
        alert('Invalid IP address format.');
        return;
    }
    if (!newHostname || newHostname.trim() === '') {
        alert('Hostname cannot be empty.');
        return;
    }

    await fetch(window.env.API_URL + '/api/v1/dhcp/leases', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            old: { mac_address: oldMac, ip_address: oldIp, hostname: oldHostname },
            new: { mac_address: newMac, ip_address: newIp, hostname: newHostname }
        }),
    });

    currentlyEditingLease = null;
    displayLeases();
    showRestartBanner();
}

window.deleteLease = async function(mac, ip, hostname) {
    if (!confirm(`Are you sure you want to delete lease for ${mac} (${ip})?`)) return;

    await fetch(`${window.env.API_URL}/api/v1/dhcp/leases`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mac_address: mac, ip_address: ip, hostname: hostname }),
    });

    displayLeases();
    showRestartBanner();
}

displayLeases();
