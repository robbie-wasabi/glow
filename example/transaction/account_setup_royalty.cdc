
// This transaction is a template for a transaction
// to create a new link in their account to be used for receiving royalties
// This transaction can be used for any fungible token, which is specified by the `vaultPath` argument
// 
// If the account wants to receive royalties in FLOW, they'll use `/storage/flowTokenVault`
// If they want to receive it in USDC, they would use FiatToken.VaultStoragePath
// and so on. 
// The path used for the public link is a new path that in the future, is expected to receive
// and generic token, which could be forwarded to the appropriate vault
import FungibleToken from "../contract/FungibleToken.cdc"
import MetadataViews from "../contract/MetadataViews.cdc"
import FlowToken from "../contract/FlowToken.cdc"

transaction(vaultPath: StoragePath) {

    prepare(signer: AuthAccount) {

        // Create new ft vault at vaultPath
        if signer.borrow<&FungibleToken.Vault>(from: vaultPath) == nil {
            // Create a new flowToken Vault and put it in storage
            signer.save(<-FlowToken.createEmptyVault(), to: vaultPath)

            // Create a public capability to the Vault that only exposes
            // the deposit function through the Receiver interface
            signer.link<&{FungibleToken.Receiver}>(
                /public/flowTokenReceiver,
                target: vaultPath
            )

            // Create a public capability to the Vault that only exposes
            // the balance field through the Balance interface
            signer.link<&{FungibleToken.Balance}>(
                /public/flowTokenBalance,
                target: vaultPath
            )
        }

        // Return early if the account doesn't have a FungibleToken Vault
        if signer.borrow<&FungibleToken.Vault>(from: vaultPath) == nil {
            panic("A vault for the specified fungible token path does not exist")
        }

        // Create a public capability to the Vault that only exposes
        // the deposit function through the Receiver interface
        let capability = signer.link<&{FungibleToken.Receiver, FungibleToken.Balance}>(
            MetadataViews.getRoyaltyReceiverPublicPath(),
            target: vaultPath
        )!

        // Make sure the capability is valid
        if !capability.check() { panic("Beneficiary capability is not valid!") }
    }
}