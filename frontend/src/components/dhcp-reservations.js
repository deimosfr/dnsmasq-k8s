// Sort state
let reservationsSortState = { column: 'hostname', direction: 'asc' };

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
function updateReservationsSortIcons(column, direction) {
    ['mac_address', 'ip_address', 'hostname'].forEach(col => {
        const icon = document.getElementById(`sort-icon-res-${col}`);
        if (icon) {
            if (col === column) {
                icon.className = direction === 'asc' ? 'bi bi-chevron-up' : 'bi bi-chevron-down';
            } else {
                icon.className = 'bi bi-chevron-expand';
            }
        }
    });
}

// Sort reservations
window.sortReservations = function(column) {
    if (reservationsSortState.column === column) {
        reservationsSortState.direction = reservationsSortState.direction === 'asc' ? 'desc' : 'asc';
    } else {
        reservationsSortState.column = column;
        reservationsSortState.direction = 'asc';
    }
    displayReservations();
};

document.getElementById('add-reservation-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const macAddress = document.getElementById('mac-address').value;
    const ipAddress = document.getElementById('ip-address').value;
    const hostname = document.getElementById('hostname').value;
    const comment = document.getElementById('comment').value;

    await addReservation(macAddress, ipAddress, hostname, comment);
    document.getElementById('add-reservation-form').reset();
    displayReservations();
    showRestartBanner();
});

async function addReservation(macAddress, ipAddress, hostname, comment) {
    const response = await fetch(`${window.env.API_URL}/api/v1/dhcp/reservations`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            mac_address: macAddress,
            ip_address: ipAddress,
            hostname: hostname,
            comment: comment,
        }),
    });
    if (!response.ok) {
        const data = await response.json();
        alert(`Error: ${data.error}`);
    }
}

async function getReservations() {
    const response = await fetch(`${window.env.API_URL}/api/v1/dhcp/reservations`);
    const data = await response.json();
    return data.reservations;
}

async function displayReservations() {
    const reservations = await getReservations();
    const tbody = document.getElementById('reservations-table-body');
    tbody.innerHTML = '';

    const header = document.getElementById('dhcp-reservations-header');
    if (header) {
        header.innerHTML = `Current Reservations <span class="badge rounded-pill bg-success ms-2">${reservations ? reservations.length : 0}</span>`;
    }

    if (!reservations || reservations.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" class="text-center">No reservations found</td></tr>';
        return;
    }

    // Sort reservations
    const sortedReservations = sortData(reservations, reservationsSortState.column, reservationsSortState.direction);

    sortedReservations.forEach((res, index) => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${res.mac_address}</td>
            <td>${res.ip_address}</td>
            <td>${res.hostname}</td>
            <td>${res.comment || ''}</td>
            <td>
                <i class="bi bi-pencil text-success me-3" style="cursor: pointer;" onclick="editReservation(${index}, '${res.mac_address}', '${res.ip_address}', '${res.hostname}', '${res.comment || ''}')"></i>
                <i class="bi bi-x-lg text-danger" style="cursor: pointer;" onclick="deleteReservation('${res.mac_address}', '${res.ip_address}', '${res.hostname}')"></i>
            </td>
        `;
        tbody.appendChild(row);
    });

    // Update sort icons
    updateReservationsSortIcons(reservationsSortState.column, reservationsSortState.direction);
}

window.editReservation = function(index, mac, ip, hostname, comment) {
    const tbody = document.getElementById('reservations-table-body');
    const row = tbody.children[index];
    
    row.innerHTML = `
        <td><input type="text" class="form-control form-control-sm" id="edit-res-mac-${index}" value="${mac}"></td>
        <td><input type="text" class="form-control form-control-sm" id="edit-res-ip-${index}" value="${ip}"></td>
        <td><input type="text" class="form-control form-control-sm" id="edit-res-host-${index}" value="${hostname}"></td>
        <td><input type="text" class="form-control form-control-sm" id="edit-res-comment-${index}" value="${comment}"></td>
        <td>
            <i class="bi bi-check-lg text-success me-3" style="cursor: pointer;" onclick="saveReservation(${index}, '${mac}', '${ip}', '${hostname}', '${comment}')"></i>
            <i class="bi bi-x-circle text-secondary" style="cursor: pointer;" onclick="cancelReservationEdit()"></i>
        </td>
    `;
}

window.cancelReservationEdit = function() {
    displayReservations();
}

window.saveReservation = async function(index, oldMac, oldIp, oldHostname, oldComment) {
    const newMac = document.getElementById(`edit-res-mac-${index}`).value;
    const newIp = document.getElementById(`edit-res-ip-${index}`).value;
    const newHostname = document.getElementById(`edit-res-host-${index}`).value;
    const newComment = document.getElementById(`edit-res-comment-${index}`).value;

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

    await fetch(`${window.env.API_URL}/api/v1/dhcp/reservations`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            old: { mac_address: oldMac, ip_address: oldIp, hostname: oldHostname, comment: oldComment },
            new: { mac_address: newMac, ip_address: newIp, hostname: newHostname, comment: newComment }
        }),
    });

    displayReservations();
    showRestartBanner();
}

window.deleteReservation = async function(mac, ip, hostname) {
    if (!confirm(`Are you sure you want to delete reservation for ${mac} (${ip})?`)) return;

    await fetch(`${window.env.API_URL}/api/v1/dhcp/reservations`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ mac_address: mac, ip_address: ip, hostname: hostname }),
    });

    displayReservations();
    showRestartBanner();
}

// Initial load
displayReservations();
