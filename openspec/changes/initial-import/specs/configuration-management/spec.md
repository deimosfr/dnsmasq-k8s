## ADDED Requirements

### Requirement: View Configuration
The system SHALL provide a read-only view of the current `dnsmasq.conf` configuration.

#### Scenario: View Config
- **WHEN** the user navigates to the configuration page
- **THEN** the current content of the configuration file is displayed

### Requirement: Edit Configuration
The system SHALL allow authorized users to edit the global `dnsmasq` configuration.

#### Scenario: Edit Config
- **WHEN** the user enables edit mode
- **AND** submits valid configuration changes
- **THEN** the changes are saved to the configuration

### Requirement: Validate Configuration Configuration
The system SHALL validate the syntax of the configuration before saving.

#### Scenario: Validation Success
- **WHEN** the user submits a valid configuration
- **THEN** the save operation proceeds

#### Scenario: Validation Failure
- **WHEN** the user submits an invalid configuration (e.g., syntax error)
- **THEN** the system rejects the save
- **AND** displays an error message
