# statsd-exporter-convert

`statsd-exporter-convert` assists with converting from the legacy text based [statsd_exporter](https://github.com/prometheus/statsd_exporter) configuration to the YAML based one.

## Usage

Clone this repository and build it using Go 1.9+. [Binary releases](https://github.com/bakins/statsd-exporter-convert/releases) are also available.

Run it with a single file as the only argument:

```shell
statsd-exporter-convert <path/to/my/mappings>
```

It will print the YAML to STDOUT.

Note: only minor validation is done on the input.

## Example

If we converted the  [included mapping file](./mapping.txt) which contains:

```
test.dispatcher.*.*.*
name="dispatcher_events_total"
processor="$1"
action="$2"
outcome="$3"
job="test_dispatcher"

*.signup.*.*
name="signup_events_total"
provider="$2"
outcome="$3"
job="${1}_server"
```

We would get 

```yaml
mappings:
- match: test.dispatcher.*.*.*
  name: dispatcher_events_total
  labels:
    action: $2
    job: test_dispatcher
    outcome: $3
    processor: $1
- match: '*.signup.*.*'
  name: signup_events_total
  labels:
    job: ${1}_server
    outcome: $3
    provider: $2
```

## LICENSE
see [LICENSE](./LICENSE)


