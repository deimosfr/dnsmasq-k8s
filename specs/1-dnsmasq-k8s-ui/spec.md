# Feature Specification: Dnsmasq Kubernetes UI

**Feature Branch**: `1-dnsmasq-k8s-ui`  
**Created**: 2025-11-19
**Status**: Draft  
**Input**: User description: "I want a container with dnsmasq that runs on Kubernetes. A web interface should be implemented to manage dnsmasq config, the DNS config and manual entries, the DHCP config and DHCP leases. The web interface should be able to read config, update and store in Kubernetes configmaps. It should watch as well the config files and reload dnsmasq config when something is updated. I also need a helm chart to deploy the global solution, pay attention to the fact that the chart should not override the config set with webui."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Dnsmasq Configuration (Priority: P1)

As a system administrator, I want to view the current dnsmasq configuration (including DNS and DHCP settings) through a web interface, so that I can easily check the state of the service.

**Why this priority**: This is the most basic functionality and is required for any further management tasks.

**Independent Test**: The web interface should display the current configuration, which can be verified against the configuration stored in the Kubernetes configmaps.

**Acceptance Scenarios**:

1. **Given** a running dnsmasq-k8s-ui instance, **When** I navigate to the configuration page, **Then** I should see the current dnsmasq, DNS, and DHCP configurations.
2. **Given** the configuration is updated in the configmap, **When** I refresh the configuration page, **Then** I should see the updated configuration.

### User Story 2 - Update Dnsmasq Configuration (Priority: P2)

As a system administrator, I want to update the dnsmasq configuration through the web interface, so that I can easily make changes without manually editing Kubernetes configmaps.

**Why this priority**: This is the core management functionality of the web interface.

**Independent Test**: After updating the configuration in the web interface, the corresponding Kubernetes configmap should be updated, and dnsmasq should reload the new configuration.

**Acceptance Scenarios**:

1. **Given** I am on the configuration page, **When** I change a configuration value and save it, **Then** the corresponding Kubernetes configmap should be updated.
2. **Given** the configuration is updated, **When** I check the dnsmasq service, **Then** it should be using the new configuration.
3. **Given** I enter an invalid configuration value, **When** I try to save it, **Then** the web interface should display an error message and prevent me from saving the invalid configuration.

### User Story 3 - Manage DHCP Leases (Priority: P3)

As a system administrator, I want to view and manage DHCP leases through the web interface, so that I can easily see which devices are on the network and manage their leases.

**Why this priority**: This is an important feature for managing a DHCP server.

**Independent Test**: The web interface should display the current DHCP leases, and I should be able to add, remove, and modify leases.

**Acceptance Scenarios**:

1. **Given** I am on the DHCP leases page, **When** a new device gets a lease, **Then** the new lease should appear in the list.
2. **Given** I am on the DHCP leases page, **When** I delete a lease, **Then** the lease should be removed from the list and the device should no longer be able to use that IP address.

### User Story 4 - Deploy with Helm (Priority: P1)

As a system administrator, I want to deploy the entire solution using a Helm chart, so that I can easily install and manage the application on a Kubernetes cluster.

**Why this priority**: This is essential for easy adoption and management of the solution.

**Independent Test**: The Helm chart should successfully deploy the dnsmasq container and the web interface.

**Acceptance Scenarios**:

1. **Given** I have a Kubernetes cluster, **When** I install the Helm chart, **Then** the dnsmasq and web interface pods should be running.
2. **Given** the application is deployed, **When** I update a value in the web interface, **Then** the configuration should not be overridden by the Helm chart on the next release.

### User Story 5 - Advanced Configuration (Priority: P2)

As an advanced user, I want to be able to edit the raw dnsmasq configuration file, so that I can access all possible configuration options.

**Why this priority**: This provides flexibility for users who need more control over the configuration.

**Independent Test**: The web interface should provide a text editor to edit the raw configuration file. After saving the changes, the configmap should be updated and dnsmasq should be reloaded.

**Acceptance Scenarios**:

1. **Given** I am on the advanced configuration page, **When** I make a change to the raw configuration and save it, **Then** the configmap should be updated.
2. **Given** the raw configuration is updated, **When** I check the dnsmasq service, **Then** it should be using the new configuration.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a web interface to manage dnsmasq configuration.
- **FR-002**: The web interface MUST be able to read the current dnsmasq configuration from Kubernetes configmaps.
- **FR-003**: The web interface MUST be able to update the dnsmasq configuration in Kubernetes configmaps.
- **FR-004**: The system MUST watch for changes in the configmaps and automatically reload the dnsmasq configuration.
- **FR-005**: The system MUST provide a Helm chart for deployment.
- **FR-006**: The Helm chart MUST NOT override configuration changes made through the web interface.
- **FR-007**: The web interface MUST allow management of DNS configuration, including manual entries.
- **FR-008**: The web interface MUST allow management of DHCP configuration and leases.
- **FR-009**: The system MUST run dnsmasq in a container on Kubernetes.
- **FR-010**: The web interface MUST provide a curated UI for common configuration options.
- **FR-011**: The web interface MUST provide a raw text editor for advanced users to edit the configuration directly.
- **FR-012**: The web interface MUST perform client-side validation of configuration changes.

### Edge Cases

- What happens when a user enters an invalid configuration in the raw text editor? The system should provide a way to validate the configuration before applying it.
- What happens if the configmap is edited manually while the web interface is open? The web interface should detect the change and prompt the user to reload the configuration.
- What happens if the dnsmasq service fails to reload after a configuration change? The system should report the error to the user and provide a way to revert to the last known good configuration.

### Key Entities *(include if feature involves data)*

- **DnsmasqConfig**: Represents the main dnsmasq configuration.
- **DNSConfig**: Represents the DNS configuration, including manual entries.
- **DHCPConfig**: Represents the DHCP configuration.
- **DHCPLease**: Represents a single DHCP lease.
-- **ConfigMap**: Kubernetes resource used to store configuration.

### Assumptions

- The web interface will be built using Go with the Gin framework.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A system administrator can view and update the entire dnsmasq configuration through the web interface in under 2 minutes.
- **SC-002**: The system can handle 100 concurrent users on the web interface without degradation.
- **SC-003**: 95% of configuration changes made through the web interface are applied to the dnsmasq service within 10 seconds.
- **SC-004**: The Helm chart can be used to deploy the entire solution to a Kubernetes cluster in under 5 minutes.