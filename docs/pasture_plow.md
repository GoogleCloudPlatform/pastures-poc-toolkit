## pasture plow

Initializes environment configuration

### Synopsis

This command will create an environment and define its properties in a pasture configuration file, which is located by default at $HOME/.pastures/pasture.yaml.

```
pasture plow [flags]
```

### Options

```
  -b, --billing-account string   GCP billing account ID
  -d, --domain string            GCP organization domain name
      --fabric-version string    Cloud Foundation Fabric FAST version (default "v29.0.0")
  -g, --group-owner string       Name of Cloud Identity group that owns the pastures
  -h, --help                     help for plow
  -l, --location string          GCP multi-region location code (default "US")
  -p, --prefix string            Prefix for resources with unique names (max 9 characters)
      --rehydrate                Restore previous Pastures configuration from saved version in GCS bucket
      --seed-version string      Version of pasture seed terraform modules to use
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
      --verbose         controls Terraform output verbosity (default "false")
```

### SEE ALSO

* [pasture](pasture.md)	 - A POC toolkit for Google Cloud
