# firehose-stats

This is a plugin for the Cloud Foundry CLI to provide relatively real time metrics and statistics from the firehose.

The purpose of this plugin is to provide a quick and easy way to visually understand relevant statistics about the
Cloud Foundry deployment. This negates the need for having a system like Elasticsearch setup to consume all the metrics.

## Usage

User must be logged in as an admin.

```
cf firehose-stats
```
or

```
cf fs
```

## Install
```
cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org
cf install-plugin "FirehoseStats" -r CF-Community
```

If you are doing development and want to install it locally, run `scripts/install.sh`


## Uninstall
```
cf uninstall FirehoseStats
```

## Testing

Run `scripts/test.sh`