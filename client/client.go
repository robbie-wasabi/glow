package client

import (
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

	. "github.com/rrossilli/glow/util"

	. "github.com/rrossilli/glow/model"
)

const (
	ECDSA_P256 = "ECDSA_P256"
	SHA3_256   = "SHA3_256"

	ENV_EMBEDDED = "embedded"
	ENV_EMULATOR = "emulator"
	ENV_TESTNET  = "testnet"
	ENV_MAINNET  = "mainnet"

	DEFAULT_LOG_LEVEL            = 3
	DEFAULT_EMULATOR_SVC_ACCOUNT = "emulator-svc"
)

type GlowClientBuilder struct {
	InMemory     bool
	InitAccts    bool
	DepContracts bool
	GasLim       uint64
	HashAlgo     crypto.HashAlgorithm
	SigAlgo      crypto.SignatureAlgorithm

	// env vars
	LogLvl       int
	FlowJSONPath string
	Env          string
}

func (b *GlowClientBuilder) InitAccounts(l bool) *GlowClientBuilder {
	b.InitAccts = l
	return b
}

func (b *GlowClientBuilder) DeployContracts(l bool) *GlowClientBuilder {
	b.DepContracts = l
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

func NewGlowClientBuilder(env, flowJSONPath string, logLvl int) *GlowClientBuilder {
	if env == "" {
		env = ENV_EMBEDDED
	}

	inMemory := false
	depContracts := false
	initAccounts := false
	hashAlgo := crypto.StringToHashAlgorithm(SHA3_256)
	sigAlgo := crypto.StringToSignatureAlgorithm(ECDSA_P256)

	if env == ENV_EMBEDDED {
		env = ENV_EMULATOR
		inMemory = true
		depContracts = true
		initAccounts = true
	}

	return &GlowClientBuilder{
		Env:          env,
		InMemory:     inMemory,
		InitAccts:    initAccounts,
		DepContracts: depContracts,
		LogLvl:       logLvl,
		GasLim:       9999,
		FlowJSONPath: flowJSONPath,
		HashAlgo:     hashAlgo,
		SigAlgo:      sigAlgo,
	}
}

type GlowClient struct {
	env      string
	FlowJSON FlowJSON
	Logger   output.Logger
	Services *services.Services
	State    *flowkit.State
	HashAlgo crypto.HashAlgorithm
	SigAlgo  crypto.SignatureAlgorithm
	gasLimit uint64
}

func NewGlowClient() *GlowClientBuilder {
	env := os.Getenv("ENV")

	flowJSONPath := os.Getenv("FJSON")
	absPath := ROOT + flowJSONPath

	log := os.Getenv("LOG")
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

	c := NewGlowClientBuilder(env, absPath, logLvl)

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
	fmt.Printf("b.FlowJSONPath: %v\n", b.FlowJSONPath)
	state, err := flowkit.Load([]string{b.FlowJSONPath}, loader)
	if err != nil {
		logger.Error("\nFlowkit was unable to load project configuration: make sure to 'export FJSON=<flow_json_path>'")
		panic(err)
	}
	flowJSON := parseFlowJSON(b.FlowJSONPath)

	logger.Info("\n==================================")
	logger.Info("STARTING CLIENT")
	logger.Info(fmt.Sprintf("env: %v", b.Env))
	logger.Info(fmt.Sprintf("in memory: %v", b.InMemory))
	logger.Info(fmt.Sprintf("flow.json: %v", b.FlowJSONPath))
	logger.Info("==================================\n")

	var service *services.Services
	if b.InMemory {
		svcAcct, _ := state.EmulatorServiceAccount()
		gw := gateway.NewEmulatorGateway(svcAcct)
		service = services.NewServices(gw, state, logger)
	} else {
		network, err := state.Networks().ByName(b.Env)
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

	wrappedClient := GlowClient{
		env:      b.Env,
		FlowJSON: flowJSON,
		Logger:   logger,
		Services: service,
		State:    state,
		HashAlgo: b.HashAlgo,
		SigAlgo:  b.SigAlgo,
		gasLimit: b.GasLim,
	}

	if b.InitAccts {
		wrappedClient.initAccounts()
	}

	if b.DepContracts {
		wrappedClient.deployContracts()
	}

	return &wrappedClient
}

// Submit transactions to initialize accounts sourced from flow.json
func (c GlowClient) initAccounts() {
	svcAcct := c.GetSvcAcct()
	accounts := c.AccountsSorted()
	for i, a := range accounts {
		// skip svc account
		if i == 0 {
			continue
		}

		acct, err := c.CreateAccount(
			a.CryptoPrivateKey(),
			svcAcct,
		)
		if err != nil {
			panic(err)
		}

		c.Logger.Info(fmt.Sprintf("%s CREATED", acct.Address))
	}
}

// Submit transactions to deploy contracts to existing accounts sourced from flow.json
func (c GlowClient) deployContracts() {
	acctNames := c.AccountNamesSorted() // sorted list of account names
	for _, a := range acctNames {
		d := c.GetAccountDeployment(a)
		for _, d := range d {
			acct := c.GetAccount(a)
			contract := c.FlowJSON.GetContract(d)
			txRes, err := c.DeployContract(
				cadence.String(d),
				RemoveFirstChar(contract.Source), // must remove leading "." due to project structure
				acct,
			)
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
