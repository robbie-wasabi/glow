
// this script verifies some data with a signature
pub fun main(pubKey: String, signature: String, dataToVerify: String): Bool {
    let pk = PublicKey(
        publicKey: pubKey.decodeHex(),
        signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
    )

    return pk.verify(
        signature: signature.decodeHex(),
        signedData: dataToVerify.utf8,
        domainSeparationTag: "FLOW-V0.0-user",
        hashAlgorithm: HashAlgorithm.SHA3_256
    )
}