# Pukcab module: Cloudflare

Module name: `cloudflare`

This module enables you toback up DNS zone files for all zones on a Cloudflare account.

# Requirements

- The backup host must be able to access the Cloudflare API

# Configuration

|Key|Type|Description|
|---|----|-----------|
|`cloudflare_email`|string|The email address of your Cloudflare user account.|
|`cloudflare_api_key`|string|The API key for your Cloudflare user account.|

## Example

```json
{
    "name": "cloudflare",
    "config": {
        "cloudflare_email": "example@example.com",
        "cloudflare_api_key": "123456789",
    }
}
```
