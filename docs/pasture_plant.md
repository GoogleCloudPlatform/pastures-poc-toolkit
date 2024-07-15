## pasture plant

Creates a POC environment from a template

### Synopsis

Plant creates a POC environment in a FAST foundation sandbox using
a seed template (e.g. data-cloud). Example:
	
	pasture plant data-cloud --region us-central1 --pasture-size small
	
A list of seed templates is shown by running:
	
	pasture plant --help

### Options

```
      --dry-run   Displays the desired state of the POC
  -h, --help      help for plant
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
```

### SEE ALSO

* [pasture](pasture.md)	 - A POC toolkit for Google Cloud
* [pasture plant data-cloud](pasture_plant_data-cloud.md)	 - Deploy a Data Cloud pasture with optional jumpstarts
