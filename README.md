# Embedded Reflash Tool

Small utility built primarily to replace AMI MegaRAC with OpenBMC in the field.
To build for Aspeed type SoCs:

Pre Go 1.22

```bash
GOARCH=arm GOARM=5 go build -ldflags="-s -w"
```

Go 1.22 or later

```bash
GOARCH=arm GOARM=6,softfloat go build -ldflags="-s -w" # AST2500
GOARCH=arm GOARM=7,softfloat go build -ldflags="-s -w" # AST2600
```
