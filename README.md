# Ansible Galaxy collections grabber

A simple CLI tool to download Ansible Galaxy Collections locally for offline/airgapped use.

## Build

The `go.mod` file requires Go 1.21+ but there are no specific requirements on the latest compiler version; it can be built by running `make` in the project's root directory.

## Usage

To download a set of collections as listed in an input JSON file (`input.json`) to a specified directory (`./output`):

```bash
$> ./ansible-galaxy-grabber --collections=@input.json --directory=./output
```

The input file can be in either JSON or YAML format and should specify the namespace and the name of the collections to download.
It can optionally provide one or more comma-separated constraints on the versions of the collection; when povided, the constraints are used to filter out versions that should not be downloaded; the version criteria can be specified as per https://blog.gopheracademy.com/advent-2015/semver/

An example JSON file is as follows:
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

The same example in YAML format:

```yaml
---
- namespace: crivetimihai
  collection: virtualization
  constraint: ">= 1.0.17, < 1.0.24, != 1.0.20"
- namespace: crivetimihai
  collection: development
  constraint: "= 1.0.1"
```

NOTE: the input JSON and YAML file can be either provided as a file on disk (in which case you need to provide the path prefixed with `@`, as shown in the example above) or be passed inline. 
When passing an inline YAML, make sure that it starts with `---` (as in the example above) or the command line parser won't be able to autodetect the encoding.


