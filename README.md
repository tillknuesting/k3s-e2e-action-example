# k3s-e2e-action-example

A GitHub Actions example that runs **end-to-end (E2E) tests against a local [K3s](https://k3s.io/) cluster** as part of CI.

## What it does

The [`E2E Tests (k3s)`](.github/workflows/e2e-test-k3s.yml) workflow runs the full pipeline on every relevant event:

1. Checks out the repository and sets up Go and Docker Buildx.
2. **On push to `main`:** builds the operator image and pushes it to the GitHub Container Registry (GHCR).
3. Starts a local Docker registry and a local K3s cluster via [`testing/k3s-setup.sh`](testing/k3s-setup.sh).
4. Waits for the Kubernetes API server and all nodes to become ready.
5. **On push to `main`:** pulls the freshly built GHCR image into the local registry; **on pull requests:** builds the image into the local registry.
6. Deploys the `hello` app to the K3s cluster and waits for the rollout to finish.
7. Runs the Go E2E test via `go test`.

## Workflow triggers

- **`push` to `main`** (paths: `operator/**`, `testing/**`, workflow file) → full pipeline incl. GHCR publish.
- **`pull_request` to `main`** (same path filters) → full pipeline without publishing, using a freshly built local image.
- **`workflow_dispatch`** → manual runs.

## Repository layout

| Path | Purpose |
| --- | --- |
| `operator/` | The application under test (Docker build context). |
| `deployment/` | Kubernetes `Deployment` / `Service` manifests for the app under test. |
| `testing/go/` | Go module with the E2E test (`e2e_test.go`) and module deps. |
| `testing/k3s-setup.sh` | Starts the local registry and the K3s cluster. |
| `testing/registries.yaml` | Registry mirror config used by K3s. |
| `.github/workflows/e2e-test-k3s.yml` | The E2E GitHub Actions workflow. |

## E2E test

The Go test in [`testing/go/e2e_test.go`](testing/go/e2e_test.go) waits for the `hello` deployment rollout to complete and then verifies that the pod is running.

## Matrix testing

The workflow runs the E2E pipeline against a matrix of K3s versions —
`v1.24.7+k3s1`, `v1.25.3+k3s1` and `v1.26.4+k3s1` — via the job's
`strategy.matrix`. Each version runs the full pipeline (build → local
registry → K3s → deploy → E2E test) on its own runner; failures are
reported per version (`fail-fast: false`).

## Running it

- Open a pull request touching `operator/` or `testing/` (or the workflow itself) to see the full E2E pipeline run on `ubuntu-24.04`.
- Push to `main` to additionally publish the image to GHCR and run the same pipeline against the fresh image.

## Notes

- The cluster runs in the GitHub-hosted runner via the official K3s install script; no external cluster is needed.
- The local registry (`registry.localhost:5000`) is mirrored into K3s via `testing/registries.yaml`.
- The workflow waits (with timeouts) for the Kubernetes API server and nodes before deploying and testing.
