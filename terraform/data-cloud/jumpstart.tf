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
  is_big = var.pasture_size == "large" ? true : false
}

# module "doc_summary" {
#   count  = var.enable_summarization == "true" && !local.is_big ? 1 : 0
#   source = "github.com/GoogleCloudPlatform/terraform-genai-doc-summarization?ref=v0.1.1"

#   project_id = module.projects.projects["data"].id
#   region     = var.region
# }

module "doc_summary_with_tuning" {
  count  = var.enable_summarization == "true" && local.is_big ? 1 : 0
  source = "github.com/GoogleCloudPlatform/terraform-genai-knowledge-base"

  project_id          = module.projects.projects["data"].id
  region              = var.region
  documentai_location = lower(var.locations.bq)
}

module "rag" {
  count  = var.enable_rag == "true" ? 1 : 0
  source = "github.com/GoogleCloudPlatform/terraform-genai-rag"

  project_id = module.projects.projects["data"].id
  region     = var.region
}

# TODO: swap out the below source for tf registry on v8 or higher

module "warehouse" {
  count   = var.enable_warehouse == "true" ? 1 : 0
  source = "github.com/terraform-google-modules/terraform-google-bigquery//modules/data_warehouse"
  # source  = "terraform-google-modules/bigquery/google//modules/data_warehouse"
  # version = "~> 7.0"

  project_id      = module.projects.projects["data"].id
  region          = var.region
  dataform_region = var.region
}

module "data_science" {
  count = var.enable_analytics == "true" ? 1 : 0
  source = "github.com/GoogleCloudPlatform/terraform-google-analytics-lakehouse?ref=v0.4.0"

  project_id      = module.projects.projects["data"].id
  region          = var.region
}