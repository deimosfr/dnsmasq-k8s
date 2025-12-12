## ADDED Requirements

### Requirement: Service Status
The system SHALL display the status of managed services (dnsmasq, supervisor).

#### Scenario: View Status
- **WHEN** the user views the dashboard
- **THEN** the current running status of services is shown

### Requirement: Control Services
The system SHALL allow users to start, stop, and restart the dnsmasq service.

#### Scenario: Restart Service
- **WHEN** the user clicks the "Restart" action for dnsmasq
- **THEN** the service is restarted via supervisor
- **AND** the new status is reflected in the UI
