# Changelog

All notable changes to this module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.0.3] - 2023-01-11

- Reorder private keys so addresses (retrieved by transforming private keys) to reflect the numbering in LocalNet appropriately. The address for val1, based on config1, will have the lexicographically first address. This makes debugging easier.

## [0.0.0.2] - 2023-01-10

- Removed `BaseConfig` from `configs`
- Centralized `PersistenceGenesisState` and `ConsensusGenesisState` into `GenesisState`
- Removed `is_client_only` since it's set programmatically in the CLI

## [0.0.0.1] - 2022-12-29

- Updated all `config*.json` files with the missing `max_mempool_count` value
- Added `is_client_only` to `config1.json` so Viper knows it can be overridden. The config override is done in the Makefile's `client_connect` target. Setting this can be avoided if we merge the changes in https://github.com/pokt-network/pocket/compare/main...issue/cli-viper-environment-vars-fix

## [0.0.0.0] - 2022-12-22

- Introduced this `CHANGELOG.md`
