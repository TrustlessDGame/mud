# Change Log

## 2.0.0-next.0

### Minor Changes

- [#1147](https://github.com/latticexyz/mud/pull/1147) [`66cc35a8`](https://github.com/latticexyz/mud/commit/66cc35a8ccb21c50a1882d6c741dd045acd8bc11) Thanks [@dk1a](https://github.com/dk1a)! - Create gas-report package, move gas-report cli command and GasReporter contract to it

- [#1061](https://github.com/latticexyz/mud/pull/1061) [`a7b30c79`](https://github.com/latticexyz/mud/commit/a7b30c79bcc78530d2d01858de46a0fb87954fda) Thanks [@dk1a](https://github.com/dk1a)! - Rename `MudV2Test` to `MudTest` and move from `@latticexyz/std-contracts` to `@latticexyz/store`.

  ```solidity
  // old import
  import { MudV2Test } from "@latticexyz/std-contracts/src/test/MudV2Test.t.sol";
  // new import
  import { MudTest } from "@latticexyz/store/src/MudTest.sol";
  ```

  Refactor `StoreSwitch` to use a storage slot instead of `function isStore()` to determine which contract is Store:

  - Previously `StoreSwitch` called `isStore()` on `msg.sender` to determine if `msg.sender` is a `Store` contract. If the call succeeded, the `Store` methods were called on `msg.sender`, otherwise the data was written to the own storage.
  - With this change `StoreSwitch` instead checks for an `address` in a known storage slot. If the address equals the own address, data is written to the own storage. If it is an external address, `Store` methods are called on this address. If it is unset (`address(0)`), store methods are called on `msg.sender`.
  - In practice this has the same effect as before: By default the `World` contracts sets its own address in `StoreSwitch`, while `System` contracts keep the Store address undefined, so `Systems` write to their caller (`World`) if they are executed via `call` or directly to the `World` storage if they are executed via `delegatecall`.
  - Besides gas savings, this change has two additional benefits:
    1. it is now possible for `Systems` to explicitly set a `Store` address to make them exclusive to that `Store` and
    2. table libraries can now be used in tests without having to provide an explicit `Store` argument, because the `MudTest` base contract redirects reads and writes to the internal `World` contract.

### Patch Changes

- [#1075](https://github.com/latticexyz/mud/pull/1075) [`904fd7d4`](https://github.com/latticexyz/mud/commit/904fd7d4ee06a86e481e3e02fd5744224376d0c9) Thanks [@holic](https://github.com/holic)! - Add store sync package

- [#1153](https://github.com/latticexyz/mud/pull/1153) [`8d51a034`](https://github.com/latticexyz/mud/commit/8d51a03486bc20006d8cc982f798dfdfe16f169f) Thanks [@dk1a](https://github.com/dk1a)! - Clean up Memory.sol, make mcopy pure

- [#1168](https://github.com/latticexyz/mud/pull/1168) [`48909d15`](https://github.com/latticexyz/mud/commit/48909d151b3dfceab128c120bc6bb77de53c456b) Thanks [@dk1a](https://github.com/dk1a)! - bump forge-std and ds-test dependencies

- [#1179](https://github.com/latticexyz/mud/pull/1179) [`53522998`](https://github.com/latticexyz/mud/commit/535229984565539e6168042150b45fe0f9b48b0f) Thanks [@holic](https://github.com/holic)! - - bump to viem 1.3.0 and abitype 0.9.3
  - move `@wagmi/chains` imports to `viem/chains`
  - refine a few types
- Updated dependencies [[`8d51a034`](https://github.com/latticexyz/mud/commit/8d51a03486bc20006d8cc982f798dfdfe16f169f), [`48909d15`](https://github.com/latticexyz/mud/commit/48909d151b3dfceab128c120bc6bb77de53c456b), [`66cc35a8`](https://github.com/latticexyz/mud/commit/66cc35a8ccb21c50a1882d6c741dd045acd8bc11), [`f03531d9`](https://github.com/latticexyz/mud/commit/f03531d97c999954a626ef63bc5bbae51a7b90f3), [`53522998`](https://github.com/latticexyz/mud/commit/535229984565539e6168042150b45fe0f9b48b0f), [`0c4f9fea`](https://github.com/latticexyz/mud/commit/0c4f9fea9e38ba122316cdd52c3d158c62f8cfee)]:
  - @latticexyz/common@2.0.0-next.0
  - @latticexyz/gas-report@2.0.0-next.0
  - @latticexyz/schema-type@2.0.0-next.0
  - @latticexyz/config@2.0.0-next.0

All notable changes to this project will be documented in this file.
See [Conventional Commits](https://conventionalcommits.org) for commit guidelines.

# [1.42.0](https://github.com/latticexyz/mud/compare/v1.41.0...v1.42.0) (2023-04-13)

### Bug Fixes

- **cli:** account for getRecord's trimming ([#574](https://github.com/latticexyz/mud/issues/574)) ([9c5317a](https://github.com/latticexyz/mud/commit/9c5317afb2c4a9ac2fbaca424f90f30575dba671))

### Features

- add support for key schemas ([#480](https://github.com/latticexyz/mud/issues/480)) ([37aec2e](https://github.com/latticexyz/mud/commit/37aec2e0a8adf378035fa5b54d752cc6888378d2))
- **cli:** add encode function to all tables ([#498](https://github.com/latticexyz/mud/issues/498)) ([564604c](https://github.com/latticexyz/mud/commit/564604c0c03d675e007d176ec735d8fb76976771))
- **cli:** add module config to CLI ([#494](https://github.com/latticexyz/mud/issues/494)) ([263c828](https://github.com/latticexyz/mud/commit/263c828d3eb6f43d5e635c28026f4a68fbf7a19b))
- **cli:** allow static arrays as abi types in store config and tablegen ([#509](https://github.com/latticexyz/mud/issues/509)) ([588d037](https://github.com/latticexyz/mud/commit/588d0370d4c7d13667ff784ecb170edf59aa119e))
- **cli:** improve storeArgument, refactor cli ([#500](https://github.com/latticexyz/mud/issues/500)) ([bb68670](https://github.com/latticexyz/mud/commit/bb686702da75401d9ea4a8c8effcf3a15fa53b49))
- **cli:** set storeArgument to true by default ([#553](https://github.com/latticexyz/mud/issues/553)) ([cb1ecbc](https://github.com/latticexyz/mud/commit/cb1ecbcd036ead1b1ba0b717c7531d15beaeb106))
- **cli:** use a central codegen dir for tablegen and worldgen ([#585](https://github.com/latticexyz/mud/issues/585)) ([7500b11](https://github.com/latticexyz/mud/commit/7500b119d727a7155fa1033b2fc3ca729a51d033))
- **cli:** use abi types in store config ([#507](https://github.com/latticexyz/mud/issues/507)) ([12a739f](https://github.com/latticexyz/mud/commit/12a739f953d0929f7ffc8657fa22bc9e68201d75))
- **cli:** use json for gas report output ([#607](https://github.com/latticexyz/mud/issues/607)) ([bea12ca](https://github.com/latticexyz/mud/commit/bea12cac16a2e0cdbb9623571cf0b02a5ed969a2))
- **config:** separate config from cli ([#600](https://github.com/latticexyz/mud/issues/600)) ([cd224a5](https://github.com/latticexyz/mud/commit/cd224a5244ee55316d4b95a21007a8076adefe6e))
- **store:** add metadata to the schema table ([#550](https://github.com/latticexyz/mud/issues/550)) ([55ab704](https://github.com/latticexyz/mud/commit/55ab704c36a8ba5fd021cc80abb32b3f69e97b73))
- use IErrors in IStore and IWorldCore ([#573](https://github.com/latticexyz/mud/issues/573)) ([4f9ed7b](https://github.com/latticexyz/mud/commit/4f9ed7ba22ea978623b6d54e9731081580c2ad8f))
- v2 event decoding ([#415](https://github.com/latticexyz/mud/issues/415)) ([374ed54](https://github.com/latticexyz/mud/commit/374ed542c3387a4ec2b36ab68ae534419aa58763))
- **world,store:** add updateInField ([#525](https://github.com/latticexyz/mud/issues/525)) ([0ac76fd](https://github.com/latticexyz/mud/commit/0ac76fd57484f54860157b79678b8b9eb7a86997))
- **world:** add naive ReverseMappingHook/Module ([#487](https://github.com/latticexyz/mud/issues/487)) ([36aaaef](https://github.com/latticexyz/mud/commit/36aaaef3a69914b962a3ef0847aa144134e89d28))
- **world:** add support for modules, add RegistrationModule, add CoreModule ([#482](https://github.com/latticexyz/mud/issues/482)) ([624cbbc](https://github.com/latticexyz/mud/commit/624cbbc6722823e83594f3df38d72682a1cecd99))
- **world:** add UniqueEntityModule ([#552](https://github.com/latticexyz/mud/issues/552)) ([983e26a](https://github.com/latticexyz/mud/commit/983e26a0ee0c0521e99d09dd25ebb9937e7c4ded))
- **world:** allow registration of function selectors in the World, split out RegisterSystem from World ([#481](https://github.com/latticexyz/mud/issues/481)) ([ba0166f](https://github.com/latticexyz/mud/commit/ba0166fb6cd7de63ddc6f4f500ff90c05da67b09))
- **world:** simplify access control to namespaces instead of routes ([#467](https://github.com/latticexyz/mud/issues/467)) ([945f2ef](https://github.com/latticexyz/mud/commit/945f2ef4a09c2fd1f9c4bb0418a1569fc31e0776))

# [1.41.0](https://github.com/latticexyz/mud/compare/v1.40.0...v1.41.0) (2023-03-09)

### Bug Fixes

- **cli:** add missing await to tablegen, fix formatting ([#472](https://github.com/latticexyz/mud/issues/472)) ([4313c27](https://github.com/latticexyz/mud/commit/4313c277b10c0334716e5c3728ffeaef643c1e6b))

### Features

- add pushToField to Store and World ([#434](https://github.com/latticexyz/mud/issues/434)) ([b665efc](https://github.com/latticexyz/mud/commit/b665efcb407992779c93a400199dee8ffdc03cb7))
- **cli:** add setMetadata to autogen of table libraries ([#466](https://github.com/latticexyz/mud/issues/466)) ([1e129fe](https://github.com/latticexyz/mud/commit/1e129fe9ced354e838d3d9afc9839aba82fbf210))
- **cli:** add v2 deployment script ([#450](https://github.com/latticexyz/mud/issues/450)) ([1db37a5](https://github.com/latticexyz/mud/commit/1db37a5c6b736fdc5f806653b78f76b02239f2bb))
- **cli:** user types and route/path separation ([#454](https://github.com/latticexyz/mud/issues/454)) ([758bf03](https://github.com/latticexyz/mud/commit/758bf0388c9e282c58b2890cadb4a59e00978d26))

# [1.40.0](https://github.com/latticexyz/mud/compare/v1.39.0...v1.40.0) (2023-03-03)

### Features

- add StoreMetadata table for table name and field names to Store and World ([#428](https://github.com/latticexyz/mud/issues/428)) ([ae39ace](https://github.com/latticexyz/mud/commit/ae39acec12e3263bcb9a8cd4f0a4db5d98b96333))
- v2 - add store, world and schema-type, cli table code generation ([#422](https://github.com/latticexyz/mud/issues/422)) ([cb731e0](https://github.com/latticexyz/mud/commit/cb731e0937e614bb316e6bc824813799559956c8))

### BREAKING CHANGES

- This commit removes the deprecated `mud deploy` CLI command. Use `mud deploy-contracts` instead.
