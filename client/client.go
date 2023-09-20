package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/config"
	"github.com/onflow/flow-cli/flowkit/gateway"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-emulator/emulator"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/afero"

	"github.com/rs/zerolog"

	"github.com/rrossilli/glow/model"
	"github.com/rrossilli/glow/tmp"
)

const (
	ECDSA_P256 = "ECDSA_P256"
	SHA3_256   = "SHA3_256"

	NETWORK_EMBEDDED = "embedded"
	NETWORK_EMULATOR = "emulator"
	NETWORK_TESTNET  = "testnet"
	NETWORK_MAINNET  = "mainnet"

	DEFAULT_LOG_LEVEL            = 3
	DEFAULT_EMULATOR_SVC_ACCOUNT = "emulator-svc"
)

// Responsible for building instances of GlowClient.
type GlowClientBuilder struct {
	InMemory, ShouldCreateAccounts, ShouldDeployContracts bool
	GasLim                                                uint64
	HashAlgo                                              crypto.HashAlgorithm
	SigAlgo                                               crypto.SignatureAlgorithm
	LogLvl                                                int
	Root, NetworkName                                     string
}

// Toggles the account creation feature.
func (b *GlowClientBuilder) CreateAccounts(l bool) *GlowClientBuilder {
	b.ShouldCreateAccounts = l
	return b
}

// Toggles the contract deployment feature.
func (b *GlowClientBuilder) DeployContracts(l bool) *GlowClientBuilder {
	b.ShouldDeployContracts = l
	return b
}

func (b *GlowClientBuilder) HashAlgorithm(algo string) *GlowClientBuilder {
	b.HashAlgo = crypto.StringToHashAlgorithm(algo)
	return b
}

func (b *GlowClientBuilder) SigAlgorithm(algo string) *GlowClientBuilder {
	b.SigAlgo = crypto.StringToSignatureAlgorithm(algo)
	return b
}

func (b *GlowClientBuilder) LogLevel(level int) *GlowClientBuilder {
	b.LogLvl = level
	return b
}

func (b *GlowClientBuilder) GasLimit(limit uint64) *GlowClientBuilder {
	b.GasLim = limit
	return b
}

// Initializes a new GlowClientBuilder with default settings.
func NewGlowClientBuilder(network, root string, logLvl int) *GlowClientBuilder {
	if network == "" {
		network = NETWORK_EMBEDDED
	}

	inMemory := false
	shouldDeployContracts := false
	shouldCreateAccounts := false
	hashAlgo := crypto.StringToHashAlgorithm(SHA3_256)
	sigAlgo := crypto.StringToSignatureAlgorithm(ECDSA_P256)

	if network == NETWORK_EMBEDDED {
		network = NETWORK_EMULATOR
		inMemory = true
		shouldDeployContracts = true
		shouldCreateAccounts = true
	}

	return &GlowClientBuilder{
		NetworkName:           network,
		InMemory:              inMemory,
		ShouldCreateAccounts:  shouldCreateAccounts,
		ShouldDeployContracts: shouldDeployContracts,
		LogLvl:                logLvl,
		GasLim:                9999,
		Root:                  root,
		HashAlgo:              hashAlgo,
		SigAlgo:               sigAlgo,
	}
}

// Encapsulates Flow network interactions.
type GlowClient struct {
	network  config.Network
	root     string
	FlowJSON model.FlowJSON
	Logger   output.Logger
	FlowKit  *flowkit.Flowkit
	State    *flowkit.State
	HashAlgo crypto.HashAlgorithm
	SigAlgo  crypto.SignatureAlgorithm
	SvcAcct  model.Account
	gasLimit uint64
}

// Returns the network configuration.
func (c *GlowClient) GetNetwork() config.Network {
	return c.network
}

func NewGlowClient() *GlowClientBuilder {
	network := os.Getenv("GLOW_NETWORK")
	root := os.Getenv("GLOW_ROOT")
	log := os.Getenv("GLOW_LOG")

	var logLvl int
	if log != "" {
		lvl, err := strconv.Atoi(log)
		if err != nil {
			panic(err)
		}
		logLvl = lvl
	} else {
		logLvl = DEFAULT_LOG_LEVEL
	}

	c := NewGlowClientBuilder(network, root, logLvl)

	return c
}

func parseFlowJSON(file string) (flowJSON model.FlowJSON) {
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byteValue, &flowJSON)

	return flowJSON
}

// Initializes the GlowClient with the configurations set in the builder.
func (b *GlowClientBuilder) Start() *GlowClient {
	logger := output.NewStdoutLogger(b.LogLvl)
	loader := &afero.Afero{Fs: afero.NewOsFs()}

	fJSONPath := fmt.Sprintf("%s/flow.json", b.Root) // assumes that flow.json is at root
	state, err := flowkit.Load([]string{fJSONPath}, loader)
	if err != nil {
		// logger.Error(fmt.Sprintf("\nFlowkit was unable to load project configuration at path: %s", b.Root))
		panic(err)
	}

	flowJSON := parseFlowJSON(fJSONPath)

	network, err := state.Networks().ByName(b.NetworkName)
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("\nGlow Client Starting: Network=%v, InMemory=%v, Root=%v", b.NetworkName, b.InMemory, b.Root))

	var kit *flowkit.Flowkit
	var gw gateway.Gateway
	if b.InMemory {
		var memlog bytes.Buffer
		writer := io.Writer(&memlog)
		emulatorLogger := zerolog.New(writer).Level(zerolog.DebugLevel)
		emulatorOpts := []emulator.Option{
			emulator.WithLogger(emulatorLogger),
		}

		svcAcct, err := state.EmulatorServiceAccount()
		if err != nil {
			panic(err)
		}

		pk, err := svcAcct.Key.PrivateKey()
		if err != nil {
			panic(err)
		}

		emulatorKey := &gateway.EmulatorKey{
			PublicKey: (*pk).PublicKey(),
			SigAlgo:   b.SigAlgo,
			HashAlgo:  b.HashAlgo,
		}

		gw = gateway.NewEmulatorGatewayWithOpts(emulatorKey, gateway.WithLogger(&emulatorLogger), gateway.WithEmulatorOptions(emulatorOpts...))
	} else {
		gw, err = gateway.NewGrpcGateway(*network)
		if err != nil {
			panic(err)
		}

	}

	kit = flowkit.NewFlowkit(state, *network, gw, logger)
	svcAcct := flowJSON.GetSvcAcct(b.NetworkName)
	wrappedClient := GlowClient{
		network:  *network,
		root:     b.Root,
		FlowJSON: flowJSON,
		Logger:   logger,
		FlowKit:  kit,
		State:    state,
		HashAlgo: b.HashAlgo,
		SigAlgo:  b.SigAlgo,
		gasLimit: b.GasLim,
		SvcAcct:  svcAcct,
	}

	if b.ShouldCreateAccounts {
		wrappedClient.createAccounts()
	}

	if b.ShouldDeployContracts {
		wrappedClient.deployContracts()
	}

	return &wrappedClient
}

// Initializes accounts on the Flow network
func (c *GlowClient) createAccounts() {
	c.Logger.Info("Creating Accounts:")

	accounts := c.FlowJSON.AccountsSorted()
	for i, a := range accounts {
		// skip svc account
		if i == 0 {
			continue
		}

		acct, err := c.CreateAccount(
			a.CryptoPrivateKey(),
		)
		if err != nil {
			panic(err)
		}

		c.Logger.Info(fmt.Sprintf("Account=%s Created", acct.Address))
	}
}

// Deploys smart contracts to accounts on the Flow network
func (c *GlowClient) deployContracts() {
	c.Logger.Info("Deploy Contracts:")

	acctNames := c.FlowJSON.AccountNamesSorted(c.network.Name) // sorted list of account names
	for _, a := range acctNames {
		d := c.FlowJSON.GetAccountDeployment(c.network.Name, a)
		for _, d := range d {
			// get acct and deploy contract
			acct := c.FlowJSON.GetAccount(a)
			contract := c.GetContractCdc(d)
			txRes, err := c.NewTxFromString(
				tmp.TX_CONTRACT_DEPLOY,
				acct,
				contract.NameAsCadenceString(),
				cadence.String(hex.EncodeToString(contract.CdcBytes())),
			).SignAndSend()
			if err != nil {
				panic(err)
			}
			if txRes.Error != nil {
				panic(txRes.Error)
			}
			c.Logger.Info(fmt.Sprintf("Contract=%s Deployed", d))
		}
	}
}
