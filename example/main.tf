/**
 * Copyright 2019 Google LLC
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
  project = var.project_id
}

provider "google" {
  version = "~> 2.11.0"
  project = local.project
}

resource "google_folder" "test" {
  display_name = "validator-test"
  parent       = "organizations/${var.org_id}"
}

resource "google_project" "my_project" {
  name       = "test-project"
  project_id = "test-project-912"
  folder_id  = "folders/880142660913"

  labels  = {
    "project-label-key-a" = "project-label-val-a"
  }

  # org_id     = "organizations/${var.org_id}"
  # folder_id  = google_folder.test.name
}

resource "google_compute_disk" "my-disk" {
  name    = "my-disk"
  project = local.project
  type    = "pd-ssd"
  zone    = "us-central1-a"
  image   = "debian-8-jessie-v20170523"

  labels = {
    foo = "bar"
  }
}

resource "google_compute_firewall" "my-test-firewall" {
  name    = "my-test-firewall"
  network = "default"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }

  source_tags = ["web"]
}

resource "random_id" "bucket" {
  byte_length = 8
}

resource "google_storage_bucket" "my-bucket" {
  name     = "my-bucket-${random_id.bucket.hex}"
  project  = local.project
  location = "US"

  labels = {
    foo = "bar"
  }

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }

  cors {
    origin = ["*"]
    method = ["POST"]
  }
}

/* Uncomment and change emails to try out IAM policies.
resource "google_project_iam_member" "owner-a" {
  project = "${local.project}"
  role    = "roles/owner"
  member  = "user:example-a@google.com"
}

resource "google_project_iam_member" "viewer-a" {
  project = "${local.project}"
  role    = "roles/viewer"
  member  = "user:example-a@google.com"
}

resource "google_project_iam_member" "viewer-b" {
  project = "${local.project}"
  role    = "roles/viewer"
  member  = "user:example-b@google.com"
}

resource "google_project_iam_binding" "editors" {
  project = "${local.project}"
  role    = "roles/editor"
  members  = [
    "user:example-a@google.com",
    "user:example-b@google.com"
  ]
}
*/
