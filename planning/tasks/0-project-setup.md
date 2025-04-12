# Task 0: Project Setup

## Objective
Set up the foundational structure and environment for the Beluga project, ensuring a smooth development process in Golang.

## Tasks

### 0.1 Initialize Git Repository
- Ensure the project is under version control.
- Verify the `.gitignore` file includes necessary patterns for Golang projects.

### 0.2 Set Up Go Modules
- Run `go mod init beluga` to initialize the Go module.
- Add any required dependencies using `go get`.

### 0.3 Directory Structure
- Create a standard Golang project structure:
  - `cmd/` for application entry points.
  - `pkg/` for reusable packages.
  - `internal/` for internal packages.
  - `configs/` for configuration files.
  - `scripts/` for automation scripts.
  - `docs/` for documentation.

### 0.4 Linting and Formatting
- Set up `golangci-lint` for linting.
- Ensure `gofmt` is used for consistent code formatting.

### 0.5 Testing Framework
- Configure a testing framework (e.g., `testing` package).
- Add a basic test file to validate the setup.

### 0.6 Continuous Integration
- Set up a CI pipeline (e.g., GitHub Actions) to:
  - Run tests.
  - Check code formatting and linting.

### 0.7 Documentation
- Add a `README.md` with project setup instructions.
- Include a `CONTRIBUTING.md` for guidelines on contributing to the project.

### 0.8 Environment Configuration
- Create a `.env.example` file for environment variables.
- Add `.env` to `.gitignore` to prevent sensitive data from being committed.

### 0.9 Initial Commit
- Commit the initial project setup to the repository.