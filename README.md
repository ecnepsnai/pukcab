# Pukcab

*It's _backup_ spelled backwards*

Pukcab is a modular backup solution for simple, file-based backup needs.

## Modules

- **Cloudflare**: backup DNS BIND zones files for all domains associated with a Cloudflare account.
- **HTTP**: backup any file specified by a URL using HTTP.
- **pfSense**: backup the configuration of a pfSense device.
- **SCP**: backup any file over SCP.
- **TAR**: backup any local file or directory into a TAR (tape archive).
- **CMD**: backup the output from any command.

*See the README file in each module in the `modules` directory for detailed information and configuration.*

## Usage

Pukcab is controlled using a JSON configuration file with the following properties:

|Key|Type|Description|
|---|----|-----------|
|`modules`|array|Array of modules and their associated configuration. The same module can be repeated multiple times.|
|`output_dir`|string|The directory where files should be saved.|
|`artifact_retention`|number|The number of days for backed up files to be retained.|

For example:

```json
{
    "modules": [],
    "output_dir": "/mnt/backup",
    "artifact_retention": 5
}
```

Then, run pukcab with that configuration file: `./pukcab config.json`
