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

variable "organization" {
  description = "Organization details."
  type = object({
    domain      = string
    id          = number
    customer_id = string
  })
}

variable "prefix" {
  description = "Prefix used for resources that need unique names. Use 9 characters or less."
  type        = string
  validation {
    condition     = try(length(var.prefix), 0) < 10
    error_message = "Use a maximum of 9 characters for prefix."
  }
}

variable "billing_account" {
  description = "Billing account id. If billing account is not part of the same org set `is_org_level` to `false`. To disable handling of billing IAM roles set `no_iam` to `true`."
  type = object({
    id           = string
    is_org_level = optional(bool, true)
    no_iam       = optional(bool, false)
  })
  nullable = false
}

variable "locations" {
  description = "Optional locations for GCS, BigQuery, and logging buckets created here."
  type = object({
    bq      = optional(string, "EU")
    gcs     = optional(string, "EU")
    logging = optional(string, "global")
    pubsub  = optional(list(string), [])
  })
  nullable = false
  default  = {}
}

# variable "groups" {
#   # https://cloud.google.com/docs/enterprise/setup-checklist
#   description = "Group names or emails to grant organization-level permissions. If just the name is provided, the default organization domain is assumed."
#   type = object({
#     gcp-billing-admins      = string
#     gcp-devops              = string
#     gcp-network-admins      = string
#     gcp-organization-admins = string
#     gcp-security-admins     = string
#     gcp-support             = string
#   })
#   nullable = false
# }

variable "region" {
  description = "GCP region for regional resources to be hosted"
  type        = string
  nullable    = false
  default     = "us-central1"
}

variable "state_bucket" {
  description = "GCS bucket to hold pasture state file"
  type        = string
  nullable    = false
}

variable "state_dir" {
  description = "Top-level directory in GCS bucket where state will be stored"
  type        = string
  nullable    = false
}

variable "pasture_size" {
  description = "Course-grain sizing estimate of pasture"
  type        = string
  nullable    = false
  default     = "small"
}