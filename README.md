# xk6-sftp

A plugin for the k6 load testing tool that adds support for SFTP (Secure File Transfer Protocol) operations. This extension allows you to perform SFTP actions such as uploading, downloading, and managing files on an SFTP server as part of your load testing scripts. It is useful for testing the performance and reliability of systems that rely on SFTP for file transfers.

## Install

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download [xk6](https://github.com/grafana/xk6):
```bash
go install go.k6.io/xk6/cmd/xk6@latest
```

2. [Build](https://github.com/grafana/xk6#command-usage) the k6 binary:
```bash
xk6 build --with github.com/InditexTech/xk6-sftp@latest
```

### Development

For building a `k6` binary with the plugin from the local code, you can run:

```bash
make build
```

For testing and running the plugin locally, an SFTP server is required. The default target in the Makefile will:

- Run an SFTP server with Docker Compose (make sure you have [Docker](https://docs.docker.com/engine/install/) & [Docker Compose](https://docs.docker.com/compose/install/) installed in your system).
- Download the dependencies.
- Format your code.
- Run the integration tests.
- Run the [example](examples/main.js) script.

```bash
git clone git@github.com:InditexTech/xk6-sftp.git
cd xk6-sftp
make
```

## Usage

This extension provides the following JS methods for interacting with the SFTP server:

```javascript
import xk6sftp from "k6/x/sftp";

// Create a new SFTP client
const client = xk6sftp.newClient("username", "password", "host", 3322);

export default function () {
  // Upload a file, creating target directories if necessary
  let result = client.uploadFile("localPath", "remotePath");
  
  // Download a file
  result = client.downloadFile("remotePath", "localPath");
  
  // Delete a file
  result = client.deleteFile("remotePath")
  
  // All the ...File methods return a JS object such as:
  // {"success":true,"message":"file uploaded successfully"}
}

export function teardown() {
  if (client !== null) {
    // Close the SFTP connection
    client.close();
  }
}
```

See the [examples](./examples) folder for a more detailed usage example.
