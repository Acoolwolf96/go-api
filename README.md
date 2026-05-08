# CI/CD Pipeline on Self-Hosted Infrastructure



A working CI/CD pipeline that runs on my own hardware. When I push code to GitHub, it automatically gets tested, packaged into a container, and deployed to a Kubernetes cluster.



## Components



| Component | What it does |

| Go | REST API code |

| Docker | Container packaging |

| Docker Registry | Stores images |

| K3s | Runs containers |

| GitHub Actions | Runs tests, builds, pushes |

| Argo CD | Syncs cluster with git |



## Repositories



| Repo | Purpose |

| `go-api` | Application source code, tests, Dockerfile, CI workflow |

| `go-api-gitops` | Kubernetes manifests — source of truth for cluster state |



## Pipeline Flow



1. Code is pushed to `main` on the `go-api` GitHub repository

2. GitHub Actions triggers the self-hosted runner on the VM

3. Tests run with `go test ./...` — pipeline stops here if any test fails

4. Docker image is built and pushed to the private registry with two tags: the full git SHA and `latest`

5. The CI pipeline clones the `go-api-gitops` manifest repository and updates `deployment.yaml` with the exact image SHA

6. The updated manifest is committed and pushed back to `go-api-gitops`

7. Argo CD detects the manifest change and automatically syncs the cluster

8. K3s pulls the new image from the private registry and rolls out the updated deployment.



## Key Concepts Demonstrated



- **GitOps**: The cluster state is fully defined in a Git repository. Argo CD ensures the cluster always matches the manifest.

- **Immutable deployments**: Every deployment uses a specific image SHA, not a mutable tag. Full traceability from commit to running pod.

- **Self-hosted infrastructure**: The entire pipeline runs on owned hardware

- **Separation of concerns**: Application code and deployment manifests live in separate repositories.

