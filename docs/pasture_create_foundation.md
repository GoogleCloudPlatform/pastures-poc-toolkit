## pasture create foundation

Deploy a foundation-only pasture with no blueprints

### Synopsis

Creates a foundation landing zone from the FAST framework.
Projects can optionally be deployed as features into the landing zone. An
example of how to use this pasture:
	
	pasture create foundation

```
pasture create foundation [flags]
```

### Options

```
  -h, --help   help for foundation
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
      --dry-run         Displays the desired state of the POC
      --verbose         controls Terraform output verbosity (default "false")
```

### SEE ALSO

* [pasture create](pasture_create.md)	 - Creates a POC environment from a template
