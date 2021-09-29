# PHP support for DocFX YAML

This project contains a converter to go from PHPDocumentor XML output to
DocFX YAML.

Usage:

```
go install github.com/googleapis/phpdocyaml
phpdocyaml -namespace '\Google\Cloud\Vision' -version 1.0.0 -structure path/to/structure.xml
```

See `phpdocyaml -h` for more usage information.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Disclaimer

This is not an official Google Product. It may change in backwards-incompatible
ways at any time without warning.