# Contributing to [fgm/container]

We welcome contributions to fgm/container! 
Whether you're fixing a bug, proposing a feature, or improving documentation, your help is appreciated.


## Getting Started

1.  Fork the repository
2.  Clone your fork
3.  Create a branch
4.  Make your changes
6.  Push to your fork
7.  Create a pull request to [fgm/container]


## Contribution Guidelines

* **Code Style:** Please follow the existing code style. We use [staticcheck], and strive to use idiomatic Go code, following the [Effective Go] best practices.
* **Tests:** If you're adding or changing code, please add or update tests as needed. Unit tests are needed, benchmarks are important if your implementation provides new implementations. 
 [Fuzz tests] are nice to have too.
* **Documentation:** Update the documentation to reflect your changes.
* **Commit Messages:** Use clear and concise commit messages. Preferably follow the [Conventional Commits] format
* **Issue Tracking:** If you're fixing a bug, please reference the issue number in your pull request description.


## Reporting Bugs

- **Check existing issues:** Before creating a new issue, please check if the bug has already been reported, at [issues]
- **Provide detailed information:** Include the following in your bug report:
  * Steps to reproduce the bug.
  * Expected behavior.
  * Actual behavior.
  * Your processor family (arm64, amd64, ...)
  * Your operating system (Linux, darwin, ...) and version.
  * Any relevant error messages or logs.


## Suggesting Features

- **Check existing issues:** Before suggesting a new feature, please check if it has already been suggested, at [issues]
- **Provide a clear description:** Explain the feature you'd like to see and why it would be useful.
- **Consider implementation:** If possible, provide some ideas on how the feature could be implemented. Example in other languages can be useful.


## Development Setup

- **Install dependencies:**
  - This project has zero run-time dependencies,
and minimizes test-time dependencies.
  - As such, code or tests MUST NOT use non-stdlib packages.
  - Follow instructions on the [Staticcheck] site to install the linter. 
- **Run tests:** 
  - Currently, this command will run all unit and benchmark tests:
    ```bash
    go test -race -run=. -bench=. ./... 
    ```
  - If you add fuzz tests, update the documentation to describe how to run them.
- **Run the demos:**
    ```bash
    go run ./cmd/orderedmap
    go run ./cmd/queuestack
    ```
## Security Issues

Follow instructions in [SECURITY.md](SECURITY.md).

Do not disclose security vulnerabilities publicly. 

## Code of Conduct

Please note that this project is released without a Code of Conduct at this point,
but one could be added at any time.

Thank you for contributing to [fgm/container]!


[fgm/container]: https://github.com/fgm/container
[staticcheck]: https://staticcheck.io
[Effective Go]: https://go.dev/doc/effective_go
[Conventional Commits]:https://www.conventionalcommits.org/en/v1.0.0
[Fuzz tests]: https://go.dev/doc/security/fuzz/
[issues]: https://github.com/fgm/container/issues
