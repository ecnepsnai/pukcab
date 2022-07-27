# Pukcab module: SCP

Module name: `scp`

This module enables you tocopy any file from a remote host using SCP.

# Requirements

- SSH Host Key authentication must be used, password-based authentication is not supported
- The `scp` binary must be present on the host running pukcab
- The private key cannot be password protected
- Only a single file can be downloaded

# Configuration

|Key|Type|Description|
|---|----|-----------|
|`host_address`|string|The host address or target server to connect to. Do not include a port number.|
|`port`|number|Optionally specify a port number. If omitted or set to 0, 22 is used.|
|`username`|string|The username to identify as to the remote host.|
|`private_key`|string|The private key in PEM format, including headers. Replace all newlines with `\n`.|
|`host_public_key`|string|The SSH public key of this host. Must include the algorithm. Example: `ssh-rsa AAAAB3Nza...`.|
|`file_path`|string|The path of the file to download on the remote host.|
|`scp_path`|string|Optionally specify a SCP binary to use. If omitted pukcab will search $PATH.|

## Example

```json
{
    "name": "scp",
    "config": {
        "host_address": "10.0.0.1",
        "port": 22,
        "username": "example",
        "private_key": "-----BEGIN OPENSSH PRIVATE KEY-----\n<omitted>\n-----END OPENSSH PRIVATE KEY-----",
        "host_public_key": "ssh-ed25519 AAAAC3<omitted>",
        "file_path": "/config/config.boot"
    }
}
```
