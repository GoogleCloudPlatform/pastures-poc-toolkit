## pasture plant data-cloud

Deploy a Data Cloud pasture with optional jumpstarts

### Synopsis

Creates a data-cloud landing zone in a FAST foundation sandbox.
Jumpstarts can optionally be deployed as features into the landing zone. An
example of how to use this pasture:
	
	pasture plant data-cloud --jumpstart data-warehouse

```
pasture plant data-cloud [flags]
```

### Options

```
      --analytics-lakehouse   Enable the Analytics lakehouse jumpstart
      --data-warehouse        Enable the BigQuery data warehouse jumpstart
      --genai-rag             Enable the Vertex AI RAG jumpstart
  -h, --help                  help for data-cloud
      --knowledge-base        Enable the Vertex AI knowledge base jumpstart
  -s, --pasture-size string   Size of pasture environment - must be 'big' or 'small'
  -r, --region string         Region for GCP resources to be deployed (default "us-central1")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.pastures.yaml)
      --dry-run         Displays the desired state of the POC
```

### SEE ALSO

* [pasture plant](pasture_plant.md)	 - Creates a POC environment from a template
