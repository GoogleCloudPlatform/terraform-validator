resource "google_organization_iam_binding" "binding" {
  org_id = "123456789"
  role    = "roles/browser"

  members = [
    "user:alice@gmail.com",
  ]
}
