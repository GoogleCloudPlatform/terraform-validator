provider "google" {
  version = "~> 1.20"
  credentials = "{{.Provider.credentials}}"
}
