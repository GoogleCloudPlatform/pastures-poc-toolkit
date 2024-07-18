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

> [!TIP]
You could simply re-run your `pasture` command and not re-encounter the error.

Pastures is intended to be a short-lived PoC execution. The chances for error recurrence becomes less and less after you deploy stages 0 and 1 of Cloud Foundation Fabric.

For a more comprehensive solution, you may also try:

1. To use a [workaround](https://stackoverflow.com/a/62827358) to force Google API calls in Cloud Shell to use an IP from the private.googleapis.com range (199.36.153.8/30 )
1. To deploy the foundation code from a local machine that supports IPv6.