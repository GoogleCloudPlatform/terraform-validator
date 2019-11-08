provider "google" {
  version     = "~> {{.Provider.version}}"
  credentials = "{{.Provider.credentials}}"
}

resource "google_storage_bucket" "my-test-bucket" {
  name     = "test-bucket"
  location = "EU"

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}
