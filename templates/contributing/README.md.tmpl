# Developing the Provider

Either:

- [Devbox](https://www.jetpack.io/devbox/docs/installing_devbox/) for dependencies and development
- install dependencies manually, see [Requirements](../../README.md).
- GitHub's Codespaces which comes with `go` preinstalled out of the box
    - to run in Codespaces open this repo on GitHub, find green `<> Code` dropdown button and select `Codespaces` tab

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
Run `make format`.

## Testing

Tests are run as GitHub actions. Tests are defined in the folder `./internal/`.

The tests literally create projects and code deployments and impact sources etc. on sleuth staging. So, if the tests don't pass, there is a good chance that the problem isn't on the side of this code but on the side of Sleuth.

Since the tests contact Sleuth staging directly, they also depend on various objects to exist there, ie. there must be precisely 1 PagerDuty integration, the API key for the org must be correct, ... .

The tests are also run as 1 group, they are always run all together, but they are also **run for various terraform versions**, that is why it appears that there are multiple tests:

![img.png](tests.png)

Originally, the idea was to run `make testacc` to run the tests locally, but this does not work for me at all, maybe it will work for you.

What did work for me is to call (but I had to have Sleuth running locally on http://dev.sleuth.io/):

```shell
TF_ACC=1 SLEUTH_BASEURL="http://dev.sleuth.io" SLEUTH_API_KEY="f4d4c4******" go test -v -cover ./internal/...
```

You can also point the url to staging, to get the true results.

After the tests are run, they **do delete the object they've created.** This is great because they clean up after themselves, but also not so great, because you can't inspect what truly went wrong since the objects don't exist anymore.


**Note:** Tests create real resources, and often cost money to run.
**Note:** Sometimes you might need to call `go build -v .`.

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
