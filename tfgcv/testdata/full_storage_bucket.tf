resource "google_storage_bucket" "full-list-default" {
  name     = "image-store-bucket"
  location = "EU"

  bucket_policy_only = true
  cors {
    origin          = ["test-origin1", "test-origin2"]
    method          = ["test-method1", "test-method2"]
    response_header = ["test-response_header1", "test-response_header2"]
    max_age_seconds = 42
  }
  cors {
    origin          = ["test-origin1", "test-origin2"]
    method          = ["test-method1", "test-method2"]
    response_header = ["test-response_header1", "test-response_header2"]
    max_age_seconds = 42
  }
  encryption {
    default_kms_key_name = "test-default_kms_key_name"
  }
  force_destroy = true
  labels = {
    label_foo1 = "label-bar1"
  }
  lifecycle_rule {
    action {
      type          = "test-type"
      storage_class = "test-storage_class"
    }
    condition {
      age                   = 42
      created_before        = "test-created_before"
      is_live               = true
      matches_storage_class = ["test-matches_storage_class1", "matches_storage_class2"]
      num_newer_versions    = 42
      with_state            = "LIVE"
    }
  }
  lifecycle_rule {
    action {
      type          = "test-type"
      storage_class = "test-storage_class"
    }
    condition {
      age                   = 42
      created_before        = "test-created_before"
      is_live               = true
      matches_storage_class = ["test-matches_storage_class1", "matches_storage_class2"]
      num_newer_versions    = 42
      with_state            = "LIVE"
    }
  }
  logging {
    log_bucket        = "test-log_bucket"
    log_object_prefix = "test-log_object_prefix"
  }
  requester_pays = true
  retention_policy {
    is_locked        = true
    retention_period = 42
  }
  storage_class = "test-storage_class"
  versioning {
    enabled = true
  }
  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}
