module github.com/larsartmann/go-composable-business-types/money

go 1.26.3

require (
	github.com/bojanz/currency v1.4.4
	github.com/larsartmann/go-composable-business-types/locale v0.0.0-00010101000000-000000000000
)

require (
	github.com/cockroachdb/apd/v3 v3.2.1 // indirect
	github.com/larsartmann/go-composable-business-types v0.4.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

replace (
	github.com/larsartmann/go-composable-business-types => ../
	github.com/larsartmann/go-composable-business-types/locale => ../locale
)
