# @latticexyz/store-sync

## 2.0.0-next.1

### Major Changes

- [#1198](https://github.com/latticexyz/mud/pull/1198) [`e86fbc12`](https://github.com/latticexyz/mud/commit/e86fbc1260f698c6a7b6a92c901fefd186c329ff) Thanks [@holic](https://github.com/holic)! - Adds store indexer service package with utils to query the indexer service.

  You can run the indexer locally by checking out the MUD monorepo, installing/building everything, and running `pnpm start:local` from `packages/store-indexer`.

  To query the indexer in the client, you can create a tRPC client with a URL pointing to the indexer service and call the available tRPC methods:

  ```ts
  import { createIndexerClient } from "@latticexyz/store-sync/trpc-indexer";

  const indexer = createIndexerClient({ url: indexerUrl });
  const result = await indexer.findAll.query({
    chainId: publicClient.chain.id,
    address,
  });
  ```

  If you're using `syncToRecs`, you can just pass in the `indexerUrl` option as a shortcut to the above:

  ```ts
  import { syncToRecs } from "@latticexyz/store-sync/recs";

  syncToRecs({
    ...
    indexerUrl: "https://your.indexer.service",
  });
  ```

### Patch Changes

- [#1197](https://github.com/latticexyz/mud/pull/1197) [`9e5baf4f`](https://github.com/latticexyz/mud/commit/9e5baf4fff0c60615b8f2b4645fb11cb78cb0bd8) Thanks [@holic](https://github.com/holic)! - Add RECS sync strategy and corresponding utils

  ```ts
  import { createPublicClient, http } from 'viem';
  import { syncToRecs } from '@latticexyz/store-sync';
  import storeConfig from 'contracts/mud.config';
  import { defineContractComponents } from './defineContractComponents';

  const publicClient = createPublicClient({
    chain,
    transport: http(),
    pollingInterval: 1000,
  });

  const { components, singletonEntity, latestBlock$, blockStorageOperations$, waitForTransaction } = await syncToRecs({
    world,
    config: storeConfig,
    address: '0x...',
    publicClient,
    components: defineContractComponents(...),
  });
  ```

- Updated dependencies [[`c963b46c`](https://github.com/latticexyz/mud/commit/c963b46c7eaceebc652930936643365b8c48a021), [`b02f9d0e`](https://github.com/latticexyz/mud/commit/b02f9d0e43089e5f9b46d817ea2032ce0a1b0b07), [`afdba793`](https://github.com/latticexyz/mud/commit/afdba793fd84abf17eef5ef59dd56fabe353c8bd)]:
  - @latticexyz/store@2.0.0-next.1
  - @latticexyz/schema-type@2.0.0-next.1
  - @latticexyz/recs@2.0.0-next.1
  - @latticexyz/block-logs-stream@2.0.0-next.1
  - @latticexyz/common@2.0.0-next.1
  - @latticexyz/protocol-parser@2.0.0-next.1

## 2.0.0-next.0

### Minor Changes

- [#1075](https://github.com/latticexyz/mud/pull/1075) [`904fd7d4`](https://github.com/latticexyz/mud/commit/904fd7d4ee06a86e481e3e02fd5744224376d0c9) Thanks [@holic](https://github.com/holic)! - Add store sync package

- [#1176](https://github.com/latticexyz/mud/pull/1176) [`eeb15cc0`](https://github.com/latticexyz/mud/commit/eeb15cc06fcbe80c37ba3926d9387f6bd5947234) Thanks [@holic](https://github.com/holic)! - - Replace `blockEventsToStorage` with `blockLogsToStorage` that exposes a `storeOperations` callback to perform database writes from store operations. This helps encapsulates database adapters into a single wrapper/instance of `blockLogsToStorage` and allows for wrapping a block of store operations in a database transaction.

  - Add `toBlock` option to `groupLogsByBlockNumber` and remove `blockHash` from results. This helps track the last block number for a given set of logs when used in the context of RxJS streams.

- [#1185](https://github.com/latticexyz/mud/pull/1185) [`69a96f10`](https://github.com/latticexyz/mud/commit/69a96f109065ae2564a340208d5f9a0be3616747) Thanks [@holic](https://github.com/holic)! - `blockLogsToStorage(sqliteStorage(...))` converts block logs to SQLite operations. You can use it like:

  ```ts
  import { drizzle } from "drizzle-orm/better-sqlite3";
  import Database from "better-sqlite3";
  import { BaseSQLiteDatabase } from "drizzle-orm/sqlite-core";
  import { createPublicClient } from "viem";
  import { blockLogsToStorage } from "@latticexyz/store-sync";
  import { sqliteStorage } from "@latticexyz/store-sync/sqlite";

  const database = drizzle(new Database('store.db')) as any as BaseSQLiteDatabase<"sync", void>;
  const publicClient = createPublicClient({ ... });

  blockLogs$
    .pipe(
      concatMap(blockLogsToStorage(sqliteStorage({ database, publicClient }))),
      tap(({ blockNumber, operations }) => {
        console.log("stored", operations.length, "operations for block", blockNumber);
      })
    )
    .subscribe();
  ```

### Patch Changes

- [#1179](https://github.com/latticexyz/mud/pull/1179) [`53522998`](https://github.com/latticexyz/mud/commit/535229984565539e6168042150b45fe0f9b48b0f) Thanks [@holic](https://github.com/holic)! - - bump to viem 1.3.0 and abitype 0.9.3
  - move `@wagmi/chains` imports to `viem/chains`
  - refine a few types
- Updated dependencies [[`904fd7d4`](https://github.com/latticexyz/mud/commit/904fd7d4ee06a86e481e3e02fd5744224376d0c9), [`b98e5180`](https://github.com/latticexyz/mud/commit/b98e51808aaa29f922ac215cf666cf6049e692d6), [`4bb7e8cb`](https://github.com/latticexyz/mud/commit/4bb7e8cbf0da45c85b70532dc73791e0e2e1d78c), [`ca50fef8`](https://github.com/latticexyz/mud/commit/ca50fef8108422a121d03571fb4679060bd4891a), [`eeb15cc0`](https://github.com/latticexyz/mud/commit/eeb15cc06fcbe80c37ba3926d9387f6bd5947234), [`72b80697`](https://github.com/latticexyz/mud/commit/72b806979db6eb2880772193898351d657b94f75), [`8d51a034`](https://github.com/latticexyz/mud/commit/8d51a03486bc20006d8cc982f798dfdfe16f169f), [`48909d15`](https://github.com/latticexyz/mud/commit/48909d151b3dfceab128c120bc6bb77de53c456b), [`66cc35a8`](https://github.com/latticexyz/mud/commit/66cc35a8ccb21c50a1882d6c741dd045acd8bc11), [`f03531d9`](https://github.com/latticexyz/mud/commit/f03531d97c999954a626ef63bc5bbae51a7b90f3), [`a7b30c79`](https://github.com/latticexyz/mud/commit/a7b30c79bcc78530d2d01858de46a0fb87954fda), [`53522998`](https://github.com/latticexyz/mud/commit/535229984565539e6168042150b45fe0f9b48b0f), [`0c4f9fea`](https://github.com/latticexyz/mud/commit/0c4f9fea9e38ba122316cdd52c3d158c62f8cfee)]:
  - @latticexyz/block-logs-stream@2.0.0-next.0
  - @latticexyz/protocol-parser@2.0.0-next.0
  - @latticexyz/store@2.0.0-next.0
  - @latticexyz/common@2.0.0-next.0
  - @latticexyz/schema-type@2.0.0-next.0
  - @latticexyz/store-cache@2.0.0-next.0
