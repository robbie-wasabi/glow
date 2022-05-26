package tmp

const (
	TX_CONTRACT_DEPLOY = `
        transaction(name: String, code: String) {
            prepare(signer: AuthAccount) {
                signer.contracts.add(name: name, code: code.decodeHex())
            }
        }
    `

	TX_CONTRACT_REMOVE = `
        transaction(name: String) {
            prepare(signer: AuthAccount) {
                signer.contracts.remove(name: name)
            }
        }
    `

	TX_CONTRACT_UPDATE = `
        transaction(name: String, code: String) {
            prepare(signer: AuthAccount) {
                signer.contracts.update__experimental(name: name, code: code.decodeHex())
            }
        }
    `

	TX_CREATE_ACCOUNT = `
		transaction(publicKey: String) {
			prepare(signer: AuthAccount) {
				let account = AuthAccount(payer: signer)
				let accountKey = PublicKey(
					publicKey: publicKey.decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
				)
				account.keys.add(
					publicKey: accountKey,
					hashAlgorithm: HashAlgorithm.SHA3_256,
					weight: 1000.0
				)
			}
		}
	`

	TX_FLOW_TRANSFER = `
        import FungibleToken from 0xFungibleToken
        import FlowToken from 0xFlowToken
        
        transaction(amount: UFix64, recipient: Address) {
            let sentVault: @FungibleToken.Vault
            prepare(signer: AuthAccount) {
                let vaultRef = signer.borrow<&FlowToken.Vault>(from: /storage/flowTokenVault)
                    ?? panic("failed to borrow reference to sender vault")
            
                self.sentVault <- vaultRef.withdraw(amount: amount)
            }
            
            execute {
                let receiverRef =  getAccount(recipient)
                .getCapability(/public/flowTokenReceiver)
                .borrow<&{FungibleToken.Receiver}>()
                    ?? panic("failed to borrow reference to recipient vault")
            
                receiverRef.deposit(from: <-self.sentVault)
            }
        }
    `

	// todo...
)
