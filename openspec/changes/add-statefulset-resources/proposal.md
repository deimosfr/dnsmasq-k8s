# Change: Add StatefulSet Resources Support

## Why
Currently, the StatefulSet template does not allow configuring CPU and memory resources (requests and limits). This prevents users from setting appropriate resource constraints and guarantees for the dnsmasq workload, which is a best practice for Kubernetes deployments.

## What Changes
- Add `resources` configuration block to `values.yaml` with safe defaults (or empty structure).
- Update `statefulset.yaml` template to inject the `resources` block into the container spec.

## Impact
- **Affected specs**: `workload`
- **Affected code**:
    - `chart/values.yaml`
    - `chart/templates/statefulset.yaml`
