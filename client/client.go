package client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/gateway"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/afero"

	// . "github.com/rrossilli/glow/util"

	. "github.com/rrossilli/glow/model"
	. "github.com/rrossilli/glow/tmp"
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

type GlowClientBuilder struct {
	InMemory     bool
	IniAccts     bool
	DepContracts bool
	GasLim       uint64
	HashAlgo     crypto.HashAlgorithm
	SigAlgo      crypto.SignatureAlgorithm

	// network vars
	LogLvl  int
	Root    string
	Network string
}

func (b *GlowClientBuilder) InitAccounts(l bool) *GlowClientBuilder {
	b.IniAccts = l
	return b
}

func (b *GlowClientBuilder) DeployContracts(l bool) *GlowClientBuilder {
	b.DepContracts = l
	return b
}

// todo: specify which contracts to deploy

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

func NewGlowClientBuilder(network, root string, logLvl int) *GlowClientBuilder {
	if network == "" {
		network = NETWORK_EMBEDDED
	}

	inMemory := false
	depContracts := false
	initAccounts := false
	hashAlgo := crypto.StringToHashAlgorithm(SHA3_256)
	sigAlgo := crypto.StringToSignatureAlgorithm(ECDSA_P256)

	if network == NETWORK_EMBEDDED {
		network = NETWORK_EMULATOR
		inMemory = true
		depContracts = true
		initAccounts = true
	}

	return &GlowClientBuilder{
		Network:      network,
		InMemory:     inMemory,
		IniAccts:     initAccounts,
		DepContracts: depContracts,
		LogLvl:       logLvl,
		GasLim:       9999,
		Root:         root,
		HashAlgo:     hashAlgo,
		SigAlgo:      sigAlgo,
	}
}

type GlowClient struct {
	network  string
	root     string
	FlowJSON FlowJSON
	Logger   output.Logger
	Services *services.Services
	State    *flowkit.State
	HashAlgo crypto.HashAlgorithm
	SigAlgo  crypto.SignatureAlgorithm
	SvcAcct  Account
	gasLimit uint64
}

func (c *GlowClient) GetNetwork() string {
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

// source flow.json
func parseFlowJSON(file string) (flowJSON FlowJSON) {
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byteValue, &flowJSON)

	return flowJSON
}

// Start Client
func (b *GlowClientBuilder) Start() *GlowClient {
	logger := output.NewStdoutLogger(b.LogLvl)
	loader := &afero.Afero{Fs: afero.NewOsFs()}
	fJSONPath := fmt.Sprintf("%s/flow.json", b.Root) // assumes that flow.json is at root
	// fmt.Printf("fJSONPath: %v\n", fJSONPath)
	state, err := flowkit.Load([]string{fJSONPath}, loader)
	if err != nil {
		logger.Error(fmt.Sprintf("\nFlowkit was unable to load project configuration at path: %s", b.Root))
		panic(err)
	}
	flowJSON := parseFlowJSON(fJSONPath)

	logger.Info("\n==================================")
	logger.Info("STARTING GLOW CLIENT:\n")
	logger.Info(fmt.Sprintf("NETWORK: %v", b.Network))
	logger.Info(fmt.Sprintf("IN MEMORY: %v", b.InMemory))
	logger.Info(fmt.Sprintf("ROOT: %v", b.Root))

	var service *services.Services
	if b.InMemory {
		svcAcct, _ := state.EmulatorServiceAccount()
		gw := gateway.NewEmulatorGateway(svcAcct)
		service = services.NewServices(gw, state, logger)
	} else {
		network, err := state.Networks().ByName(b.Network)
		if err != nil {
			panic(err)
		}
		host := network.Host
		gw, err := gateway.NewGrpcGateway(host)
		if err != nil {
			panic(err)
		}
		service = services.NewServices(gw, state, logger)
	}

	svcAcct := flowJSON.GetSvcAcct(b.Network)

	wrappedClient := GlowClient{
		network:  b.Network,
		root:     b.Root,
		FlowJSON: flowJSON,
		Logger:   logger,
		Services: service,
		State:    state,
		HashAlgo: b.HashAlgo,
		SigAlgo:  b.SigAlgo,
		gasLimit: b.GasLim,
		SvcAcct:  svcAcct,
	}

	if b.IniAccts {
		wrappedClient.initAccounts()
	}

	if b.DepContracts {
		wrappedClient.deployContracts()
	}

	logger.Info("==================================")

	return &wrappedClient
}

// Submit transactions to initialize accounts sourced from flow.json
func (c *GlowClient) initAccounts() {
	c.Logger.Info("\nCREATING ACCOUNTS:\n")
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

		c.Logger.Info(fmt.Sprintf("%s CREATED", acct.Address))
	}
}

// Submit transactions to deploy contracts to existing accounts sourced from flow.json
func (c *GlowClient) deployContracts() {
	c.Logger.Info("\nDEPLOYING CONTRACTS:\n")
	acctNames := c.FlowJSON.AccountNamesSorted(c.network) // sorted list of account names
	for _, a := range acctNames {
		d := c.FlowJSON.GetAccountDeployment(c.network, a)
		for _, d := range d {
			// get acct and deploy contract
			acct := c.FlowJSON.GetAccount(a)
			contract := c.GetContractCdc(d)
			txRes, err := c.NewTxFromString(
				TX_CONTRACT_DEPLOY,
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
			c.Logger.Info(fmt.Sprintf("%s DEPLOYED", d))
		}
	}
}
