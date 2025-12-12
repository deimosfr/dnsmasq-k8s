# Project Context

## Purpose
Dnsmasq K8s is a web interface for managing dnsmasq DNS and DHCP services within Kubernetes clusters. It aims to simplify network management by providing an intuitive UI for editing DNS records, DHCP leases, and service configurations, with changes automatically persisted to Kubernetes ConfigMaps.

## Tech Stack
- **Frontend**:
    - HTML5
    - Vanilla JavaScript (ES6+ with Web Components)
    - CSS (Vanilla with variables)
    - Bootstrap 5.3 (via CDN)
- **Backend**:
    - Go 1.25
    - Gin Gonic (Web Framework)
    - Kubernetes Client-Go (K8s API interaction)
- **Infrastructure**:
    - Kubernetes
    - Helm
    - Docker
    - Supervisor (Process management)

## Project Conventions

### Code Style
- **Go**: Follows standard Go formatting (`gofmt`) and idioms.
- **JavaScript**:
    - Uses ES modules and classes.
    - Web Components for UI elements (e.g., `app-navbar`).
    - No build step for frontend code (served as static files).
- **CSS**:
    - Uses `theme.css` for custom styles and variables.
    - Relies on Bootstrap utility classes for layout and common components.

### Architecture Patterns
- **Backend**:
    - Layered architecture: `api` (routes), `controllers` (handlers), `services` (business logic), `models` (data structures).
    - Stateless design where possible, relying on Kubernetes ConfigMaps as the source of truth for configuration.
- **Frontend**:
    - Single Page Application (SPA) feel using simple page transitions or dynamic content loading where applicable, though primarily multi-page structure served by Go backend.
    - Direct API communication from frontend to backend endpoints.

### Testing Strategy
- **Backend**: Go unit tests (`_test.go` files) using `testing` package and `testify`.
- **Frontend**: Explicit testing strategy not observed in initial scan (likely manual or basic unit tests if present).

### Git Workflow
- Standard feature branch workflow.

## Domain Context
- **DNS**: Managing A, CNAME, TXT records. Understanding of `dnsmasq.conf` syntax.
- **DHCP**: Managing static reservations (MAC to IP) and viewing active leases.
- **Kubernetes Integration**: The application acts as a controller/interface for `dnsmasq` running in a pod, syncing state to ConfigMaps (`dnsmasq-custom-dns`, `dnsmasq-reservations`, `dnsmasq-leases`).

## Important Constraints
- **Host Network**: Often runs with `hostNetwork: true` to handle DHCP properly.
- **Capabilities**: Requires `NET_ADMIN` capability for network operations.
- **ConfigMaps**: Critical dependency; application reads/writes to specific ConfigMaps.

## External Dependencies
- **Kubernetes API**: Crucial for all state persistence.
- **Bootstrap CDN**: Frontend relies on external CDN for styles and JS.
