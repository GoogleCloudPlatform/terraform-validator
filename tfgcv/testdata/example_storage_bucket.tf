resource "google_storage_bucket" "image-store" {
  name     = "image-store-bucket"
  location = "EU"

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
}
