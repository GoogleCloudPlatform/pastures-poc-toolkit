# FAQ

### What is the purpose of Pastures?

Pastures is intended to accelerate and simplify full-featured experimentation on Google Cloud using production-ready foundations that can be quickly promoted or destroyed once completed. It relies extensively on established frameworks and jumpstarts to provide a workable environment for most use cases, and eliminates landing zone toil/rework that is often encountered when promoting a PoC use case to production.

### What does Pastures deploy?

Pastures deploys three discrete subassemblies:

1. [FAST](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/README.md) Foundation - Bootstrap and Resource Management stages are deployed with secure-by-default configurations to establish a hierarchy for subsequent stages to be deployed

2. Seed Template Projects - An opinionated implementation of a Google Cloud solution pillar, e.g. Data Cloud, Open Cloud, etc. These seeds provide a full-featured, proof-of-concept landing zone for use case experimentation

3. [Optional] Seed Jumpstarts - To further accelerate use case prototyping, each seed can optionally be paired with a [Google Cloud Jumpstart](https://cloud.google.com/architecture/all-jss-guides). These jumpstarts provide best-practice solutions designed with common use cases in mind.

### What is the Fabric FAST foundation?

The spirit of Pastures is to accelerate experimentation by leveraging established frameworks, shorterning the use case experimentation lifecycle from PoC -> pilot -> production. We believe that the [Cloud Foundation Fabric](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric) FAST implementation delivers exactly that; rapid, secure-by-default landing zones to host all lifecycles of a workload. Pastures deploys [`stage-0`](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/0-bootstrap/README.md) and [`stage-1`](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/1-resman/README.md) to provide sane and secure organizational defaults, and a top-level folder (`Sandbox`) to host pasture deployments.

### Why are a GCP Organization and Billing Account required?

These are basic primitives of GCP [required](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/1-resman/README.md) by the FAST foundation.

### Why are the prerequisite permissions required?

These permissions are [required](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/stages/1-resman/README.md) by the FAST foundation in order to bootstrap a greenfield GCP organization. You can read more about how FAST leverages IAM authoritative policies to protect the organization and it's resources.

### Where is a pasture's state persisted?

Each of the FAST stages' remote state are stored in discrete GCS buckets. The Seed state and Pasture vars files are stored in the common FAST `outputs` GCS bucket. This persistence approach provides portability for Pastures, although we highly recommend running in a Cloud Shell environment.

### What are Google Cloud Jumpstarts?

Jumpstarts are preconditioned resource collections by the Cloud Architecture Center. Take a look at this [Google blog post](https://cloud.google.com/blog/products/application-modernization/introducing-google-cloud-jump-start-solutions) covering the essence of Jumpstarts

### What all is deleted when I "burn" a pasture?

`pasture burn` will delete a seed template project, any optional Jumpstarts deployed, and any manual resources created inside the seed project. Currently, a destroy sequence _will not_ delete a FAST foundation. This is to accommodate future experimentation in a safe-by-default environment, since a foundation is always recommended for enterprise workloads on GCP.
