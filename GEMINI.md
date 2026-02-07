# Gemini Code Assistant Context

This document provides context for the Gemini Code Assistant to understand the `cagent` project.

## Project Overview

`cagent` is a powerful, easy-to-use, customizable multi-agent runtime written in Go. It orchestrates AI agents with specialized capabilities and tools, and the interactions between agents. Agents are configured using YAML files.

The project is structured as a Go application with a command-line interface (CLI) built using the `cobra` library. The main entry point is `main.go`, which calls the `root.Run()` function in `cmd/root/root.go`.

### Key Features

-   **Multi-agent architecture:** Create specialized agents for different domains.
-   **Rich tool ecosystem:** Agents can use external tools and APIs via the Model Context Protocol (MCP).
-   **Smart delegation:** Agents can automatically route tasks to the most suitable specialist.
-   **YAML configuration:** Declarative model and agent configuration. Numerous examples can be found in the `/examples` directory.
-   **Advanced reasoning:** Built-in "think", "todo" and "memory" tools for complex problem-solving.
-   **Multiple AI providers:** Support for OpenAI, Anthropic, Gemini and Docker Model Runner.

## Project Structure

-   `main.go`: The main entry point for the application.
-   `cmd/`: Contains the CLI command definitions, built using the `cobra` library.
-   `pkg/`: Contains the core application logic, including the agent runtime, model interactions, and tool handling.
-   `examples/`: Provides a wide variety of example agent configurations in YAML format, demonstrating different use cases and features.
-   `docs/`: Contains detailed documentation on usage, contributing, and telemetry.

## Building and Running

The project uses `Taskfile.yml` for task automation.

### Key Commands

-   **Build:** `task build` - Compiles the Go source code and creates a binary in the `bin` directory.
-   **Test:** `task test` - Runs the test suite.
-   **Lint:** `task lint` - Runs the `golangci-lint` linter.
-   **Format:** `task format` - Formats the Go source code.
-   **Run:** `cagent run <agent-file.yaml>` - Runs an agent from a YAML configuration file.

## Development Conventions

-   **Code Style:** The project follows standard Go formatting and linting practices, enforced by `gofumpt` and `golangci-lint`.
-   **Testing:** Tests are written using the standard Go testing framework.
-   **Configuration:** Agents and models are configured using YAML files. Examples can be found in the `/examples` directory.
-   **Documentation:** The `docs` directory contains detailed documentation on usage, contributing, and telemetry.