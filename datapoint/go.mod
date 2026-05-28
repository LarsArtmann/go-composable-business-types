module github.com/larsartmann/go-composable-business-types/datapoint

go 1.26.3

require (
	github.com/larsartmann/go-branded-id v0.3.0
	github.com/larsartmann/go-composable-business-types v0.4.0
	github.com/larsartmann/go-composable-business-types/nanoid v0.0.0-00010101000000-000000000000
)

require (
	github.com/sixafter/aes-ctr-drbg v1.19.1 // indirect
	github.com/sixafter/nanoid v1.64.2 // indirect
	github.com/sixafter/prng-chacha v1.16.2 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
)

replace (
	github.com/larsartmann/go-composable-business-types => ../
	github.com/larsartmann/go-composable-business-types/nanoid => ../nanoid
)
