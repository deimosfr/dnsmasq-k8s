# Security

## ADDED Requirements

#### Requirement: Basic Authentication
The system MUST support Basic Authentication (RFC 7617) for all API access.
Credentials MUST be validated against a configured list of users.
Static assets (HTML, JS, CSS) SHOULD be accessible without authentication.

#### Scenario: Valid Credentials
Given the system is configured with user "admin" and password "secret"
When a request is made with "Authorization: Basic YWRtaW46c2VjcmV0"
Then the request SHOULD be allowed.

#### Scenario: Invalid Credentials
Given the system is configured with auth enabled
When a request is made with invalid credentials or no credentials
Then the system MUST respond with 401 Unauthorized
And the response MUST NOT include "WWW-Authenticate" header (to avoid browser prompt)
And the frontend SHOULD redirect the user to the login page.

#### Requirement: Disabled Authentication
The system MUST allowed disabling authentication.

#### Scenario: Authentication Disabled
Given authentication is disabled in configuration
When a request is made without credentials
Then the request SHOULD be allowed.
