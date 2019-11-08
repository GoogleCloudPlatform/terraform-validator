resource "google_compute_instance" "full_list_default_1" {
  # Required arguments
  name         = "test1"
  machine_type = "n1-standard-1"
  boot_disk {
    auto_delete             = true
    device_name             = "test-device_name"
    disk_encryption_key_raw = "test-disk_encryption_key_raw"
    initialize_params {
      # TODO: panic in google.resolveImageImageExists if it is not a global image
      image = "projects/debian-cloud/global/images/debian-9"
      size  = 42
      type  = "pd-standard"
    }
    # TODO: generate new test case.
    # Got: "boot_disk.0.kms_key_self_link": conflicts with boot_disk.0.disk_encryption_key_raw
    # kms_key_self_link = "test-kms_key_self_link"
    # TODO: report a bug.
    # Got: An argument named "mode" is not expected here.
    # mode = "READ_ONLY"
    source = "test-source"
  }
  network_interface {
    access_config {
      nat_ip = "test-nat_ip"
    }
    access_config {
      network_tier = "STANDARD"
    }
    access_config {
      public_ptr_domain_name = "test-public_ptr_domain_name"
    }
    alias_ip_range {
      ip_cidr_range         = "test-ip_cidr_range"
      subnetwork_range_name = "test-subnetwork_range_name"
    }
    network    = "default"
    network_ip = "test-network_ip"
  }
  network_interface {
    subnetwork         = "test-subnetwork"
    subnetwork_project = "test-subnetwork_project"
  }
  # Optional arguments
  allow_stopping_for_update = true
  attached_disk {
    # Required arguments
    source = "test-source"
    # Optional arguments
    device_name             = "test-device_name"
    disk_encryption_key_raw = "test-disk_encryption_key_raw"
    kms_key_self_link       = "test-kms_key_self_link"
    mode                    = "READ_ONLY"
  }
  attached_disk {
    source = "test-source2"
  }
  can_ip_forward      = true
  deletion_protection = true
  description         = "test-description"
  guest_accelerator {
    type  = "test-guest_accelerator-type1"
    count = 42
  }
  guest_accelerator {
    type  = "test-guest_accelerator-type2"
    count = 42
  }
  hostname = "test-hostname"
  labels = {
    label_foo1 = "label-bar1"
  }
  metadata = {
    metadata_foo1 = "metadata-bar1"
  }
  # TODO: metadata_startup_script mix up with metadata. Need to test
  # with a new test case.
  # metadata_startup_script { }
  min_cpu_platform = "test-min_cpu_platform"
  scheduling {
    preemptible         = true
    on_host_maintenance = "test-on_host_maintenance"
    automatic_restart   = true
    node_affinities {
      key      = "test-key"
      operator = "IN"
      values   = ["test-values1", "test-values2"]
    }
  }
  scratch_disk {
    interface = "SCSI"
  }
  scratch_disk {
    interface = "SCSI"
  }
  service_account {
    email  = "test-email"
    scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  }
  shielded_instance_config {
    enable_secure_boot          = true
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }
  tags = ["foo", "bar"]
  zone = "us-central1-a"
}

# test boot_disk.kms_key_self_link
resource "google_compute_instance" "full_list_default_2" {
  name         = "test2"
  machine_type = "n1-standard-1"
  boot_disk {
    # TODO: this is not found in the result asset.
    kms_key_self_link = "test-kms_key_self_link"
  }
  network_interface {
    network = "default"
  }
  zone = "us-central1-a"
}
