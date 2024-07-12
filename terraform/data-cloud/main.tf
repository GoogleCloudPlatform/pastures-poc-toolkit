/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

locals {
  dimensions = {
    "small" = {
      slots = 100
      ram   = 16
    }
    "large" = {
      slots = 200
      ram   = 32
    }
  }
}

resource "random_string" "random" {
  length  = 4
  special = false
  upper   = false
}

resource "google_folder" "data_cloud" {
  display_name = "pasture-data-cloud"
  parent       = data.google_active_folder.sandbox.name
}

module "projects" {
  source = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//blueprints/factories/project-factory?ref=v29.0.0"

  data_defaults = {
    billing_account = var.billing_account.id
    parent          = data.google_active_folder.sandbox.name
  }

  data_merges = {
    labels = {
        source = "pastures"
    }
    services = [
      "logging.googleapis.com",
      "monitoring.googleapis.com",
      "stackdriver.googleapis.com",
      "iam.googleapis.com",
      "serviceusage.googleapis.com",
      "servicemanagement.googleapis.com",
      "cloudapis.googleapis.com",
      "cloudresourcemanager.googleapis.com"
    ]
  }

  data_overrides = {
    prefix = "pasture-${var.prefix}-${random_string.random.result}"
  }

  factory_data_path = "data/projects"
}

resource "google_bigquery_reservation" "reservation" {
  project = module.data-platform.projects.project_id.processing

  name              = "pastures-data-cloud"
  location          = var.locations.bq
  slot_capacity     = local.dimensions[var.pasture_size].slots
  edition           = "ENTERPRISE_PLUS"
  ignore_idle_slots = false
  concurrency       = 0
  autoscale {
    max_slots = 400
  }
}

resource "google_bigquery_reservation_assignment" "assignment" {
  project = module.data-platform.projects.project_id.processing

  assignee    = "projects/${module.data-platform.projects.project_id.processing}"
  job_type    = "QUERY"
  reservation = google_bigquery_reservation.reservation.id
}

resource "google_bigquery_bi_reservation" "bi_reservation" {
  project = module.data-platform.projects.project_id.curated

  location = var.locations.bq
  size     = local.dimensions[var.pasture_size].ram * pow(1024, 3)
}

module "data-platform" {
  source              = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//blueprints/data-solutions/data-platform-minimal?ref=v29.0.0"
  organization_domain = var.organization.domain
  project_config = {
    billing_account_id = var.billing_account.id
    parent             = google_folder.data_cloud.folder_id
  }
  prefix         = var.prefix
  project_suffix = "-${random_string.random.result}"

  groups = {
    data-analysts  = google_cloud_identity_group.data_analysts.display_name
    data-engineers = google_cloud_identity_group.data_engineers.display_name
    data-security  = google_cloud_identity_group.data_security.display_name
  }

  location = var.locations.bq
  region   = var.region
}
