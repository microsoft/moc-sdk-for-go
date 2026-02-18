# Contributing

## How to Build

Install mingw-w64 by runnin below command if you are using WSL

`sudo apt-get install mingw-w64`

Then run 

`make`

Note: Try running `go mod tidy` if there are any issues with go packages when building

## Automated Workflows

### Go Version Updates

This repository uses a GitHub Actions workflow that automatically checks for new Go versions on the first day of each month. When a new stable Go version is released, the workflow:

1. Updates the Go version in `go.mod`
2. Updates the Go version in `.pipelines/build.yaml` (Azure Pipelines configuration)
3. Runs `go mod tidy` to update dependencies
4. Creates a pull request with the changes

The workflow can also be triggered manually from the Actions tab in GitHub if an immediate update is needed.

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
