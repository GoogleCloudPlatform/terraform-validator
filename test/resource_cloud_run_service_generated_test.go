// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package test

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestAccCloudRunService_cloudRunServiceBasicExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceBasicExample_offline"
	offline := true
	testAccCloudRunService_cloudRunServiceBasicExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudRunServiceBasicExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceBasicExample_online"
	offline := false
	testAccCloudRunService_cloudRunServiceBasicExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudRunServiceBasicExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudRunServiceBasicExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudRunServiceBasicExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tf-test-cloudrun-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, context)
}

func TestAccCloudRunService_cloudRunServiceSqlExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceSqlExample_offline"
	offline := true
	testAccCloudRunService_cloudRunServiceSqlExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudRunServiceSqlExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceSqlExample_online"
	offline := false
	testAccCloudRunService_cloudRunServiceSqlExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudRunServiceSqlExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudRunServiceSqlExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudRunServiceSqlExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tf-test-cloudrun-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "1000"
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.instance.connection_name
        "run.googleapis.com/client-name"        = "terraform"
      }
    }
  }
  autogenerate_revision_name = true
}

resource "google_sql_database_instance" "instance" {
  name             = "tf-test-cloudrun-sql%{random_suffix}"
  region           = "us-east1"
  database_version = "MYSQL_5_7"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection  = "%{deletion_protection}"
}
`, context)
}

func TestAccCloudRunService_cloudRunServiceConfigurationExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceConfigurationExample_offline"
	offline := true
	testAccCloudRunService_cloudRunServiceConfigurationExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudRunServiceConfigurationExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceConfigurationExample_online"
	offline := false
	testAccCloudRunService_cloudRunServiceConfigurationExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudRunServiceConfigurationExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudRunServiceConfigurationExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudRunServiceConfigurationExample(context map[string]interface{}) string {
	return Nprintf(`
# Example configuration of a Cloud Run service

resource "google_cloud_run_service" "default" {
  name     = "config%{random_suffix}"
  location = "us-central1"

    template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"

        # Container "entry-point" command
        # https://cloud.google.com/run/docs/configuring/containers#configure-entrypoint
        command = ["/server"]

        # Container "entry-point" args
        # https://cloud.google.com/run/docs/configuring/containers#configure-entrypoint
        args = []
        
        # [START cloudrun_service_configuration_http2]
        # Enable HTTP/2
        # https://cloud.google.com/run/docs/configuring/http2
        ports {
          name           = "h2c"
          container_port = 8080
        }
        # [END cloudrun_service_configuration_http2]

                # Environment variables
        # https://cloud.google.com/run/docs/configuring/environment-variables
        env {
          name  = "foo"
          value = "bar"
        }
        env {
          name  = "baz"
          value = "quux"
        }
        
                        resources {
          limits = {
            # CPU usage limit
            # https://cloud.google.com/run/docs/configuring/cpu
            cpu = "1000m" # 1 vCPU

            # Memory usage limit (per container)
            # https://cloud.google.com/run/docs/configuring/memory-limits
            memory = "512Mi"
          }
        }
                
              }
      
            # Timeout
      # https://cloud.google.com/run/docs/configuring/request-timeout
      timeout_seconds = 300
      
            # Maximum concurrent requests
      # https://cloud.google.com/run/docs/configuring/concurrency
      container_concurrency = 80
      
          }
    
                metadata {
            annotations = {

        # Max instances
        # https://cloud.google.com/run/docs/configuring/max-instances
        "autoscaling.knative.dev/maxScale" = 10

        # Min instances
        # https://cloud.google.com/run/docs/configuring/min-instances
        "autoscaling.knative.dev/minScale" = 1

        # If true, garbage-collect CPU when once a request finishes
        # https://cloud.google.com/run/docs/configuring/cpu-allocation
        "run.googleapis.com/cpu-throttling" = false
      }
            
            # Labels
      # https://cloud.google.com/run/docs/configuring/labels
      labels = {
        foo : "bar"
        baz : "quux"
      }
                }
              }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, context)
}

func TestAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_offline"
	offline := true
	testAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_online"
	offline := false
	testAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudRunServiceMultipleEnvironmentVariablesExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tf-test-cloudrun-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
        env {
          name = "SOURCE"
          value = "remote"
        }
        env {
          name = "TARGET"
          value = "home"
        }
      }
    }
  }

  metadata {
    annotations = {
      generated-by = "magic-modules"
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
  autogenerate_revision_name = true

  lifecycle {
    ignore_changes = [
        metadata.0.annotations,
    ]
  }
}
`, context)
}

func TestAccCloudRunService_cloudRunServiceScheduledExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceScheduledExample_offline"
	offline := true
	testAccCloudRunService_cloudRunServiceScheduledExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudRunServiceScheduledExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudRunServiceScheduledExample_online"
	offline := false
	testAccCloudRunService_cloudRunServiceScheduledExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudRunServiceScheduledExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"project":       getTestProjectFromEnv(),
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudRunServiceScheduledExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudRunServiceScheduledExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project_service" "run_api" {
  project                    = "%{project}"
  service                    = "run.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

resource "google_project_service" "iam_api" {
  project                    = "%{project}"
  service                    = "iam.googleapis.com"
  disable_on_destroy         = false
}

resource "google_project_service" "resource_manager_api" {
  project                    = "%{project}"
  service                    = "cloudresourcemanager.googleapis.com"
  disable_on_destroy         = false
}

resource "google_project_service" "scheduler_api" {
  project                    = "%{project}"
  service                    = "cloudscheduler.googleapis.com"
  disable_on_destroy         = false
}

resource "google_cloud_run_service" "default" {
  project  = "%{project}"
  name     = "tf-test-my-scheduled-service%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_project_service.run_api
  ]
}

resource "google_service_account" "default" {
  project      = "%{project}"
  account_id   = "tf-test-scheduler-sa%{random_suffix}"
  description  = "Cloud Scheduler service account; used to trigger scheduled Cloud Run jobs."
  display_name = "scheduler-sa"

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_project_service.iam_api
  ]
}

resource "google_cloud_scheduler_job" "default" {
  name             = "tf-test-scheduled-cloud-run-job%{random_suffix}"
  description      = "Invoke a Cloud Run container on a schedule."
  schedule         = "*/8 * * * *"
  time_zone        = "America/New_York"
  attempt_deadline = "320s"

  retry_config {
    retry_count = 1
  }

  http_target {
    http_method = "POST"
    uri         = google_cloud_run_service.default.status[0].url

    oidc_token {
      service_account_email = google_service_account.default.email
    }
  }

  # Use an explicit depends_on clause to wait until API is enabled
  depends_on = [
    google_project_service.scheduler_api
  ]
}

resource "google_cloud_run_service_iam_member" "default" {
  project = "%{project}"
  location = google_cloud_run_service.default.location
  service = google_cloud_run_service.default.name
  role = "roles/run.invoker"
  member = "serviceAccount:${google_service_account.default.email}"
}
`, context)
}

func TestAccCloudRunService_cloudrunServiceAccessControlExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudrunServiceAccessControlExample_offline"
	offline := true
	testAccCloudRunService_cloudrunServiceAccessControlExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudrunServiceAccessControlExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudrunServiceAccessControlExample_online"
	offline := false
	testAccCloudRunService_cloudrunServiceAccessControlExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudrunServiceAccessControlExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudrunServiceAccessControlExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudrunServiceAccessControlExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tf-test-cloud-run-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloud_run_service.default.location
  service  = google_cloud_run_service.default.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}
`, context)
}

func TestAccCloudRunService_cloudRunSystemPackagesExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudRunSystemPackagesExample_offline"
	offline := true
	testAccCloudRunService_cloudRunSystemPackagesExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudRunSystemPackagesExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudRunSystemPackagesExample_online"
	offline := false
	testAccCloudRunService_cloudRunSystemPackagesExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudRunSystemPackagesExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudRunSystemPackagesExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudRunSystemPackagesExample(context map[string]interface{}) string {
	return Nprintf(`
# Example of how to deploy a Cloud Run application with system packages

resource "google_cloud_run_service" "default" {
  name     = "tf-test-graphviz-example%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        # Replace with the URL of your graphviz image
        #   gcr.io/<YOUR_GCP_PROJECT_ID>/graphviz
        image = "gcr.io/cloudrun/hello"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Make Cloud Run service publicly accessible
resource "google_cloud_run_service_iam_member" "allow_unauthenticated" {
  service  = google_cloud_run_service.default.name
  location = google_cloud_run_service.default.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
`, context)
}

func TestAccCloudRunService_cloudrunServiceIdentityExample_generated_offline(t *testing.T) {
	testSlug := "CloudRunService_cloudrunServiceIdentityExample_offline"
	offline := true
	testAccCloudRunService_cloudrunServiceIdentityExample_shared(t, testSlug, offline)
}

func TestAccCloudRunService_cloudrunServiceIdentityExample_generated_online(t *testing.T) {
	testSlug := "CloudRunService_cloudrunServiceIdentityExample_online"
	offline := false
	testAccCloudRunService_cloudrunServiceIdentityExample_shared(t, testSlug, offline)
}

func testAccCloudRunService_cloudrunServiceIdentityExample_shared(t *testing.T, testSlug string, offline bool) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}

	t.Parallel()
	context := map[string]interface{}{
		"random_suffix": "meepmerp", // true randomization isn't needed for validator
	}

	terraformConfig := getTestPrefix() + testAccCloudRunService_cloudrunServiceIdentityExample(context)
	dir, err := ioutil.TempDir(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dstFile := path.Join(dir, "main.tf")
	err = os.WriteFile(dstFile, []byte(terraformConfig), 0666)
	if err != nil {
		t.Fatalf("error while writing to file %s, error %v", dstFile, err)
	}

	terraformWorkflow(t, dir, testSlug)
	if offline && shouldOutputGeneratedFiles() {
		generateTFVconvertedAsset(t, dir, testSlug)
		return
	}

	// need to have comparison.. perhaps test vs checked in code
	// testConvertCommand(t, dir, c.name, offline, c.compareConvertOutput)

	testValidateCommandGeneric(t, dir, testSlug, offline, true)
}

func testAccCloudRunService_cloudrunServiceIdentityExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "cloudrun_service_identity" {
  account_id   = "my-service-account"
}

resource "google_cloud_run_service" "default" {
  name     = "tf-test-cloud-run-srv%{random_suffix}"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/cloudrun/hello"
      }
      service_account_name = google_service_account.cloudrun_service_identity.email  
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}
`, context)
}
