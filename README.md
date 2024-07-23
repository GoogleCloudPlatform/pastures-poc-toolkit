# Pastures PoC Toolkit

![release](https://img.shields.io/github/v/release/googlecloudplatform/pastures-poc-toolkit) ![license](https://img.shields.io/github/license/GoogleCloudPlatform/pastures-poc-toolkit)

The Pastures Proof of Concept (PoC) utility is used to bootstrap greenfield, production-ready landing zones for use case experimentation on [Google Cloud](https://cloud.google.com/).

![demo-gif](assets/demo.gif)

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://shell.cloud.google.com/cloudshell/editor?cloudshell_git_repo=http://github.com/GoogleCloudPlatform/pastures-poc-toolkit.git)

## Prerequisites
1. [Google Cloud Organization](https://cloud.google.com/resource-manager/docs/cloud-platform-resource-hierarchy#organizations)
2. [Google Cloud Billing Account](https://cloud.google.com/billing/docs/how-to/manage-billing-account)
3. [Cloud Identity Group](https://support.google.com/cloudidentity/answer/9400082)
4. The following IAM permissions:
    - Billing Account Administrator (`roles/billing.admin`) either on the organization or the billing account (see the following section for details)
    - Logging Admin (`roles/logging.admin`)
    - Organization Role Administrator (`roles/iam.organizationRoleAdmin`)
    - Organization Administrator (`roles/resourcemanager.organizationAdmin`)
    - Project Creator (`roles/resourcemanager.projectCreator`)
    - Tag Admin (`roles/resourcemanager.tagAdmin`)
    - Owner (`roles/owner`)

The following bash script can be used to quickly assign these permissions to your account:

> [!NOTE]
> These prerequisite roles are automatically assigned to the pasture group you supply as a flag to the `configure` subcommand

```shell
# set variable for current logged in user
export PASTURE_USER=$(gcloud config list --format 'value(core.account)')

# find and set your org id
gcloud organizations list
export ORG_ID=123456

# set needed roles
export PASTURE_ROLES="roles/billing.admin roles/logging.admin \
  roles/iam.organizationRoleAdmin roles/resourcemanager.projectCreator \
  roles/resourcemanager.organizationAdmin roles/resourcemanager.tagAdmin \
  roles/owner"

for role in $PASTURE_ROLES; do
  gcloud organizations add-iam-policy-binding $ORG_ID \
    --member user:$PASTURE_USER --role $role --condition None
done
```

We recommend running `pasture` from a Cloud Shell environment. If you choose to run from your local machine, the following packages are also required:
1. [Terraform](https://developer.hashicorp.com/terraform/install)
2. [gcloud SDK](https://cloud.google.com/sdk/docs/install)

## Install

Install the binary to your `$PATH` of choice. `amd64` is currently the only supported build architecture.
<!-- x-release-please-start-version -->
```shell
sudo wget https://github.com/GoogleCloudPlatform/pastures-poc-toolkit/releases/download/v1.1.1/pastures_amd64.tar.gz -O - \
| sudo tar -zxf - -C /usr/local/bin

sudo chmod +x /usr/local/bin/pasture
```
<!-- x-release-please-end -->
## Quickstart

We recommend running this quickstart from a Cloud Shell environment.

1. Configure your local Pastures environment by:
    - Defining a prefix for resource naming
    - Specifying which Cloud Identity group e.g. `pasture-group` will own the PoC
    - Specifying your GCP Organization domain
    - Specifying your GCP Billing Account

> [!IMPORTANT]
> Ensure that the user running the CLI is a member of your Cloud Identity group e.g. `pasture-group`.

```shell
pasture configure --prefix example1 --group-owner pasture-group --domain example.com --billing-account ABCDEF-GHIJKL-MNOPQ
```

2. Create a pasture by indicating which seed template you'd like to deploy (could take ~15 mins to complete):

```shell
pasture create data-cloud --region us-central1 --pasture-size small
```

## Cleanup

Destruction of a pasture is scoped to the seed template. All resources deployed by `pasture` or out of band will be deleted. Currently, `pasture destroy` requires the same paramters inputs that were used with the corresponding `pasture create`:

```shell
pasture destroy data-cloud --region us-central1 --pasture-size small
```

## Known Issues

See [Known Issues](docs/known_issues.md).

## State

Pastures will make every effort to persist environmental state and variable values to a GCS bucket. This supports the ability to run `pasture` from multiple locations, but demands a rehydration step in order to position the dependencies. Rehydration simply requires authorization with Google and the `prefix` originally chosen:

```shell
pasture configure --rehydrate --prefix example1
```

Afterwards, you can continue running `pasture` as your normally would.

## Pasture Templates

| Name | Description | Docs | Est. Price Calculator |
| ---- | ----------- | ---- | --------------------- |
| `data-cloud` | Landing zone for data, analytics and generative AI | [cmd](docs/pasture_create_data-cloud.md) | [Small](https://cloud.google.com/products/calculator-legacy#id=5c5c2811-605e-4bdd-94f6-d1c9a19defd5)<br>[Big](https://cloud.google.com/products/calculator-legacy#id=ab352e16-69de-4726-8e91-f1fe0475c3dc) |
| `foundation` | Generic landing zone from Fabric FAST foundation stage 0 and stage 1 | [cmd](docs/pasture_create_foundation.md) | N/A |

## Blueprints

| Name | Seed | Docs |
| ---- | ---- | ---- |
| Data Platform | [`data-cloud`](docs/pasture_create_data-cloud.md) | [Docs](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/tree/master/blueprints/data-solutions/data-platform-foundations) |

## Learn More

- [FAQ](docs/faq.md)
- [Configure Pasture](docs/pasture_configure.md)
- [Create a PoC](docs/pasture_create.md)
- [Delete a PoC](docs/pasture_destroy.md)
- [What is FAST Foundation?](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/master/fast/README.md)