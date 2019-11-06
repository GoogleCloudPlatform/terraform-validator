resource "google_project_iam_binding" "project" {
  project = "foobar"
  role    = "roles/editor"

  members = [
    "user:jane@example.com",
  ]
}
