# Data Model: Dnsmasq Kubernetes UI

## DnsmasqConfig

Represents the main dnsmasq configuration.

| Field | Type | Description |
|---|---|---|
| `domain-needed` | boolean | Don't forward plain names (without a dot or domain part) |
| `bogus-priv` | boolean | Don't forward reverse lookups for private IP ranges |
| `no-resolv` | boolean | Don't read /etc/resolv.conf |
| `server` | string | Specify upstream DNS server |

## DNSConfig

Represents the DNS configuration, including manual entries.

| Field | Type | Description |
|---|---|---|
| `address` | string | A record (e.g., `/example.com/1.2.3.4`) |
| `cname` | string | CNAME record (e.g., `/alias.example.com/example.com`) |

## DHCPConfig

Represents the DHCP configuration.

| Field | Type | Description |
|---|---|---|
| `dhcp-range` | string | The range of IP addresses to lease |
| `dhcp-option` | string | DHCP options (e.g., `option:router,192.168.1.1`) |

## DHCPLease

Represents a single DHCP lease.

| Field | Type | Description |
|---|---|---|
| `mac-address` | string | The MAC address of the device |
| `ip-address` | string | The IP address leased to the device |
| `hostname` | string | The hostname of the device |
