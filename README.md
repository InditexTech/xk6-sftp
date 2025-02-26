# xk6-sftp

The `xk6-sftp` extension is a plugin for the k6 load testing tool that adds support for SFTP (Secure File Transfer Protocol) operations. This extension allows you to perform SFTP actions such as uploading, downloading, and managing files on an SFTP server as part of your load testing scripts. It is useful for testing the performance and reliability of systems that rely on SFTP for file transfers.

## Install

### Pre-built binaries 

``` sh
make run
```

### Build from source

``` sh
make build
```

## Examples

See [examples](./examples/) folder.

## Extension API

- `newClient`: creates a new SFTP client.
- `downloadFile`: to download a file.
- `uploadFile`: to upload a file creating target directories if necessary.
- `close`: closes the SFTP connection.

