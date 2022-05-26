package tmp

const (
	SC_FLOW_BALANCE = `
        import FungibleToken from 0xFungibleToken
        import FlowToken from 0xFlowToken
        
        // get an account's flow balance 
        pub fun main(account: Address): UFix64 {
            let vaultRef = getAccount(account)
                .getCapability(/public/flowTokenBalance)
                .borrow<&FlowToken.Vault{FungibleToken.Balance}>()
                ?? panic("Could not borrow Balance reference to the Vault")
        
            let balance = vaultRef.balance
        
            return balance
        }
    `

	SC_CONTRACT_CHECK = `
        pub fun main(
            address: Address, 
            name: String
        ): Bool {
            let contract = getAccount(address).contracts.get(name: name)
                ?? panic("could not retrieve ".concat(name).concat(" from account ").concat(address.toString()))
            return true
        }
    `

	// todo...
)
