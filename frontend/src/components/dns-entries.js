// Sort state
let dnsSortState = { column: 'domain', direction: 'asc' };

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
function updateSortIcons(column, direction) {
    ['type', 'domain', 'value'].forEach(col => {
        const icon = document.getElementById(`sort-icon-${col}`);
        if (icon) {
            if (col === column) {
                icon.className = direction === 'asc' ? 'bi bi-chevron-up' : 'bi bi-chevron-down';
            } else {
                icon.className = 'bi bi-chevron-expand';
            }
        }
    });
}

// Sort DNS entries
window.sortDNSEntries = function(column) {
    if (dnsSortState.column === column) {
        dnsSortState.direction = dnsSortState.direction === 'asc' ? 'desc' : 'asc';
    } else {
        dnsSortState.column = column;
        dnsSortState.direction = 'asc';
    }
    displayDNSEntries();
};

// Validation patterns
const VALIDATION_PATTERNS = {
    address: {
        pattern: /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
        placeholder: 'IP Address',
        title: 'Please enter a valid IPv4 address'
    },
    cname: {
        pattern: /^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$/,
        placeholder: 'Target Domain',
        title: 'Please enter a valid domain name'
    },
    txt: {
        pattern: null,
        placeholder: 'Text Value',
        title: 'Enter any text value'
    }
};

// Update form validation based on DNS type
function updateFormValidation(type) {
    const valueInput = document.getElementById('dns-value');
    const validation = VALIDATION_PATTERNS[type];
    
    if (validation) {
        if (validation.pattern) {
            valueInput.pattern = validation.pattern.source;
            valueInput.title = validation.title;
        } else {
            valueInput.removeAttribute('pattern');
            valueInput.title = validation.title;
        }
        valueInput.placeholder = validation.placeholder;
    }
}

// Type select change listener
const dnsTypeSelect = document.getElementById('dns-type');
dnsTypeSelect.addEventListener('change', (e) => {
    updateFormValidation(e.target.value);
});

// Form submission
document.getElementById('add-dns-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const dnsType = document.getElementById('dns-type').value;
    const domain = document.getElementById('dns-domain').value.trim();
    const value = document.getElementById('dns-value').value.trim();

    // Client-side validation
    const validation = VALIDATION_PATTERNS[dnsType];
    if (validation && validation.pattern && !validation.pattern.test(value)) {
        alert(`Invalid value for ${dnsType.toUpperCase()} record: ${validation.title}`);
        return;
    }

    try {
        await addDNSEntry(dnsType, domain, value);
        document.getElementById('add-dns-form').reset();
        // Reset to default validation (A record)
        updateFormValidation('address');
        displayDNSEntries();
        showRestartBanner();
    } catch (error) {
        console.error('Error in add DNS entry handler:', error);
        // Error already shown in addDNSEntry
    }
});

async function addDNSEntry(dnsType, domain, value) {
    const response = await fetch('/api/v1/dns/entries', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            type: dnsType,
            domain: domain,
            value: value,
        }),
    });

    if (!response.ok) {
        const errorData = await response.json();
        alert(`Error adding DNS entry: ${errorData.error || 'Unknown error'}`);
        throw new Error(errorData.error || 'Unknown error');
    }
}

async function getDNSEntries() {
    const response = await fetch('/api/v1/dns/entries');
    const data = await response.json();
    return data;
}

async function displayDNSEntries() {
    const entries = await getDNSEntries();
    const tbody = document.getElementById('dns-entries-table-body');
    tbody.innerHTML = '';

    const header = document.getElementById('dns-entries-header');
    if (header) {
        header.innerHTML = `Current DNS Entries <span class="badge rounded-pill bg-success ms-2">${entries ? entries.length : 0}</span>`;
    }

    if (!entries || entries.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" class="text-center">No DNS entries found</td></tr>';
        return;
    }

    // Sort entries
    const sortedEntries = sortData(entries, dnsSortState.column, dnsSortState.direction);

    sortedEntries.forEach((entry, index) => {
        const row = document.createElement('tr');
        row.dataset.index = index;
        const displayType = formatDnsType(entry.type);
        // Escape values for HTML attributes
        const escapedDomain = entry.domain.replace(/'/g, '&#39;');
        const escapedValue = entry.value.replace(/'/g, '&#39;');
        
        row.innerHTML = `
            <td>${displayType}</td>
            <td>${entry.domain}</td>
            <td>${entry.value}</td>
            <td>
                <i class="bi bi-pencil text-success me-3" style="cursor: pointer;" onclick="editEntry(${index}, '${entry.type}', '${escapedDomain}', '${escapedValue}')"></i>
                <i class="bi bi-x-lg text-danger" style="cursor: pointer;" onclick="deleteEntry('${entry.type}', '${escapedDomain}', '${escapedValue}')"></i>
            </td>
        `;
        tbody.appendChild(row);
    });

    // Update sort icons
    updateSortIcons(dnsSortState.column, dnsSortState.direction);
}

function formatDnsType(type) {
    if (type === 'address') return 'A';
    if (type === 'cname') return 'CNAME';
    if (type === 'txt') return 'TXT';
    return type.toUpperCase();
}

window.editEntry = function(index, type, domain, value) {
    const tbody = document.getElementById('dns-entries-table-body');
    const row = tbody.children[index];
    const displayType = formatDnsType(type);
    
    // Unescape HTML entities
    const unescapedDomain = domain.replace(/&#39;/g, "'");
    const unescapedValue = value.replace(/&#39;/g, "'");
    
    row.innerHTML = `
        <td>
            <span class="form-control-plaintext form-control-sm">${displayType}</span>
        </td>
        <td><input type="text" class="form-control form-control-sm" id="edit-domain-${index}" value="${unescapedDomain}"></td>
        <td><input type="text" class="form-control form-control-sm" id="edit-value-${index}" value="${unescapedValue}"></td>
        <td>
            <i class="bi bi-check-lg text-success me-3" style="cursor: pointer;" onclick="saveEntry(${index}, '${type}', '${domain}', '${value}')"></i>
            <i class="bi bi-x-circle text-secondary" style="cursor: pointer;" onclick="cancelEdit()"></i>
        </td>
    `;
}

window.cancelEdit = function() {
    displayDNSEntries();
}

window.saveEntry = async function(index, oldType, oldDomain, oldValue) {
    const newType = oldType; // Type cannot be changed
    const newDomain = document.getElementById(`edit-domain-${index}`).value.trim();
    const newValue = document.getElementById(`edit-value-${index}`).value.trim();

    // Validate based on type
    const validation = VALIDATION_PATTERNS[oldType];
    if (validation && validation.pattern && !validation.pattern.test(newValue)) {
        alert(`Invalid value for ${oldType.toUpperCase()} record: ${validation.title}`);
        return;
    }

    // Unescape HTML entities from old values
    const unescapedOldDomain = oldDomain.replace(/&#39;/g, "'");
    const unescapedOldValue = oldValue.replace(/&#39;/g, "'");

    try {
        await fetch('/api/v1/dns/entries', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                old: { type: oldType, domain: unescapedOldDomain, value: unescapedOldValue },
                new: { type: newType, domain: newDomain, value: newValue }
            }),
        });

        displayDNSEntries();
        showRestartBanner();
    } catch (error) {
        alert('Error updating DNS entry: ' + error.message);
    }
}

window.deleteEntry = async function(type, domain, value) {
    // Unescape HTML entities
    const unescapedDomain = domain.replace(/&#39;/g, "'");
    const unescapedValue = value.replace(/&#39;/g, "'");
    
    if (!confirm(`Are you sure you want to delete this entry?\n${formatDnsType(type)} ${unescapedDomain} ${unescapedValue}`)) {
        return;
    }

    await fetch('/api/v1/dns/entries', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            type: type,
            domain: unescapedDomain,
            value: unescapedValue,
        }),
    });

    displayDNSEntries();
    showRestartBanner();
}

// Initialize
displayDNSEntries();