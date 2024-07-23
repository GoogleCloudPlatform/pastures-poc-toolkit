## pasture

A POC toolkit for Google Cloud

### Synopsis

Pasture is a CLI toolkit that will create POC
	landing zones in your Google Cloud organization. It relies
	on the Cloud Foundation Fabric framework to establish a GCP
	foundation, and will deploy each 'pasture' as a Sandbox project.

### Options

```
      --config string   config file (default is $HOME/.pastures.yaml)
  -h, --help            help for pasture
  -t, --toggle          Help message for toggle
      --verbose         controls Terraform output verbosity (default "false")
```

### SEE ALSO

* [pasture completions](pasture_completions.md)	 - Generate completion script
* [pasture configure](pasture_configure.md)	 - Initializes environment configuration
* [pasture create](pasture_create.md)	 - Creates a POC environment from a template
* [pasture destroy](pasture_destroy.md)	 - Removes the POC resources created by a seed.
* [pasture version](pasture_version.md)	 - Displays Pasture binary version
