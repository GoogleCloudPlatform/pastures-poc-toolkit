## pasture create

Creates a POC environment from a template

### Synopsis

Create instantiates a POC environment in a FAST foundation sandbox using
a seed template (e.g. data-cloud). Example:
	
	pasture create data-cloud --region us-central1 --pasture-size small
	
A list of seed templates is shown by running:
	
	pasture create --help

### Options

```
      --dry-run   Displays the desired state of the POC
  -h, --help      help for create
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
      --verbose         controls Terraform output verbosity (default "false")
```

### SEE ALSO

* [pasture](pasture.md)	 - A POC toolkit for Google Cloud
* [pasture create data-cloud](pasture_create_data-cloud.md)	 - Deploy a Data Cloud pasture with blueprints
* [pasture create foundation](pasture_create_foundation.md)	 - Deploy a foundation-only pasture with no blueprints
