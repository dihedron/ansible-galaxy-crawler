# Ansible Galaxy collections grabber

A simple CLI tool to download Ansible Galaxy Collections and store them for offline/airgapped use.

## Build

Requires Go 1.21+; it can be build by running `make` in the project's root directory.

## Usage

To download a set of collections as listed in an input JSON file (`input.json`) to directory a specified directory (`./output`):

```bash
$> ./ansible-galaxy-grabber --collections=@input.json --directory=./output
```

The input file can be in either JSON or YAML format and should specify the name of the namespace and the collection; it can optionally provide one (or more, comma separated) constraints of the verison of the collection, as per https://blog.gopheracademy.com/advent-2015/semver/

A sample JSON file is as follows:
```json
[
    {
        "namespace": "crivetimihai",
        "collection": "virtualization",
        "constraint": ">= 1.0.17, < 1.0.24, != 1.0.20"
    },
    {
        "namespace": "crivetimihai",
        "collection": "development",
        "constraint": "= 1.0.1"
    }
]
```

and the corresponding YAML:

```yaml
---
- namespace: crivetimihai
  collection: virtualization
  constraint: ">= 1.0.17, < 1.0.24, != 1.0.20"
- namespace: crivetimihai
  collection: development
  constraint: "= 1.0.1"
```

NOTE: the input JSON and YAML file can be either provided as a file on disk (in which case you need to provide the path prefixed with `@`, as shown in the example above) or be passed inline. When passing an inline YAML, make sure that it starts with `---` or the command line parser won't be able to autodetect the encoding.


