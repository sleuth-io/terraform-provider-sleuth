# Developing the Provider

Preferably use [Devbox](https://www.jetpack.io/devbox/docs/installing_devbox/) for dependencies and development. If you wish to install dependencies manually, see [Requirements](../../README.md).

## Getting dependencies

```sh
devbox install
devbox shell
```

## Running provider locally
You have 2 options to run the provider locally:

1. Use development overrides (preferred):
    1. copy `dev.tfrc.example` to `dev.tfrc`
    1. update `dev.tfrc` with `echo $GOBIN`
        - if `GOBIN` is empty, try `echo $GOPATH` and append it with `/bin/`. e.g.: `echo $GOPATH` is `/go`, update `dev.tfrc` to `"sleuth.io/core/sleuth" = "/go/bin/"`
    1. run `export TF_CLI_CONFIG_FILE=dev.tfrc`
    1. check that overrides are successfully applied with `terraform apply` -> you should see `Provider development overrides are in effect` in the output
    1. run `make install` to build the provider and install it to `$GOBIN`
    1. Done! You can now use the provider locally with `terraform apply`

2. Manually install the provider:
    1. run `make install_deprecated`
        - this will build the provider and put the provider binary in the `$GOPATH/bin` directory & in correct location for terraform providers
        - *Note:* if you're on macOS, change `OS_ARCH` in `Makefile` to `darwin_amd64` (Intel) or `darwin_arm64` (Apple silicon)
    2. run `terraform init` to initialize the provider and use newly built binary
    3. Done! You can now use the provider locally with `terraform apply`

## Generating documentation
Run `make docs`. This will read files in `/templates` and `examples` folder and generate documentation in `/docs` folder. Note that changes done directly in `/docs` folder will be overwritten.

## Formatting
Run `make fmt`.

## Testing
In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Running against a local instance of Sleuth
To run against a local instance of Sleuth, do the following:
1. Start Sleuth locally so that it is available on http://dev.sleuth.io
2. Copy the `main.tf.example` file as `main.tf` and edit the file to change `api_key`
3. Run terraform via `make dev`. Note this will delete the state each time.

## Debugging
When you have context availabe you should use `tflog.LEVEL|(ctx)` for logging. This will print out the log message with context and level. For example `tflog.DEBUG|(ctx)` will print out `DEBUG: [terraform-provider-sleuth] (ctx) message`.

*Note:* Levels for provider can be adjusted using `TF_LOG_PROVIDER=LEVEL`.

If context is not available you can use `fmt.Print("here")` with combination of `TF_LOG=DEBUG` env variable when running `terraform plan` or `terraform apply`.

## Adding Go Dependencies
This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.
