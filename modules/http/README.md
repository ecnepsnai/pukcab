# Pukcab module: HTTP

Module name: `http`

This module enables you to download a file and save it as an artifact.

# Requirements

*None*

# Configuration

|Key|Type|Description|
|---|----|-----------|
|`url`|string|The URL to fetch.|
|`file_name`|string|The output file name to save.|
|`allow_untrusted_certificates`|boolean|(Optional) If true, untrusted certificates are allowed.|
|`headers`|map string -> string|(Optional) Map of headers to add to the HTTP request.|

## Example

```json
{
    "name": "http",
    "config": {
        "url": "https://github.com/ecnepsnai/pukcab/archive/refs/heads/main.zip",
        "file_name": "pukcab.zip",
        "headers": {
            "X-Example-Header": "Header Value"
        }
    }
}
```
