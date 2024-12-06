package tmp

const (
	// SC_FLOW_BALANCE retrieves the Flow token balance of a given account.
	SC_FLOW_BALANCE = `
	import FungibleToken from 0xFungibleToken
	import FlowToken from 0xFlowToken

	pub fun main(account: Address): UFix64 {
		let vaultRef = getAccount(account)
			.getCapability(/public/flowTokenBalance)
			.borrow<&FlowToken.Vault{FungibleToken.Balance}>()
			?? panic("Could not borrow Balance reference")

		return vaultRef.balance
	}
	`

	// SC_CONTRACT_CHECK verifies if a given contract exists in the specified account.
	SC_CONTRACT_CHECK = `
	pub fun main(address: Address, name: String): Bool {
		let contract = getAccount(address).contracts.get(name: name)
			?? panic("Could not retrieve ".concat(name).concat(" from account ").concat(address.toString()))
		return true
	}
	`
)
