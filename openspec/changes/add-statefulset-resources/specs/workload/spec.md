## ADDED Requirements
### Requirement: Configurable Resources
The chart SHALL allow configuring Kubernetes resource requests and limits (CPU/Memory) for the main dnsmasq container via values.

#### Scenario: Resources are customized
- **WHEN** `resources` object is provided in values (e.g., limits: cpu: 100m)
- **THEN** the StatefulSet manifest contains the corresponding resources block in the container spec
