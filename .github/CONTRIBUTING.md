# Contributing

**Grit** is open source software; contributions from the community are
encouraged and appreciated. Please take a moment to read these guidelines
before submitting changes.

## Requirements

- [Go 1.16](https://golang.org/)
- [GNU make](https://www.gnu.org/software/make/) (or equivalent)

## Running the tests

    make

The default target of the make file installs all necessary dependencies and runs
the tests.

Code coverage reports can be built with:

    make coverage

To rebuild coverage reports and open them in a browser, use:

    make coverage-open

## Submitting changes

Change requests are reviewed and accepted via pull-requests on GitHub. If you're
unfamiliar with this process, please read the relevant GitHub documentation
regarding [forking a repository](https://help.github.com/articles/fork-a-repo)
and [using pull-requests](https://help.github.com/articles/using-pull-requests).

Before submitting your pull-request (typically against the `master` branch),
please run:

    make prepare

To apply any automated code-style updates, run linting checks, run the tests and
build coverage reports. Please ensure that your changes are tested and that a
high level of code coverage is maintained.
