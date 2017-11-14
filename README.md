# statsd-exporter-convert

`statsd-exporter-convert` assits with converting from the legacy text based [statsd_exporter](https://github.com/prometheus/statsd_exporter) configuration to the YAML based one.

## Usage

Clone this repository and build it using Go 1.9+.

Run it with a single file as the only argument:

```shell
statsd-exporter-convert <path/to/my/mappings>
```

It will print the YAML to STDOUT.

Note: only minor validation is done on the input.

## LICENSE
see [LICENSE](./LICENSE)


