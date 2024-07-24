# Known Issues

## Cannot assign requested address error in Cloud Shell

Error message:

When using [Google Cloud Shell](https://cloud.google.com/shell/docs) to run Pastures, you may encounter an error like

```
dial tcp [2607:f8b0:400c:c15::5f]:443: connect: cannot assign requested address
```

when Terraform calls the Google APIs.

### Cause:

This is a [known terraform issue](https://github.com/hashicorp/terraform-provider-google/issues/6782) regrading IPv6.

### Solution:


You may try:

1. To use a [workaround](https://stackoverflow.com/a/62827358) to force Google API calls in Cloud Shell to use an IP from the private.googleapis.com range (199.36.153.8/30 )
2. To deploy the foundation code from a local machine that supports IPv6.

Pastures is intended to be a short-lived PoC execution. The chances for error recurrence becomes less and less after you deploy stages 0 and 1 of Cloud Foundation Fabric. Therefore, you just re-run the pasture command.

> [!TIP]
You could simply re-run your `pasture` command and not re-encounter the error.


## Error waiting for Deleting Network

Error message:

```
The network resource 'projects/pasture-yourprefix-xxxxxxxxxx/global/networks/yourprefix-lod' is already being used by 'projects/pasture-yourprefix-xxxxxxxxxx-lod/global/addresses/cdf-pasture-datafusion'
```

### Cause:

This is a transient error encountered during `pasture destroy`. An underlying `terraform destroy` operation attempts to remove the Data Fusion deployment before the prequisite destroy operations have completed.

### Solution:

Re-run your `pasture` command, and you won't re-encounter the error.