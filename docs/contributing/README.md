# Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](../../README.md)).
The easiest way to do this is to use dosbox and run `dosbox shell`.

To compile the provider, run `make install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make docs`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

To run against a local instance of Sleuth, do the following:
1. Start Sleuth locally so that it is available on http://dev.sleuth.io 
2. Copy the `main.tf.example` file as `main.tf` and edit the file to change `api_key`
3. Run terraform via `main dev`. Note this will delete the state each time.

# Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

# Updating documentation

`/docs` folder is entirely generated using `make docs` command.
In order to write any documentation you need write it in `/templates` folder and the contents of that will get copied to the docs folder.
