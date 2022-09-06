# Glow

Flow development in GO! Rapidly test your contracts in an embedded emulator, testnet, or mainnet.

Influenced by https://github.com/bjartek/overflow

## Features

1. Sources flow.json to automatically create accounts and deploy contracts on glow client initialization.
2. Write unit tests to run against the flow testnet/mainnet or the embedded runtime emulator.
3. Create "disposable" accounts to execute tests without any state cleanup.
4. Create, sign, and send transactions in a single line of code.

---

## Setup

```bash
# install flow cli
$ brew install flow-cli

# install flow vscode extension (optional)
$ flow cadence install-vscode-extension

# install dependencies
$ make init

# give execute permissions to test.sh to run "make test" cmd (optional)
$ chmod 777 test.sh
```

### ENV Vars

```bash
# export glow root for folder "example" in current directory (pwd)
$ export GLOW_ROOT=`pwd`/example

# export network as one of the following: embedded, emulator, testnet, mainnet (default: embedded)
$ export GLOW_NETWORK=embedded

# export log to specify verbosity level of client logger output (default: 3)
$ export GLOW_LOG=3
```

Quick Tip: bash scripts like ./test.sh are useful in order to run a group of tests with a specified configuration.

### Run Example Tests

```bash
# run all tests (requires test.sh execute permissions, see "## Setup")
$ make test

# - or you can run the tests without the test.sh script
$ export GLOW_NETWORK=embedded # default
$ export GLOW_ROOT=`pwd`/example
$ go test ./example/test

# - or individually
$ export GLOW_NETWORK=embedded 
$ export GLOW_ROOT=`pwd`/example
$ go test ./example/test -run TransferFlow
```

---

## Client

The Glow client is a configurable flow CLI wrapper that makes smart contract development simple and efficient.

### Initialization

Glow sources the flow.json to create all relevant accounts, deploy all pertinent contracts, and make its contents available at runtime.

flow.json configuration info here: https://developers.flow.com/tools/flow-cli/configuration

```go
  // init glow client
  client := NewGlowClient().Start()

  // get contract
  contract := client.FlowJSON.GetContract(CONTRACT_NAME)

  // get account(s)
  account := client.FlowJSON.GetAccount(ACCOUNT_NAME)
  account := client.FlowJSON.GetSvcAcct(NETWORK_NAME)
  accounts := client.FlowJSON.GetAccounts(NETWORK_NAME)

  // get deployment
  deployment := client.FlowJSON.GetDeployment(NETWORK_NAME)
  deployments := client.FlowJSON.GetAccountDeployment(NETWORK_NAME, ACCOUNT_NAME)
```

### Keys

```go
  client := NewGlowClient().Start() 

  // create private key from seed phrase
  cryptoPrivateKey, err := client.NewPrivateKey(SEED_PHRASE)

  // create private key from string (0x prefix is optional)
  cryptoPrivKeyFromString, err := client.NewPrivateKeyFromString(PRIV_KEY_STRING)

  // create public key from string (0x prefix is optional)
  cryptoPubKeyFromString, err := client.NewPublicKeyFromString(PRIV_KEY_STRING)
```

### Accounts

Accounts in the flow.json are prefixed with the associated network:

```json
{
  "accounts": {
    "emulator-svc": {
      "address": "0xf8d6e0586b0a20c7",
      "key": "3a63ae4f8fffacd89d1b574d87fe448a0f848da7d0a45c04b60744b1c3905a14"
    },
    "emulator-account": {
      "address": "01cf0e2f2f715450", // 0x prefix is not necessary
      "key": "4a4f7a1d07b441135489823f1bcdc27ba607c1916b3b182a2b7ee91cf11eb5f6"
    },
    "testnet-svc": {
      "address": "0xbc450f7d561b7bc1",
      "key": "4d17ed74bef04b66e9c5ea299de7831a7815239188d45afe9e69a6b54dd966fd"
    },
  }
}
```

Working with accounts:

```go
    client := NewGlowClient().Start()

    // get service account
    svc = client.FlowJSON.GetAccount("svc")
    // - or with shorthand
    svc := client.SvcAcct 

    // get account by name.
    // network is inferred so "emulator-account" should be written as "account"
    acct := client.FlowJSON.GetAccount("account")

    // create throw-away account
    throwAwayAcct, err := client.CreateDisposableAccount() // creates an acct with a common seedphrase

    // create a secure account
    privKey, err := client.NewPrivateKey(SEED_PHRASE) // create a new crypto private key
    secureAcct, err := client.CreateAccount(privKey)

    // helpful functions
    address := acct.Address
    cadenceAddress := acct.CadenceAddress()
    privateKey := acct.PrivKey
    publicKey := acct.CryptoPrivateKey().PublicKey()
```

---

### Cadence

Glow has built in amenities to make development in cadence a bit simpler.

#### Imports

Contract imports are replaced at runtime. Glow supports two import strategies:

i.e.

```cadence
    // preferable for syntax highlighting using the vscode extension 
    // (https://developers.flow.com/tools/vscode-extension)
    import NonFungibleToken from "./NonFungibleToken.cdc"

    // this also works
    import NonFungibleToken from 0xNonFungibleToken 
```

### Txs and Scripts

Transaction and Script objects can be created easily with a client:

```go
    client := NewGlowClient().Start()
    svc := client.SvcAcct
    proposer := client.FlowJSON.GetAccount("proposer")

    // bytes
    tx := client.NewTx(TX_BYTES, proposer)

    // from string
    tx = client.NewTxFromString(TX_STRING, proposer)

    // from file
    tx = client.NewTxFromFile(PATH_TO_TX, proposer)

    // add args
    tx = tx.Args(
        cadence.Path{
            Domain:     "storage",
            Identifier: "flowTokenVault",
        },
    )

    // add authorizer
    tx, err := tx.AddAuthorizer(svc)

    // sign with default key (key at index 0)
    signedTx, err := tx.Sign()

    // sign with a key at specified index
    signedTx, err := tx.SignWithKeyAtIndex(KEY_INDEX)

    // send
    res, err := signedTx.Send()

    // sign with default key and send
    res, err = tx.SignAndSend()

    // tx one liner
    res, err = client.NewTx(TX_BYTES, proposer, cadence.String("TEST")).SignAndSend()

    // same thing for scripts...
    sc := client.NewSc(SC_BYTES)
    sc = client.NewScFromString(SC_STRING)
    sc = client.NewScFromFile("./script/nft_borrow.cdc")

    // exec
    res, err = sc.Exec()

    // script one liner
    res, err = sc.NewSc(SC_BYTES, cadence.String("TEST")).Exec()
```

### Signing Arbitrary Data

```go
    import (
        // ...
        "github.com/onflow/flow-go-sdk/crypto"
    )

    client := NewGlowClient().Start()
    signer := client.FlowJSON.GetAccount("account")

    // sign data with hash algo
    signedData, err := signer.SignMessage([]byte(DATA_STRING), HASH_ALGO)

    // example: sign data with SHA3_256 hash algo
    signedDataSHA3256, err := signer.SignMessage([]byte("some_message"), crypto.SHA3_256)
```

## Caveats

Rather than throwing an error, the client will always panic when it discovers
missing configuration such as transactions, scripts, contracts, flow.json, accounts, etc...

The "log()" function in cadence does not print any output in the terminal...
This is obviously not ideal but use "panic()" in scripts and txns to print desired log outputs.

```js
    // no output
    log(message)

    // output
    if true {
        panic(message)
    }
```
