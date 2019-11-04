resource "google_project_iam_member" "project" {
  project = "foobar"
  role    = "roles/editor"
  member  = "user:jane@example.com"
}
