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

Test connection
```bash
azping https://<account>.blob.core.windows.net/<container>/<blob>
```

## Authentication
Currently authentication will always be automatically derived from the environment. We're using the default azure credential discovery flow implemented in [azidentity.NewDefaultAzureCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#NewDefaultAzureCredential).

## Installation
There are no releases (yet), so the easiest way is installing using go install:

```
go install github.com/jpicht/azcat/cmd/multi/...
```

This will install the five standalone tools `azcat`, `azls`, `azping`, `azput` and `azrm` into
your `$GOPATH/bin/` directory.

There is a fifth tool called `azblob` which can be used in scenarios where more than one function is needed, but saving on size is necessary. It can either be used directly, but the command line options are a bit clunky and subject to change, or five symlinks (
`azcat`, `azls`, `azping`, `azput` and `azrm`) can be created pointing to it, and it behaves (nearly) exactly like these tools.

## Contributors

Thank you @voyvodov for updating the credential flow.