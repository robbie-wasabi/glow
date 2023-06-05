// update a contract in an account

transaction(name: String, code: String) {
    prepare(signer: AuthAccount) {
        signer.contracts.update__experimental(name: name, code: code.decodeHex())
    }
}