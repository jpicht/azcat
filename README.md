# simple azure blob tools
This is a collection of very simple tools to manipulate blobs in azure blob storage. Most functions of the azure blob store aren't available through these tools.

## Motivation
I created these tools out of the pain and frustration that comes with trying to work with the tools currently available (looking at you, [azure cli](https://github.com/Azure/azure-cli) and [azcopy](https://github.com/Azure/azure-storage-azcopy)).

For most usecases `azure cli` is absolutely over-powered and the official docker container weighs in at a whopping 1,13GiB (at time of writing).

`azcopy` is a go tool, that has a lot of good functionality. The interface is a bit clunky (anyone else who needed some time to wrap their head around `--from-to`? Still not sure what half of the possible values do), but the deal breaker for me are the [intentional](https://github.com/Azure/azure-storage-azcopy/issues/186) lack of SharedKey and SAS authentication and the choice to use the **kernel keyring**, which is **disabled in docker** by default, or the **gnome-session-keyring**, also not available in docker, to store the session token. Without any option to use *anything* else.

## Usage
List files with prefix
```bash
azls https://<account>.blob.core.windows.net/<container>/<prefix>
```

Output blob contents to STDOUT
```bash
azcat https://<account>.blob.core.windows.net/<container>/<blob>
```

Write blob contents from STDIN
```bash
mongodump --archive | azput https://<account>.blob.core.windows.net/<container>/<blob>
```

Remove blob
```bash
azrm https://<account>.blob.core.windows.net/<container>/<blob>
```

## Authentication
Currently authentication will always be automatically derived from the environment. The environment variables are conveniently named just like the `azure cli` expects.

Supported methods (are tried in this order):
1. connection string `AZURE_STORAGE_CONNECTION_STRING`
1. shared key via `AZURE_STORAGE_ACCOUNT_NAME` and `AZURE_STORAGE_ACCOUNT_KEY`
1. OAuth via `AZURE_TENANT_ID`, `AZURE_CLIENT_ID` and
  1. `AZURE_CLIENT_SECRET` or
  1. `AZURE_CLIENT_CERTIFICATE_PATH` or
  2. `AZURE_USERNAME` and `AZURE_PASSWORD`
1. metadata service `169.254.169.254`

## Installation
There are no releases (yet), so the easiest way is installing using go install:

```
go install github.com/jpicht/azcat/cmd/multi/...
```

This will install the four standalone tools `azcat`, `azls`, `azput` and `azrm` into
your `$GOPATH/bin/` directory.

There is a fifth tool called `azblob` which can be used in scenarios where more than one function is needed, but saving on size is necessary. It can either be used directly, but the command line options are a bit clunky and subject to change, or four symlinks (
`azcat`, `azls`, `azput` and `azrm`) can be created pointing to it, and it behaves (nearly) exactly like these tools.
