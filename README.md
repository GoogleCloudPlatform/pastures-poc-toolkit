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
sudo wget https://github.com/GoogleCloudPlatform/pastures-poc-toolkit/releases/download/v0.7.2/pastures_amd64.tar.gz -O - \
| sudo tar -zxf - -C /usr/local/bin

sudo chmod +x /usr/local/bin/pasture
```
<!-- x-release-please-end -->
## Quickstart

**Note: We recommend running this quickstart from a Cloud Shell environment**

1. Configure your local Pastures environment by:
    - Defining a prefix for resource naming
    - Specifying which Cloud Identity group will own the PoC
    - Specifying your GCP Organization domain
    - Specifying your GCP Billing Account

```shell
pasture plow --prefix example1 --group-owner pasture-group --domain example.com --billing-account ABCDEF-GHIJKL-MNOPQ
```

2. Create a pasture by indicating which seed template and optional jumpstarts you'd like to deploy:

**Note: This could take up to 10 minutes to deploy**

```shell
pasture plant data-cloud --region us-central1 --pasture-size small --data-warehouse
```

## Cleanup

Destruction of a planted pasture is scoped to the seed template. All resources deployed by `pasture` or out of band will be deleted. Currently, `pasture burn` requires the same paramters inputs that were used with the corresponding `pasture plant`:

```shell
pasture burn data-cloud --region us-central1 --pasture-size small --data-warehouse
```

## State

Pastures will make every effort to persist environmental state and variable values to a GCS bucket. This supports the ability to run `pasture` from multiple locations, but demands a rehydration step in order to position the dependencies. Rehydration simply requires authorization with Google and the `prefix` originally chosen:

```shell
pasture plow --rehydrate --prefix example1
```

Afterwards, you can continue running `pasture` as your normally would.

## Pasture Templates

| Name | Description | Docs | Est. Price Calculator |
| ---- | ----------- | ---- | --------------------- |
| `data-cloud` | Landing zone for data, analytics and generative AI | [cmd](docs/pasture_plant_data-cloud.md) | [Small](https://cloud.google.com/products/calculator?hl=en&dl=CiQ0Yzc1N2RjNC0yN2QyLTQyMmEtODlkZS0xYzkxNzAyM2JmNTgQCxokMzIxQTgxMTctN0Q3NC00QUU4LUE4NzAtNTJFNDIxMUMyNEYx)<br>[Big](https://cloud.google.com/products/calculator?hl=en&dl=CiRhYzk2Y2MzZS05ZWRkLTRmMDAtYWM5OS1lYmVmN2UyYjY0NTEQCxokRDkwRDA3OTEtNDE2Qi00ODNFLUJFRjctMjU3RTEwNkRCQzE5) |

## Jumpstarts

| Name | Seed | Docs |
| ---- | ---- | ---- |
| Data Warehouse | [`data-cloud`](docs/pasture_plant_data-cloud.md) | [Docs](https://cloud.google.com/architecture/big-data-analytics/data-warehouse) |
| Analytics Lakehouse | [`data-cloud`](docs/pasture_plant_data-cloud.md) | [Docs](https://cloud.google.com/architecture/big-data-analytics/analytics-lakehouse) |
| Knowledge Base | [`data-cloud`](docs/pasture_plant_data-cloud.md) | [Docs](https://cloud.google.com/architecture/ai-ml/generative-ai-knowledge-base) |
| GenAI RAG | [`data-cloud`](docs/pasture_plant_data-cloud.md) | [Docs](https://cloud.google.com/architecture/ai-ml/generative-ai-rag) |

## Learn More

- [FAQ](docs/faq.md)
- [Configure Pasture](docs/pasture_plow.md)
- [Create a PoC](docs/pasture_plant.md)
- [Delete a PoC](docs/pasture_burn.md)
- [What is FAST Foundation?](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/v29.0.0/fast/README.md)