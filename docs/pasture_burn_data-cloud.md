## pasture burn data-cloud

Deploy a Data Cloud pasture with blueprints

### Synopsis

Creates a data-cloud landing zone in a FAST foundation sandbox.
Blueprints are deployed as features into the landing zone. An
example of how to use this pasture:
	
	pasture plant data-cloud --region us-central1 --pasture-size small

```
pasture burn data-cloud [flags]
```

### Options

```
  -h, --help                  help for data-cloud
  -s, --pasture-size string   Size of pasture environment - must be 'big' or 'small'
  -r, --region string         Region for GCP resources to be deployed (default "us-central1")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
      --dry-run         Displays the desired state of the POC
      --verbose         controls Terraform output verbosity (default "false")
```

### SEE ALSO

* [pasture burn](pasture_burn.md)	 - Removes the POC resources created by a seed.
