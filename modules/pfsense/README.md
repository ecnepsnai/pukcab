# Pukcab module: PFSense

Module name: `pfsense`

This module enables you toback up the configuration from a running PFSense router.

# Requirements

- The backup host must be able to access the PFSense devices' web UI.
- The PFSense device must be configured to use TLS. Plain-Text configurations are not supported.

# Configuration

|Key|Type|Description|
|---|----|-----------|
|`host_address`|string|The host address of the PFSense device. Do not include a protocol or path.|
|`username`|string|The username to log in as.|
|`password`|string|The password for that user.|
|`allow_untrusted_certificates`|boolean|(Optional) If true then TLS errors are ignored.|
|`encrypt_password`|string|(Optional) If included the backup will be encrypted with this password.|

## Example

```json
{
    "name": "pfsense",
    "config": {
        "host_address": "192.168.1.1",
        "username": "exampleUser",
        "password": "examplePassword",
        "encrypt_password": "not_a_good_password",
    }
}
```
