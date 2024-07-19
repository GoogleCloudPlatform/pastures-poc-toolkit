## pasture burn

Removes the POC resources created by a seed.

### Synopsis

Removes the POC resources created by a seed in a previous
run of the plant command. Example:

	pasture burn data-cloud --jumpstart data-warehouse
	
A list of seed templates is shown by running:
	
	pasture burn --help

### Options

```
      --dry-run   Displays the desired state of the POC
  -h, --help      help for burn
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
```

### SEE ALSO

* [pasture](pasture.md)	 - A POC toolkit for Google Cloud
* [pasture burn data-cloud](pasture_burn_data-cloud.md)	 - Deploy a Data Cloud pasture with blueprints
