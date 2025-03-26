module test

go 1.21

toolchain go1.23.5

require (
	gopkg.in/yaml.v3 v3.0.1
	soundwave-go v0.0.0
)

replace soundwave-go => ../soundwave-go
