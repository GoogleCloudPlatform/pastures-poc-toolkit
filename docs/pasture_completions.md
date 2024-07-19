## pasture completions

Generate completion script

### Synopsis

To load completions:

	Bash:

	  $ source <(completions completion bash)

	  # To load completions for each session, execute once:
	  # Linux:
	  $ completions completion bash > /etc/bash_completion.d/completions
	  # macOS:
	  $ completions completion bash > $(brew --prefix)/etc/bash_completion.d/completions

	Zsh:

	  # If shell completion is not already enabled in your environment,
	  # you will need to enable it.  You can execute the following once:

	  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

	  # To load completions for each session, execute once:
	  $ completions completion zsh > "${fpath[1]}/_completions"

	  # You will need to start a new shell for this setup to take effect.

	fish:

	  $ completions completion fish | source

	  # To load completions for each session, execute once:
	  $ completions completion fish > ~/.config/fish/completions/completions.fish

	PowerShell:

	  PS> completions completion powershell | Out-String | Invoke-Expression

	  # To load completions for every new session, run:
	  PS> completions completion powershell > completions.ps1
	  # and source this file from your PowerShell profile.
	

```
pasture completions [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completions
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
```

### SEE ALSO

* [pasture](pasture.md)	 - A POC toolkit for Google Cloud
