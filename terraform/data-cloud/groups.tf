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

provider "google" {
  alias = "google-group"

  billing_project       = module.projects.projects["cmn"].id
  user_project_override = true
}

resource "google_cloud_identity_group" "data_analysts" {
  provider = google.google-group

  display_name         = "gcp-data-analysts"
  initial_group_config = "WITH_INITIAL_OWNER"

  parent = "customers/${var.organization.customer_id}"

  group_key {
      id = "gcp-data-analysts@${var.organization.domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group" "data_engineers" {
  provider = google.google-group

  display_name         = "gcp-data-engineers"
  initial_group_config = "WITH_INITIAL_OWNER"

  parent = "customers/${var.organization.customer_id}"

  group_key {
      id = "gcp-data-engineers@${var.organization.domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_cloud_identity_group" "data_security" {
  provider = google.google-group

  display_name         = "gcp-data-security"
  initial_group_config = "WITH_INITIAL_OWNER"

  parent = "customers/${var.organization.customer_id}"

  group_key {
      id = "gcp-data-security@${var.organization.domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}