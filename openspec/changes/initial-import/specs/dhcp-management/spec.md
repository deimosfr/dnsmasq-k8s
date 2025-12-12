## ADDED Requirements

### Requirement: Manage Static Reservations
The system SHALL allow users to define static DHCP reservations mapping MAC addresses to IP addresses.

#### Scenario: List Reservations
- **WHEN** the user navigates to the DHCP page
- **THEN** a list of existing static reservations is displayed

#### Scenario: Add Reservation
- **WHEN** the user provides a valid MAC address, IP address, and hostname
- **THEN** the reservation is saved to the configuration

### Requirement: View Active Leases
The system SHALL display current DHCP leases.

#### Scenario: View Leases
- **WHEN** the user views the DHCP status
- **THEN** a table of active leases (IP, MAC, Hostname, Expiry) is shown
