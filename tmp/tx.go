package tmp

const (
	// Deploys a contract to an account.
	TX_CONTRACT_DEPLOY = `
	transaction(name: String, code: String) {
		prepare(signer: AuthAccount) {
			signer.contracts.add(name: name, code: code.decodeHex())
		}
	}`

	// Removes a previously deployed contract from an account.
	TX_CONTRACT_REMOVE = `
	transaction(name: String) {
		prepare(signer: AuthAccount) {
			signer.contracts.remove(name: name)
		}
	}`

	// Updates an existing contract on an account.
	TX_CONTRACT_UPDATE = `
	transaction(name: String, code: String) {
		prepare(signer: AuthAccount) {
			signer.contracts.update__experimental(name: name, code: code.decodeHex())
		}
	}`

	// Creates a new Flow account with a specified public key.
	TX_CREATE_ACCOUNT = `
	transaction(publicKey: String) {
		prepare(signer: AuthAccount) {
			let account = AuthAccount(payer: signer)
			let key = PublicKey(
				publicKey: publicKey.decodeHex(),
				signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
			)
			account.keys.add(
				publicKey: key,
				hashAlgorithm: HashAlgorithm.SHA3_256,
				weight: 1000.0
			)
		}
	}`

	// Transfers Flow tokens from one account to another.
	TX_FLOW_TRANSFER = `
	import FungibleToken from 0xFungibleToken
	import FlowToken from 0xFlowToken

	transaction(amount: UFix64, recipient: Address) {
		let sentVault: @FungibleToken.Vault
		prepare(signer: AuthAccount) {
			let vault = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
				?? panic("Could not borrow sender vault reference")

			self.sentVault <- vault.withdraw(amount: amount)
		}
		execute {
			let receiver = getAccount(recipient)
				.getCapability(/public/flowTokenReceiver)
				.borrow<&{FungibleToken.Receiver}>()
				?? panic("Could not borrow recipient vault reference")

			receiver.deposit(from: <-self.sentVault)
		}
	}`
)
