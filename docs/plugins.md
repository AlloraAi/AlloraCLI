# Plugin Development Guide

## Overview
AlloraCLI supports a robust plugin system that allows you to extend functionality with custom providers, agents, and commands.

## Plugin Architecture

Plugins are Go packages that implement specific interfaces:

- **Agent Plugins**: Custom AI agents
- **Provider Plugins**: Cloud provider integrations
- **Command Plugins**: New CLI commands

## Creating an Agent Plugin

```go
package myplugin

import (
    "context"
    "github.com/AlloraAi/AlloraCLI/pkg/agents"
    "github.com/AlloraAi/AlloraCLI/pkg/config"
)

type MyAgent struct {
    agents.BaseAgent
}

func NewMyAgent(cfg config.Agent) agents.Agent {
    return &MyAgent{
        BaseAgent: agents.BaseAgent{
            Name: "my-agent",
            Config: cfg,
        },
    }
}

func (a *MyAgent) Query(ctx context.Context, query *agents.Query) (*agents.Response, error) {
    // Implement your custom logic here
    return &agents.Response{
        Text: "Custom response",
        Confidence: 0.9,
    }, nil
}
```

## Plugin Registration

Create a `plugin.yaml` file:

```yaml
name: my-plugin
version: 2.0.0
description: "My custom plugin"
author: "Your Name"
main: "./main.go"
dependencies:
  - "github.com/AlloraAi/AlloraCLI"
```

## Building Plugins

```bash
go build -buildmode=plugin -o my-plugin.so main.go
```

## Installing Plugins

```bash
allora plugin install ./my-plugin.so
```

## Plugin Examples

See the [examples/plugins](../examples/plugins) directory for complete examples.
