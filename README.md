# Glow

Flow development in GO! Rapidly test your contracts in an embedded emulator, testnet, or mainnet.

Influenced by https://github.com/bjartek/go-with-the-flow

## Features

1. Sources flow.json to automatically create accounts and deploy contracts.
2. Embedded in memory emulator makes unit testing fast and easy.
3. Create, sign, and submit transactions in a single line.
4. Use "disposable" accounts to test your contracts on Testnet and Mainnet.

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

### Test

```bash
# run all tests (requires test.sh execute permissions)
$ make test

# run all tests w/out test.sh permissions
$ export GLOW_NETWORK=embedded # default
$ export GLOW_ROOT=`pwd`/example
$ go test ./example/test

# run individual test
$ export GLOW_NETWORK=embedded # default
$ export GLOW_ROOT=`pwd`/example
$ go test ./example/test -run TransferFlow
```

### ENV Vars

In order to support multiple contexts in a single repo,

```bash
# export glow root for folder "example" in current directory (pwd)
$ export GLOW_ROOT=`pwd`/example

# export network as one of the following: embedded, emulator, testnet, mainnet (default: embedded)
$ export GLOW_NETWORK=embedded

# export log to specify verbosity level of client logger output (default: 3)
$ export GLOW_LOG=3
```

---

## Client

The Glow client is a configurable flow CLI wrapper.

### Initialization

```go
    client := NewGlowClient().Start()
```

### Accounts

Source accounts from flow.json

```go
    client := NewGlowClient().Start()

    // get service account
    svc := client.SvcAcct

    // get account by name (network is inferred)
    acct := client.FlowJSON.GetAccount("account")

    // create throw-away account
    throwAway, err := client.CreateDisposableAccount() // creates an acct with a common seedphrase

    // create a secure account
    privKey, err := client.NewPrivateKey(SOME_SEED_PHRASE) // create a new crypto private key
    secureAcct, err := client.CreateAccount(privKey)
```

---

## Cadence

Glow has built in amenities to make development in cadence a bit simpler.

### Imports

Contract imports are replaced at runtime. Glow supports two import strategies:

i.e.

1. import NonFungibleToken from 0xNonFungibleToken
2. import NonFungibleToken from "./NonFungibleToken.cdc"

### Txs and Scripts

Transaction and Script objects can be created easily with a client:

```go
    client := NewGlowClient().Start()
    svc := client.SvcAcct
    proposer := client.FlowJSON.GetAccount("proposer")

    // bytes
    tx := client.NewTx(SOME_TX_BYTES, proposer)

    // from string
    tx = client.NewTxFromString(SOME_TX_STRING, proposer)

    // from file
    tx = client.NewTxFromFile("./transaction/account_setup_royalty.cdc", proposer)

    // add args
    tx = tx.Args(
        cadence.Path{
            Domain:     "storage",
            Identifier: "flowTokenVault",
        },
    )

    // add authorizer
    tx, err := tx.AddAuthorizer(svc)

    // sign
    signedTx, err := tx.Sign()

    // send
    res, err := signedTx.Send()

    // sign and send
    res, err = tx.SignAndSend()

    // one liner
    res, err = client.NewTx(SOME_TX_BYTES, proposer, cadence.String("TEST")).SignAndSend()

    // same thing for scripts...
    sc := client.NewSc(SOME_SC_BYTES)
    sc = client.NewScFromString(SOME_SC_STRING)
    sc = client.NewScFromFile("./script/nft_borrow.cdc")

    // exec
    res, err = sc.Exec()

    // one liner
    res, err = sc.NewSc(SOME_SC_BYTES, cadence.String("TEST")).Exec()
```

## Caveats

Rather than throwing an error, the client will always panic when it discovers
missing configuration such as transactions, scripts, contracts, flow.json, accounts, etc...
