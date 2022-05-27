# env
export GLOW_ROOT=`pwd`/example
export GLOW_NETWORK=embedded
export GLOW_LOG=3

# tests
go test ./example/test -run TestTransferFlow -v 3
go test ./example/test -run TestMintNFT -v 3