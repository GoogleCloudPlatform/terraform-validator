resource "google_storage_bucket_iam_member" "member2" {
  bucket = "test-bucket"
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
