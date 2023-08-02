/*
ecs-snapshot computes and saves ECS state from a chain via snapshots, such that the client can
perform an initial sync to the ECS world state without having to process all ECS state changes
in form of events.

It exposes a [pkg.grpc.StartSnapshotServer] for getting the snapshot data and a gRPC-web HTTP server
wrapper to allow clients which are web applications to request snapshots with a POST request.

By default, ecs-snapshot attempts to connect to a local development chain and also by default,
ecs-snapshot indexes and snapshots all ECS state it finds via emitted events, which can be
restricted by specifying contract addresses.

Usage:

	ecs-snapshot [flags]

The flags are:

	    -ws-url
	        Websocket URL for chain to index and snapshot.
	    -port
	        Port to expose the gRPC server.
		-worldAddresses
			Comma-separated list of contract addresses to limit the indexing to. If left blank, index
			everything, otherwise, use this list as a filter.
		-block
			Block to start taking snapshots from. Defaults to 0.
		-snapshot-block-interval
			Block number interval for how often to take regular snapshots.
		-initial-sync-block-batch-size
			Number of blocks to fetch data for when performing an initial sync.
		-initial-sync-block-batch-sync-timeout
			Time in milliseconds to wait between calls to fetch batched log data when performing an initial sync.
		-initial-sync-snapshot-interval
			Block number interval for how often to take intermediary snapshots when performing an initial sync.
		-default-snapshot-chunk-percentage
			Default percentage for RPCs that request a snapshot in chunks. Default to 10, i.e. 10 percent chunks.
		-metrics-port
			Prometheus metrics http handler port. Defaults to port 6060.
*/
package main

import (
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"latticexyz/mud/packages/services/pkg/eth"
	"latticexyz/mud/packages/services/pkg/grpc"
	"latticexyz/mud/packages/services/pkg/logger"
	"latticexyz/mud/packages/services/pkg/snapshot"
	"latticexyz/mud/packages/services/pkg/utils"
)

func Createkey(str string) string {
	t := strings.ToUpper(strings.ReplaceAll(str, "-", "_"))
	return t
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	var (
		wsUrl                            = os.Getenv(Createkey("ws-url"))                                //, "Websocket Url")
		port                             = os.Getenv(Createkey("snapshot-port"))                         //, 50061, "gRPC Server Port")
		worldAddresses                   = os.Getenv(Createkey("world-addresses"))                       //, "0xD84379CEae14AA33C123Af12424A37803F885889", "List of world addresses to index ECS state for. Defaults to empty string which will listen for all world events from all addresses")
		block                            = os.Getenv(Createkey("block"))                                 //, 0, "Block to start taking snapshots from. Defaults to 0")
		snapshotBlockInterval            = os.Getenv(Createkey("snapshot-block-interval"))               //, 100, "Block number interval for how often to take regular snapshots")
		initialSyncBlockBatchSize        = os.Getenv(Createkey("initial-sync-block-batch-size"))         //, 10, "Number of blocks to fetch data for when performing an initial sync")
		initialSyncBlockBatchSyncTimeout = os.Getenv(Createkey("initial-sync-block-batch-sync-timeout")) //, 100, "Time in milliseconds to wait between calls to fetch batched log data when performing an initial sync")
		initialSyncSnapshotInterval      = os.Getenv(Createkey("initial-sync-snapshot-interval"))        //, 5000, "Block number interval for how often to take intermediary snapshots when performing an initial sync")
		defaultSnapshotChunkPercentage   = os.Getenv(Createkey("default-snapshot-chunk-percentage"))     //, 10, "Default percentage for RPCs that request a snapshot in chunks. Default to 10, i.e. 10 percent chunks")
		metricsPort                      = os.Getenv(Createkey("snapshot-metrics-port"))                 // 6060)                    //, "Prometheus metrics http handler port. Defaults to port 6060")
	)

	// Setup logging.
	logger.InitLogger()
	logger := logger.GetLogger()
	defer logger.Sync()

	snapshotBlockIntervalInt64, _ := strconv.Atoi(snapshotBlockInterval)
	initialSyncBlockBatchSizeInt64, _ := strconv.Atoi(initialSyncBlockBatchSize)
	initialSyncBlockBatchSyncTimeoutInt, _ := strconv.Atoi(initialSyncBlockBatchSyncTimeout)
	initialSyncSnapshotIntervalInt, _ := strconv.Atoi(initialSyncSnapshotInterval)
	defaultSnapshotChunkPercentageInt, _ := strconv.Atoi(defaultSnapshotChunkPercentage)
	portInt, _ := strconv.Atoi(port)
	metricsPortInt, _ := strconv.Atoi(metricsPort)
	blockInt, _ := strconv.Atoi(block)

	// Build a config.
	config := &snapshot.SnapshotServerConfig{
		SnapshotBlockInterval:            int64(snapshotBlockIntervalInt64),
		InitialSyncBlockBatchSize:        int64(initialSyncBlockBatchSizeInt64),
		InitialSyncBlockBatchSyncTimeout: time.Duration(initialSyncBlockBatchSyncTimeoutInt) * time.Millisecond,
		InitialSyncSnapshotInterval:      int64(initialSyncSnapshotIntervalInt),
		DefaultSnapshotChunkPercentage:   defaultSnapshotChunkPercentageInt,
	}

	// Parse world addresses to listen to.
	worlds := utils.SplitAddressList(worldAddresses, ",")
	if len(worlds) == 0 {
		logger.Info("listening for events from all world addresses")
	} else {
		logger.Info("listening for events from specific addresses", zap.String("worldAddresses", worldAddresses))
	}

	// Get an instance of ethereum client.
	ethclient := eth.GetEthereumClient(wsUrl, logger)

	// Start gRPC server.
	go grpc.StartSnapshotServer(portInt, metricsPortInt, config, logger)

	// 1. Prepare for service to run.
	utils.EnsureDir(snapshot.SnapshotDir)

	// 2. Kick off the service to catch up on state up to the current block number.
	fromBlock := big.NewInt(int64(blockInt))
	toBlock := eth.GetCurrentBlockHead(ethclient)

	initialState := snapshot.Sync(ethclient, fromBlock, toBlock, worlds, config)

	// 3. Kick off the service to start syncing with new block heads from the current one.
	snapshot.Start(initialState, ethclient, fromBlock, worlds, config, logger)
}
