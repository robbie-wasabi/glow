# glow

Flow development in GO! Primarly used to rapidly test and debug smart contracts in controlled local environments.

---

## ENV Vars

```bash
# export glow root for folder "example" in current directory (pwd)
$ export GLOW_ROOT=`pwd`/example

# export network as one of the following: embedded, emulator, testnet, mainnet
$ export GLOW_NETWORK=embedded
```

---

## Caveats

rather than throwing an error, the client will always panic when it discovers
missing configuration such as transactions, scripts, contracts, flow.json, accounts, etc...
