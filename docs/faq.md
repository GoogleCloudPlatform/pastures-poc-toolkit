# FAQ

### What is the purpose of Pastures?

Pastures is intended to accelerate and simplify full-featured experimentation on Google Cloud using production-ready foundations that can be quickly promoted or destroyed once completed. It relies extensively on established frameworks and blueprints to provide a workable environment for most use cases, and eliminates landing zone toil/rework that is often encountered when promoting a PoC use case to production.

### What does Pastures deploy?

Pastures deploys two discrete subassemblies:

1. [FAST](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/README.md) Foundation - Bootstrap and Resource Management stages are deployed with secure-by-default configurations to establish a hierarchy for subsequent stages to be deployed

2. Seed Blueprints - To further accelerate use case prototyping, each seed is paired with a [Google Cloud Blueprint](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/blueprints/README.md). These blueprints provide best-practice solutions designed with common use cases in mind.

### What is the Fabric FAST foundation?

The spirit of Pastures is to accelerate experimentation by leveraging established frameworks, shorterning the use case experimentation lifecycle from PoC -> pilot -> production. We believe that the [Cloud Foundation Fabric](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric) FAST implementation delivers exactly that; rapid, secure-by-default landing zones to host all lifecycles of a workload. Pastures deploys [`stage-0`](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/0-bootstrap/README.md) and [`stage-1`](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/1-resman/README.md) to provide sane and secure organizational defaults, and a top-level folder (`Sandbox`) to host pasture deployments.

### Why are a GCP Organization and Billing Account required?

These are basic primitives of GCP [required](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/1-resman/README.md) by the FAST foundation.

### Why are the prerequisite permissions required?

These permissions are [required](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/1-resman/README.md) by the FAST foundation in order to bootstrap a greenfield GCP organization. You can read more about how FAST leverages IAM authoritative policies to protect the organization and it's resources.

### Where is a pasture's state persisted?

Each of the FAST stages' remote state are stored in discrete GCS buckets. The Seed state and Pasture vars files are stored in the common FAST `outputs` GCS bucket. This persistence approach provides portability for Pastures, although we highly recommend running in a Cloud Shell environment.

### What are Fabric Blueprints?

Blueprints are preconditioned resource collections maintained in the Cloud Foundation Fabric [repository](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/tree/master/blueprints). Take a look at the documentation which covers the essence of Blueprints.

### What all is deleted when I destroy a pasture?

`pasture destroy` will delete a seed template project, any seed Blueprints deployed, and any manual resources created inside the seed project. Currently, a destroy sequence _will not_ delete a FAST foundation. This is to accommodate future experimentation in a safe-by-default environment, since a foundation is always recommended for enterprise workloads on GCP.
