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

resource "google_folder_iam_member" "folder_admin" {
  folder = data.google_active_folder.sandbox.id
  role   = "roles/resourcemanager.folderAdmin"
  member = "group:${var.groups.gcp-organization-admins}@${var.organization.domain}"
}

resource "google_folder_iam_member" "project_creator" {
  folder = data.google_active_folder.sandbox.id
  role   = "roles/resourcemanager.projectCreator"
  member = "group:${var.groups.gcp-organization-admins}@${var.organization.domain}"  
}

resource "google_folder_iam_member" "owner" {
  folder = data.google_active_folder.sandbox.id
  role   = "roles/owner"
  member = "group:${var.groups.gcp-organization-admins}@${var.organization.domain}"  
}