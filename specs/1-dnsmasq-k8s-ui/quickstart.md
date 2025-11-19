# Quickstart: Dnsmasq Kubernetes UI

## Prerequisites

- A Kubernetes cluster
- Helm installed

## Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/user/dnsmasq-k8s-ui.git
    ```
2.  Navigate to the chart directory:
    ```bash
    cd dnsmasq-k8s-ui/chart
    ```
3.  Install the Helm chart:
    ```bash
    helm install dnsmasq-k8s-ui .
    ```

## Accessing the Web Interface

1.  Get the service URL:
    ```bash
    kubectl get svc dnsmasq-k8s-ui
    ```
2.  Open the URL in your browser.

## Usage

- Use the web interface to manage your dnsmasq configuration.
- The configuration is stored in a Kubernetes configmap.
- Changes to the configuration are automatically applied to the dnsmasq service.
