resource "google_organization_iam_member" "binding" {
  org_id = "0123456789"
  role    = "roles/editor"
  member  = "user:alice@gmail.com"
}
