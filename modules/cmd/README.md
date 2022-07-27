# Pukcab module: CMD

Module name: `cmd`

This module enables you to save the output of a command as a file.

# Configuration

|Key|Type|Description|
|---|----|-----------|
|`exec_path`|string|The path to the executable to run.|
|`args`|array|An array of arguments to pass to the executable.|
|`env`|array|(Optional) An array of `key=value` pair strings to append to the current environment variables. Will overwrite any duplicates using the values from this array.|
|`wd`|string|(Optional) The working directory for the executable.|
|`output_name`|string|The name of the output file.|
|`include_stderr`|bool|Optionally include the output of stderr in the artifact.|

## Example

```json
{
    "name": "cmd",
    "config": {
        "exec_path": "/usr/bin/example",
        "args": ["-C"],
        "env": ["FOO=bar"],
        "wd": "/var/log",
        "output_name": "example.txt"
    }
}
```
