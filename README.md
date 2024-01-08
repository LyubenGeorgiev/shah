## Workflows

### Dev Pipeline

This workflow is triggered on pushes to branches other than the `main` branch.

#### Job: `gosec`
- **Description**: Runs the Gosec Security Scanner on the codebase
  
#### Job: `gitleaks`
- **Description**: Executes Gitleaks to check for sensitive information leakage.
  
#### Job: `test`
- **Description**: Runs tests after security checks.
- **Dependencies**: Depends on the completion of `gosec` and `gitleaks` jobs.
  
#### Job: `build`
- **Description**: Builds Docker images, scans for vulnerabilities using Trivy
- **Dependencies**: Depends on the completion of the `test` job.

### Main Pipeline

This workflow is triggered on pushes to the `main` branch.

#### Jobs: `gosec`, `gitleaks`, `test`
- **Description**: Similar jobs as defined in the `Dev Pipeline`, executing security checks, testing.

#### Job: `build`
- **Description**: Builds Docker images, scans for vulnerabilities using Trivy. If successful uploads it to https://hub.docker.com/repository/docker/lyubengeorgiev/shah
- **Dependencies**: Depends on the completion of the `test` job.
  
#### Job: `trigger-infra`
- **Description**: Triggers deployment workflow in the 'shah-infra' repository after successful build and push to dockerhub.
- **Dependencies**: Depends on the completion of the `build` job.

### Pipeline Overview

These workflows automate the CI/CD processes, including code quality checks, security scans, Docker image building, vulnerability scanning, and deployment triggers. The `Dev Pipeline` focuses on feature branch development, while the `Main Pipeline` ensures the `main` branch's integrity and triggers deployments.

For more details about each job and workflow configurations, refer to the corresponding workflow YAML files in this repository.
