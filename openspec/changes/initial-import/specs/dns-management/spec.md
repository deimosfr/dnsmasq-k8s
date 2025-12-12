## ADDED Requirements

### Requirement: Manage A Records
The system SHALL provide the ability to manage DNS A records.

#### Scenario: List A Records
- **WHEN** the user navigates to the DNS page
- **THEN** a list of existing A records is displayed from the configuration

#### Scenario: Add A Record
- **WHEN** the user enters a valid hostname and IP address
- **THEN** the record is added to the configuration
- **AND** the configuration is persisted

#### Scenario: Delete A Record
- **WHEN** the user deletes an A record
- **THEN** the record is removed from the configuration

### Requirement: Manage CNAME Records
The system SHALL provide the ability to manage DNS CNAME records.

#### Scenario: Add CNAME Record
- **WHEN** the user enters a valid alias and target domain
- **THEN** the record is added to the configuration

### Requirement: Manage TXT Records
The system SHALL provide the ability to manage DNS TXT records.

#### Scenario: Add TXT Record
- **WHEN** the user enters a valid hostname and text value
- **THEN** the record is added to the configuration
