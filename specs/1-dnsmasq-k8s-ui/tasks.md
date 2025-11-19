# Tasks: Dnsmasq Kubernetes UI

**Input**: Design documents from `specs/1-dnsmasq-k8s-ui/`

## Phase 1: Setup (Shared Infrastructure)

- [X] T001 Create project structure per implementation plan
- [X] T002 Initialize Go project with Gin framework in `backend/`
- [X] T003 [P] Setup frontend project with Bootstrap in `frontend/`
- [X] T004 [P] Create Dockerfile based on debian slim stable

---

## Phase 2: Foundational (Blocking Prerequisites)

- [X] T005 Setup Kubernetes client-go in `backend/`
- [X] T006 [P] Implement basic API server in `backend/src/api/main.go`
- [X] T007 [P] Create basic frontend layout in `frontend/src/index.html`

---

## Phase 3: User Story 1 - View Dnsmasq Configuration (Priority: P1) 🎯 MVP

**Goal**: View the current dnsmasq configuration through a web interface.

**Independent Test**: The web interface should display the current configuration.

### Tests for User Story 1

- [X] T008 [P] [US1] Backend test for reading configuration from ConfigMap in `backend/src/services/config_test.go`
- [X] T009 [P] [US1] Frontend test for displaying configuration in `frontend/tests/config_test.js`

### Implementation for User Story 1

- [X] T010 [US1] Implement service to read configuration from ConfigMap in `backend/src/services/config.go`
- [X] T011 [US1] Implement API endpoint to get configuration in `backend/src/api/config.go`
- [X] T012 [P] [US1] Create frontend component to display configuration in `frontend/src/components/config.js`
- [X] T013 [US1] Integrate frontend component with API endpoint in `frontend/src/pages/config.js`

---

## Phase 4: User Story 4 - Deploy with Helm (Priority: P1)

**Goal**: Deploy the entire solution using a Helm chart.

**Independent Test**: The Helm chart should successfully deploy the application.

### Implementation for User Story 4

- [X] T014 [US4] Create Helm chart in `chart/`
- [X] T015 [US4] Add templates for deployment, service, and configmap to `chart/templates/`
- [X] T016 [US4] Add default values to `chart/values.yaml`

---

## Phase 5: User Story 2 - Update Dnsmasq Configuration (Priority: P2)

**Goal**: Update the dnsmasq configuration through the web interface.

**Independent Test**: The web interface should allow updating the configuration.

### Tests for User Story 2

- [X] T017 [P] [US2] Backend test for updating configuration in ConfigMap in `backend/src/services/config_test.go`
- [X] T018 [P] [US2] Frontend test for updating configuration in `frontend/tests/config_test.js`

### Implementation for User Story 2

- [X] T019 [US2] Implement service to update configuration in ConfigMap in `backend/src/services/config.go`
- [X] T020 [US2] Implement API endpoint to update configuration in `backend/src/api/config.go`
- [X] T021 [P] [US2] Create frontend component to edit configuration in `frontend/src/components/config-edit.js`
- [X] T022 [US2] Integrate frontend component with API endpoint in `frontend/src/pages/config.js`

---

## Phase 6: User Story 5 - Advanced Configuration (Priority: P2)

**Goal**: Edit the raw dnsmasq configuration file.

**Independent Test**: The web interface should provide a text editor to edit the raw configuration.

### Tests for User Story 5

- [X] T023 [P] [US5] Backend test for updating raw configuration in `backend/src/services/config_test.go`
- [X] T024 [P] [US5] Frontend test for editing raw configuration in `frontend/tests/advanced_config_test.js`

### Implementation for User Story 5

- [X] T025 [US5] Implement service to update raw configuration in `backend/src/services/config.go`
- [X] T026 [US5] Implement API endpoint to update raw configuration in `backend/src/api/config.go`
- [X] T027 [P] [US5] Create frontend component with a text editor in `frontend/src/components/advanced-config.js`
- [X] T028 [US5] Integrate frontend component with API endpoint in `frontend/src/pages/advanced-config.js`

---

## Phase 7: User Story 3 - Manage DHCP Leases (Priority: P3)

**Goal**: View and manage DHCP leases through the web interface.

**Independent Test**: The web interface should display and allow management of DHCP leases.

### Tests for User Story 3

- [X] T029 [P] [US3] Backend test for reading DHCP leases in `backend/src/services/dhcp_test.go`
- [X] T030 [P] [US3] Frontend test for displaying DHCP leases in `frontend/tests/dhcp_test.js`

### Implementation for User Story 3

- [X] T031 [US3] Implement service to read DHCP leases in `backend/src/services/dhcp.go`
- [X] T032 [US3] Implement API endpoint to get DHCP leases in `backend/src/api/dhcp.go`
- [X] T033 [P] [US3] Create frontend component to display DHCP leases in `frontend/src/components/dhcp-leases.js`
- [X] T034 [US3] Integrate frontend component with API endpoint in `frontend/src/pages/dhcp.js`

---

## Phase 8: Polish & Cross-Cutting Concerns

- [X] T035 [P] Add documentation to the `docs/` directory
- [X] T036 Code cleanup and refactoring
- [X] T037 [P] Add more unit tests to increase code coverage
- [X] T038 Run quickstart.md validation

---

## Dependencies & Execution Order

- **Setup (Phase 1)** -> **Foundational (Phase 2)**
- **Foundational (Phase 2)** -> **User Story 1 (Phase 3)**, **User Story 4 (Phase 4)**
- **User Story 1 (Phase 3)** -> **User Story 2 (Phase 5)**
- **User Story 2 (Phase 5)** -> **User Story 5 (Phase 6)**
- **User Story 1 (Phase 3)** -> **User Story 3 (Phase 7)**
- All user stories -> **Polish (Phase 8)**
