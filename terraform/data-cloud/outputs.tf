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
  _tpl_providers = "${path.module}/templates/providers.tf.tpl"
  providers = {
    "data-cloud" = templatefile(local._tpl_providers, {
      name   = "data-cloud"
      bucket = var.state_bucket
    })
  }
}

resource "google_storage_bucket_object" "providers" {
  for_each = local.providers
  bucket   = var.state_bucket
  name     = "${var.state_dir}/${each.key}-providers.tf"
  content  = each.value
}

output "project_id" {
  description = "Details of pcreated projects"
  value       = module.data-platform.projects.project_id.landing
}
