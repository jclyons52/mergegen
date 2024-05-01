# Mergegen

This repo is a test of using code generation to build merge functions for structs.

to run the utility build it 

```
go build generator.go
```

Then run the binary:

```
./generator -src=./config.go -type=Config -output=config_merge.go
```