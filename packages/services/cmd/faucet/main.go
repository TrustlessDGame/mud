/*
faucet acts as a configurable faucet with in-service integrations for MUD transactions. Optionally, requests for drip can be gated behind Twitter signature verification.

By default, faucet attempts to connect to a local development chain.

Usage:

    faucet [flags]

The flags are:

    -ws-url
        Websocket URL for sending optional integrated MUD transactions to set Component values on, for example, successful drip.
    -port
        Port to expose the gRPC server.
	-dev
		Flag to run the faucet in dev mode, where verification is not required. Default to false.
	-faucet-private-key
		Private key to use for faucet.
	-drip-amount
		Drip amount in ETH. Default to 0.01 ETH
	-drip-frequency
		Drip frequency per account in minutes. Default to 60 minutes.
	-drip-limit
		Drip limit in ETH per drip frequency interval. Default to 1 ETH
	-twitter
		Flag to run the faucet in Twitter mode, where to receive a drip you have to tweet a signature. Default to false.
	-num-latest-tweets
		Number of latest tweets to check per user when verifying drip tweet. Default to 5.
	-name-system-address
		Address of NameSystem to set an address/username mapping when verifying drip tweet. Not specified by default.
	-metrics-port
		Prometheus metrics http handler port. Defaults to port 6060.
*/

package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"latticexyz/mud/packages/services/pkg/eth"
	"latticexyz/mud/packages/services/pkg/faucet"
	"latticexyz/mud/packages/services/pkg/grpc"
	"latticexyz/mud/packages/services/pkg/logger"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	"golang.org/x/oauth2/clientcredentials"
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
		// General flags.
		wsUrl = os.Getenv(Createkey("ws-url"))      //, "ws://localhost:8545", "Websocket Url")
		port  = os.Getenv(Createkey("faucet-port")) //, 50081, "gRPC Server Port")
		// Dev mode.
		devMode = os.Getenv(Createkey("dev")) //,  //true,)  "Flag to run the faucet in dev mode, where verification is not required. Default to false")
		// Faucet configuration flags.
		faucetPrivateKey = os.Getenv(Createkey("faucet-private-key")) //, "0x58ac750167fecf1f4daafa31fc6ef0afa54d20d42b498a3c0b9358457f9d97c1", "Private key to use for faucet")
		// Drip configuration flags.
		dripAmount    = os.Getenv(Createkey("drip-amount"))    //, 0.01, "Drip amount in ETH. Default to 0.01 ETH")
		dripFrequency = os.Getenv(Createkey("drip-frequency")) //, 60, "Drip frequency per account in minutes. Default to 60 minutes")
		dripLimit     = os.Getenv(Createkey("drip-limit"))     //, 1, "Drip limit in ETH per drip frequency interval. Default to 1 ETH")
		// Flags for using twitter to verify drip requests.
		twitterMode       = os.Getenv(Createkey("twitter"))             //, false, "Flag to run the faucet in Twitter mode, where to receive a drip you have to tweet a signature. Default to false")
		numLatestTweets   = os.Getenv(Createkey("num-latest-tweets"))   //, 5, "Number of latest tweets to check per user when verifying drip tweet. Default to 5")
		nameSystemAddress = os.Getenv(Createkey("name-system-address")) //, "", "Address of NameSystem to set an address/username mapping when verifying drip tweet. Not specified by default")
		metricsPort       = os.Getenv(Createkey("metrics-port"))        //, 6061, "Prometheus metrics http handler port. Defaults to port 6060")
	)

	// Setup logging.
	logger.InitLogger()
	logger := logger.GetLogger()
	defer logger.Sync()

	dripAmountF, _ := strconv.ParseFloat(dripAmount, 10)
	dripFrequencyF, _ := strconv.ParseFloat(dripFrequency, 10)
	dripLimitF, _ := strconv.ParseFloat(dripLimit, 10)
	devModeBool, _ := strconv.ParseBool(devMode)
	twitterModeBool, _ := strconv.ParseBool(twitterMode)
	numLatestTweetsInt, _ := strconv.Atoi(numLatestTweets)
	dripFrequencyInt, _ := strconv.Atoi(dripFrequency)
	portInt, _ := strconv.Atoi(port)
	metricsPortInt, _ := strconv.Atoi(metricsPort)

	// Create a drip config.
	dripConfig := &faucet.DripConfig{
		DripAmount:               float64(dripAmountF),
		DripFrequency:            float64(dripFrequencyF),
		DripLimit:                float64(dripLimitF),
		DevMode:                  devModeBool,
		TwitterMode:              twitterModeBool,
		NumLatestTweetsForVerify: numLatestTweetsInt,
		NameSystemAddress:        nameSystemAddress,
	}
	logger.Info("using a drip configuration",
		zap.Float64("amount", dripConfig.DripAmount),
		zap.Float64("frequency", dripConfig.DripFrequency),
		zap.Float64("limit", dripConfig.DripLimit),
		zap.Bool("dev", dripConfig.DevMode),
		zap.Bool("twitter", dripConfig.TwitterMode),
	)

	// Ensure that a twitter <-> address store is setup.
	faucet.SetupStore()

	// Oauth2 configures a client that uses app credentials to keep a fresh token.
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	// Get a connection to Twitter API with a client.
	twitterClient := twitter.NewClient(config.Client(context.Background()))

	// Get an instance of ethereum client.
	ethClient := eth.GetEthereumClient(wsUrl, logger)

	// Create a private key ECDSA object.
	privateKey, err := crypto.HexToECDSA(faucetPrivateKey)
	if err != nil {
		logger.Fatal("error creating ECDSA object from private key string", zap.String("privateKey", faucetPrivateKey))
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		logger.Fatal("error casting public key to ECDSA")
	}

	// Kick off a worked that will reset the faucet limit at the specified interval.
	// Note: the duration here matches whatever time units are used in 'drip-frequency'.
	go faucet.ReplenishFaucetWorker(time.NewTicker(time.Duration(dripFrequencyInt)*time.Minute), make(chan struct{}))

	// Start the faucet gRPC server.
	grpc.StartFaucetServer(portInt, metricsPortInt, twitterClient, ethClient, privateKey, publicKey, dripConfig, logger)
}
