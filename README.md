# Pastures PoC Toolkit

![release](https://img.shields.io/github/v/release/googlecloudplatform/pastures-poc-toolkit) ![license](https://img.shields.io/github/license/GoogleCloudPlatform/pastures-poc-toolkit)

The Pastures Proof of Concept (PoC) utility is used to bootstrap greenfield, production-ready landing zones for use case experimentation on [Google Cloud](https://cloud.google.com/).

![demo-gif](assets/demo.gif)

[![Open in Cloud Shell](https://gstatic.com/cloudssh/images/open-btn.svg)](https://shell.cloud.google.com/cloudshell/editor?cloudshell_git_repo=http://github.com/GoogleCloudPlatform/pastures-poc-toolkit.git)

## 1. Setup Prerequisites

### Google Cloud
1. Create your [Google Cloud Organization](https://cloud.google.com/resource-manager/docs/cloud-platform-resource-hierarchy#organizations)
2. Turn on your [Google Cloud Billing Account](https://cloud.google.com/billing/docs/how-to/manage-billing-account)
3. Create a free [Cloud Identity](https://support.google.com/cloudidentity/answer/7389973) user e.g. `myadmin@example.com` within that Google Cloud organization who has the [Organization Policy Administrator](https://cloud.google.com/iam/docs/understanding-roles#orgpolicy.policyAdmin) IAM role.
4. Create a [Cloud Identity Group](https://support.google.com/cloudidentity/answer/9400082) e.g. `pastures-group` where your user e.g. `myadmin@example.com` is a member of that group.

### Command Line
Pastures requires the following command line utilities:
- [Terraform <sup>1</sup>](https://developer.hashicorp.com/terraform/install)
- [gcloud SDK](https://cloud.google.com/sdk/docs/install)

<hr>

> [!IMPORTANT]
 <sup>1</sup> Cloud Shell installs `terraform` and `gcloud` binaries by default, but the `terraform` binary defaults to an older version that won't work with Pastures. You may quickly override the default `terraform` binary by copying and pasting these shell commands in your Cloud Shell:

Run the following in your terminal to fix Cloud Shell.

```shell
cd ~
curl -O https://releases.hashicorp.com/terraform/1.9.4/terraform_1.9.4_linux_amd64.zip
unzip -o terraform_1.9.4_linux_amd64.zip
sudo mv terraform /usr/bin
```

## 2. Install Pastures CLI

Install the binary to your `$PATH` of choice. `amd64` is currently the only supported build architecture. That means that you can run Pastures on Cloud Shell, Linux, and the Windows Subsystem for Linux. macOS universal binary support is available in the [Makefile](Makefile).

Run the following in your terminal to install Pastures.
<!-- x-release-please-start-version -->
```shell
BASE_URL=https://github.com/GoogleCloudPlatform/pastures-poc-toolkit
RELEASE_URL=releases/download/v1.1.3/pastures_amd64.tar.gz
sudo wget -q $BASE_URL/$RELEASE_URL -O - | sudo tar -zxf - -C /usr/local/bin
sudo chmod +x /usr/local/bin/pasture
echo -e "\n\033[32mYou are now ready to use Pastures..."
echo -e "`pasture version`\n\033[0m\n"
```
<!-- x-release-please-end -->

## 3. Execute Data Cloud Quickstart

We recommend running this quickstart from a Cloud Shell environment.

> [!IMPORTANT]
Ensure that the user running the CLI e.g. `myadmin@example.com` is a member of your Cloud Identity group e.g. `pasture-group`.

1. Configure your local Pastures environment by:
    - Defining a prefix for resource naming
    - Specifying which Cloud Identity group e.g. `pasture-group` will own the PoC
    - Specifying your GCP Organization domain
    - Specifying your GCP Billing Account

```shell
pasture configure \
--prefix example1 \
--group-owner pasture-group \
--domain example.com \
--billing-account ABCDEF-GHIJKL-MNOPQ
```

2. Create a pasture by indicating which seed template you'd like to deploy (could take ~15 mins to complete):

```shell
pasture create data-cloud \
--region us-central1 \
--pasture-size small
```

## 4. Cleanup

Destruction of a pasture is scoped to the seed template. All resources deployed by `pasture` or out of band will be deleted. Currently, `pasture destroy` requires the same paramters inputs that were used with the corresponding `pasture create`:

<!-- TODO Do we need --pasture-size during destroy? Won't region be enough? Isn't one seed deployable per region? -->

```shell
pasture destroy data-cloud \
--region us-central1 \
--pasture-size small
```

## Known Issues

See [Known Issues](docs/known_issues.md).

## State

Pastures will make every effort to persist environmental state and variable values to a GCS bucket. This supports the ability to run `pasture` from multiple locations, but demands a rehydration step in order to position the dependencies. Rehydration simply requires authorization with Google and the `prefix` originally chosen:

```shell
pasture configure \
--rehydrate \
--prefix example1
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