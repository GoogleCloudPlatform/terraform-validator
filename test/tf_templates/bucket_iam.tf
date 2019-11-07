resource "google_storage_bucket_iam_member" "member" {
  bucket = "test-bucket"
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
