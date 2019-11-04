resource "google_project_iam_policy" "project" {
  project     = "foobar"
  policy_data = "${data.google_iam_policy.admin.policy_data}"
}

data "google_iam_policy" "admin" {
  binding {
    role = "roles/editor"

    members = [
      "user:jane@example.com",
    ]
  }
}
