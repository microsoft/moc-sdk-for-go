# GitHub Copilot Instructions for moc-sdk-for-go

This file contains instructions for GitHub Copilot agents working on this repository.

## Pull Request Status Checks

**CRITICAL**: Before marking a Pull Request as ready for review, you MUST:

1. **Verify all status checks are passing (green)**
   - Check that Azure Pipelines build is successful
   - Verify CodeQL security scanning passes
   - Ensure all linting checks pass
   - Confirm all unit tests pass

2. **Run local validation before pushing changes**
   ```bash
   # Run the full build, format, and test suite
   make
   
   # Or run individual checks:
   make vendor      # Update dependencies
   make format      # Format code
   make build       # Build all packages and wrapper
   make unittest    # Run unit tests
   make golangci-lint  # Run linting
   ```

3. **Monitor CI/CD pipeline results**
   - After pushing changes, wait for all CI/CD workflows to complete
   - Review any failures in Azure Pipelines or GitHub Actions
   - Fix any issues before requesting review

4. **Key build requirements**
   - This project requires Go (version specified in `go.mod`)
   - Windows builds require mingw-w64 for cross-compilation
   - All code must pass golangci-lint checks
   - Unit tests must pass for `./pkg/client/...` and `./services/security/...`

## Status Check Sources

The following status checks must be green before marking PR as ready for review:

### Azure Pipelines
- **Build Job**: Compiles all packages including Windows DLL wrapper
- **Lint Job**: Runs golangci-lint with `.golangci.yml` configuration
- **Static Analysis**: Security and code quality checks

### GitHub Actions
- **CodeQL Analysis**: Security vulnerability scanning for Go code
- **CLA Check**: Contributor License Agreement verification

## Development Guidelines

1. **Never mark a PR as ready for review if any status checks are failing**
2. **Always run local tests before pushing** to catch issues early
3. **Review pipeline logs** if checks fail to understand the root cause
4. **Fix all build and test failures** before requesting review
5. **Ensure code is properly formatted** using `make format`

## Common Commands

```bash
# Full build and test (recommended before pushing)
make

# Individual operations
make vendor          # Update Go modules
make format          # Format Go code with gofmt
make build           # Build all packages
make test            # Run all tests
make unittest        # Run unit tests only
make golangci-lint   # Run linter

# Clean build artifacts
make clean
```

## Notes for Automated PRs

For automated PRs (e.g., dependency updates, Go version updates):
- The PR may be created in draft mode by default
- Before converting to ready for review, verify:
  - All builds complete successfully
  - All tests pass
  - No new linting errors are introduced
  - The changes are minimal and focused

## Repository-Specific Requirements

- Uses Go modules (`GO111MODULE=on`)
- Private repos require `GOPRIVATE=github.com/microsoft`
- Cross-compiles for Windows using mingw-w64
- Generates SBOM (Software Bill of Materials) as part of build
- Follows Microsoft Open Source Code of Conduct

---

**Remember**: Green status checks are mandatory before PR approval. No exceptions.
