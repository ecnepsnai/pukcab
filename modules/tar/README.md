# Pukcab module: Tape Archive (tar)

Module name: `tar`

This module enables you to create a gzipped-tarball of given source files or directories

# Requirements

- The `tar` or `gtar` executable must be installed on the backup host

# Configuration

|Key|Type|Description|
|---|----|-----------|
|`tar_path`|string|(Optional) Path to tar executable to use. Defaults to `tar`.|
|`tarball_name`|string|The output tarball name. Be sure to include the `tar.gz` or `.tgz` extension.|
|`sources`|[]string|Array of paths to add to the tarball.|

## Example

```json
{
    "name": "tar",
    "config": {
        "tarball_name": "network-scripts.tar.gz",
        "sources": [
            "/etc/sysconfig/network-scripts/"
        ],
    }
}
```
