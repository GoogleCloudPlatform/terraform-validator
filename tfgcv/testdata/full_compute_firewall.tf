resource "google_compute_firewall" "full_list_default_1" {
  name    = "test-firewall1"
  network = "${google_compute_network.default.name}"

  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }
  description        = "test-description"
  destination_ranges = ["test-destination_ranges1", "test-destination_ranges2"]
  direction          = "INGRESS"
  disabled           = true
  # TODO: beta feature
  # Got: An argument named "enable_logging" is not expected here.
  # enable_logging = true
  priority = 42
}

resource "google_compute_firewall" "full_list_default_2" {
  name    = "test-firewall2"
  network = "${google_compute_network.default.name}"

  deny {
    protocol = "icmp"
  }
  deny {
    protocol = "tcp"
    ports    = ["80", "8080", "1000-2000"]
  }
  source_ranges           = ["test-source_range1", "test-source_range2"]
  source_service_accounts = ["test-source_service_account1", "test-source_service_account2"]
  target_service_accounts = ["test-target_service_account1", "test-target_service_account2"]
}

resource "google_compute_firewall" "full_list_default_3" {
  name    = "test-firewall3"
  network = "${google_compute_network.default.name}"

  deny {
    protocol = "icmp"
  }
  source_tags = ["web"]
  target_tags = ["test-target_tag1", "test-target_tag2"]
}

resource "google_compute_network" "default" {
  name = "test-network"
}
