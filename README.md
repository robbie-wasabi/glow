# Glow

_now sponsored by [FLOW](https://flow.com/)!_

Glow is a comprehensive Golang library designed to streamline the development, testing, and deployment of Flow smart contracts, transactions, and scripts across various Flow networks. By leveraging a configuration-driven architecture, Glow simplifies account creation, contract deployment, and state management, allowing developers to interact with the Flow blockchain through a consistent, programmatic interface.

Originally conceived as a personal passion project, Glow seeks to alleviate common pain points encountered during Flow smart contract development. With inspiration drawn from [bjartek/overflow](https://github.com/bjartek/overflow), Glow aspires to be a valuable resource for testing contracts and workflows, providing a well-structured and easily extensible toolkit.

## Key Features

1. **Configuration-Driven Initialization:**  
   Glow automatically reads and parses the `flow.json` configuration file, using it to instantiate accounts and deploy contracts upon client initialization.

2. **Versatile Testing Environments:**  
   Write and execute unit tests against the Flow mainnet, testnet, or even an embedded runtime emulator. This flexibility allows developers to validate contract logic and transaction behavior in multiple environments without complex setup.

3. **Disposable Accounts:**  
   Easily create short-lived accounts intended for one-off tests. These “throwaway” accounts enable quick iterative testing without the overhead of cleaning up state or reconfiguring environments.

4. **Concise Transaction and Script Execution:**  
   Compose, sign, and submit transactions—or run scripts—in a single line of code. This streamlined workflow reduces boilerplate and encourages a more intuitive development experience.

---

## Getting Started

Below are the foundational steps required to set up the development environment and begin leveraging Glow’s capabilities:

### Installation and Setup

```bash
# Install the Flow CLI, an essential component for interacting with the Flow blockchain.
$ brew install flow-cli

# Optionally install the Flow VSCode extension for enhanced Cadence syntax highlighting and editing.
$ flow cadence install-vscode-extension

# Install all necessary Go modules and dependencies specified by the project.
$ make init

# (Optional) Grant execute permissions to 'test.sh' to facilitate running tests via 'make test'.
$ chmod 777 test.sh
```

### Environment Variables

Glow respects environment variables to determine project context and verbosity:

```bash
# Specify the project’s root directory. For example, if you have an "example" folder in the current directory:
$ export GLOW_ROOT=`pwd`/example

# Define the target Flow network. Supported values: embedded, emulator, testnet, mainnet.
# The default is 'embedded' which leverages the in-memory emulator.
$ export GLOW_NETWORK=embedded

# Control the verbosity level of the Glow client’s logger.
# The default value of 3 provides a moderate amount of detail.
$ export GLOW_LOG=3
```

**Tip:** You may find it useful to write small shell scripts (e.g., `./test.sh`) for running tests with a predetermined set of environment variables, ensuring consistency and convenience.

### Running the Example Tests

Below are some approaches to execute the provided example tests:

```bash
# Run all tests using the Makefile target, assuming test.sh has execute permissions.
$ make test

# Alternatively, run tests directly without using test.sh:
$ export GLOW_NETWORK=embedded # This is the default network if none is specified.
$ export GLOW_ROOT=`pwd`/example
$ go test ./example/test

# Run a single test by name for more targeted debugging:
$ export GLOW_NETWORK=embedded
$ export GLOW_ROOT=`pwd`/example
$ go test ./example/test -run TransferFlow
```

---

## Glow Client Overview

The Glow client encapsulates all core functionalities—reading configuration data, initializing accounts, deploying contracts, and running transactions or scripts. It acts as a user-friendly abstraction on top of lower-level Flow interactions.

### Initialization & Configuration

On startup, the Glow client ingests `flow.json` to:

- Identify and initialize accounts.
- Deploy configured contracts to their respective accounts.
- Expose network-aware references to resources for runtime usage.

For more information on `flow.json` configuration, refer to [Flow CLI Configuration Documentation](https://developers.flow.com/tools/flow-cli/configuration).

**Example:**

```go
// Initialize the Glow client.
// This step reads flow.json, configures accounts, and deploys contracts.
client := NewGlowClient().Start()

// Retrieve a contract definition from the config.
contract := client.FlowJSON.GetContract("MyContract")

// Access accounts defined in flow.json.
acct := client.FlowJSON.GetAccount("SomeAccount")

// Directly fetch the service account for the configured network.
svcAcct := client.FlowJSON.GetSvcAcct("emulator")

// Fetch a map of accounts associated with a specific network.
accounts := client.FlowJSON.GetAccounts("testnet")

// Retrieve deployment details, either globally or specific to one account.
deployment := client.FlowJSON.GetDeployment("mainnet")
deployments := client.FlowJSON.GetAccountDeployment("testnet", "TestAccount")
```

### Working with Keys

Glow provides convenient utilities for creating and managing cryptographic keys.

**Examples:**

```go
client := NewGlowClient().Start()

// Derive a private key from a seed phrase.
cryptoPrivateKey, err := client.NewPrivateKey("my seed phrase ...")

// Derive a private key directly from a hex string (with or without '0x' prefix).
cryptoPrivKeyFromString, err := client.NewPrivateKeyFromString("YOUR_PRIVATE_KEY_STRING")

// Similarly, derive a public key from a hex string.
cryptoPubKeyFromString, err := client.NewPublicKeyFromString("YOUR_PUBLIC_KEY_STRING")
```

### Accounts

Accounts in `flow.json` are keyed by their network, enabling network-specific configurations. For example:

```json
{
  "accounts": {
    "emulator-svc": {
      "address": "0xf8d6e0586b0a20c7",
      "key": "3a63ae4f8fffacd89d1b574d87fe448a0f848da7d0a45c04b60744b1c3905a14"
    },
    "emulator-account": {
      "address": "01cf0e2f2f715450",
      "key": "4a4f7a1d07b441135489823f1bcdc27ba607c1916b3b182a2b7ee91cf11eb5f6"
    },
    "testnet-svc": {
      "address": "0xbc450f7d561b7bc1",
      "key": "4d17ed74bef04b66e9c5ea299de7831a7815239188d45afe9e69a6b54dd966fd"
    }
  }
}
```

**Examples of Account Manipulation:**

```go
client := NewGlowClient().Start()

// Retrieve the primary service account for the current network.
svc := client.SvcAcct // Shorthand for service account retrieval.

// Get a named account. The network is inferred, so "emulator-account" can be referenced by "account".
acct := client.FlowJSON.GetAccount("account")

// Create a temporary "disposable" account for ephemeral tests.
tmpAcct, err := client.CreateDisposableAccount()

// Create a secure account from a newly generated private key.
privKey, err := client.NewPrivateKey("some seed phrase")
secureAcct, err := client.CreateAccount(privKey)

// Access helpful properties and methods:
address := acct.Address
cadenceAddress := acct.CadenceAddress()
privateKey := acct.PrivKey
publicKey := acct.CryptoPrivateKey().PublicKey()
```

---

### Cadence Integration

Glow streamlines Cadence development through automatic import resolution and other conveniences.

#### Imports

Contracts can be imported in two primary ways:

1. Relative imports, referencing `.cdc` files directly:
   ```cadence
   import NonFungibleToken from "./NonFungibleToken.cdc"
   ```

2. Address imports, referencing contract addresses that are resolved at runtime:
   ```cadence
   import NonFungibleToken from 0xNonFungibleToken
   ```

Both approaches are supported. The first option is often preferable for local development as it integrates smoothly with the VSCode Flow extension’s syntax highlighting and code navigation features.

### Transactions and Scripts

Glow provides multiple methods to create and submit transactions or execute scripts. Whether you supply the code inline, load it from a file, or define it as raw bytes, the process is uniform and concise.

**Transaction Examples:**

```go
client := NewGlowClient().Start()
proposer := client.FlowJSON.GetAccount("proposer")

// Construct a transaction from a byte slice.
tx := client.NewTx(TX_BYTES, proposer)

// Load a transaction from a string or file.
tx = client.NewTxFromString(TX_STRING, proposer)
tx = client.NewTxFromFile("./transactions/my_transaction.cdc", proposer)

// Add arguments to your transaction.
tx = tx.Args(cadence.Path{
  Domain:     "storage",
  Identifier: "flowTokenVault",
})

// Add an authorizer account.
tx, err := tx.AddAuthorizer(client.SvcAcct)

// Sign the transaction with the default key (index 0).
signedTx, err := tx.Sign()

// Alternatively, sign with a specific key index.
signedTx, err = tx.SignWithKeyAtIndex(1)

// Finally, submit the transaction to the network.
res, err := signedTx.Send()

// Or simply sign and send in one step.
res, err = tx.SignAndSend()

// One-liner for convenience:
res, err = client.NewTx(TX_BYTES, proposer, cadence.String("TEST_ARG")).SignAndSend()
```

**Script Examples:**

```go
client := NewGlowClient().Start()

// Construct a script similarly to transactions.
sc := client.NewSc(SC_BYTES)
sc = client.NewScFromString(SC_STRING)
sc = client.NewScFromFile("./scripts/query_nft.cdc")

// Execute the script and capture the result.
res, err := sc.Exec()

// One-liner for scripts:
res, err = client.NewSc(SC_BYTES, cadence.String("TEST_ARG")).Exec()
```

### Signing Arbitrary Data

For cryptographic operations beyond transactions and scripts, Glow supports signing arbitrary data:

```go
import (
  // ...
  "github.com/onflow/flow-go-sdk/crypto"
)

client := NewGlowClient().Start()
signer := client.FlowJSON.GetAccount("account")

// Sign arbitrary data using a chosen hashing algorithm.
signedData, err := signer.SignMessage([]byte("some_data"), crypto.SHA3_256)
```

---

## Important Caveats and Notes

Below are some crucial points to keep in mind when using Glow.

### Configuration Errors

If Glow encounters an invalid configuration (e.g., missing required contracts, malformed `flow.json`, or absent accounts), it will panic rather than return an error. While this behavior may evolve in the future, it currently serves as a fail-fast mechanism to ensure noticeable runtime errors.

### Logging Within Cadence

Within Cadence code, `log()` statements do not produce visible output in the terminal. To observe runtime values, consider using `panic()` calls, as these are surfaced in the output and can serve as an ad-hoc debugging mechanism.

**Example:**

```cadence
// This will not produce any output:
log("This message won't be visible.")

// This will produce a terminal output:
panic("This message will be displayed upon transaction/script failure.")
```

### Emulator-Specific Account Creation

When utilizing the Flow emulator, account addresses are predetermined rather than dynamically generated. As a result, each address must be explicitly defined in `flow.json` in the order they are expected to appear. For example, if the first three emulator accounts are:

```text
0xf8d6e0586b0a20c7
0x01cf0e2f2f715450
0x179b6b1cb6755e31
```

Then `flow.json` should list them accordingly:

```json
"emulator-svc": {
  "address": "0xf8d6e0586b0a20c7",
  "key": "<service_account_key>"
},
"emulator-acct-1": {
  "address": "0x01cf0e2f2f715450",
  "key": "<some_key>"
},
"emulator-acct-2": {
  "address": "0x179b6b1cb6755e31",
  "key": "<some_key>"
}
```

This ensures proper alignment between Glow’s expectations and the underlying emulator environment.

---

By combining a configuration-driven approach, flexible testing environments, and user-friendly abstractions over transactions, scripts, and accounts, Glow provides a more streamlined, productive workflow for Flow blockchain developers. It aims to simplify common tasks, allowing you to focus on the logic and integrity of your smart contracts rather than wrestling with boilerplate setup.
