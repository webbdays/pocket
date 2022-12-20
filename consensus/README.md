# Consensus Module <!-- omit in toc -->

This document is meant to be a supplement to the living specification of [1.0 Pocket's P2P Specification](https://github.com/pokt-network/pocket-network-protocol/tree/main/consensus) primarily focused on the implementation, and additional details related to the design of the codebase and information related to development.

## Table of Contents <!-- omit in toc -->

- [Interface](#interface)
- [Implementation](#implementation)
  - [Code Architecture - Consensus Module](#code-architecture---p2p-module)
  - [Code Organization](#code-organization)
- [Testing](#testing)
  - [Running Unit Tests](#running-unit-tests)

## Interface

This module aims to implement the interface specified in `pocket/shared/modules/consensus_module.go` using the specification above.

## Implementation

### Code Architecture - P2P Module

```mermaid

```

### Code Organization

```bash
consensus
├── README.md                               # Self link to this README                   
├── module.go                               # The implementation of the Consensus Interface
├── types
│   ├── consensus_genesis.go              # implementation
│   ├── converters.go         
│   ├── errors.go               
│   ├── types.go                          # implementation of the Network interface using RainTree's specification
│   ├── utils.go
│   └── proto
│       └── block.proto
│       └── consensus_config.proto
│       └── consensus_genesis.proto
│       └── hotstuff_types.proto                 
└── utils.go
```

## Testing

_TODO: The work to add the tooling used to help with unit test generation is being tracked in #314._

### Running Unit Tests

```bash
make test_consensus
```

