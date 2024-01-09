package cdm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/job"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
)

func getCdmJobResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CdmV11Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDM v1.1 client, err=%s", err)
	}
	clusterId, jobName, err := cdm.ParseJobInfoFromId(state.Primary.ID)
	if err != nil {
		return nil, err
	}

	return job.Get(client, clusterId, jobName, job.GetJobsOpts{})
}

func TestAccResourceCdmJob_basic(t *testing.T) {
	var obj job.JobCreateOpts
	resourceName := "huaweicloud_cdm_job.test"
	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmJob_basic(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "job_type", "NORMAL_JOB"),
					resource.TestCheckResourceAttr(resourceName, "source_connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "destination_connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "config.0.retry_type", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "config.0.scheduler_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_job_config.fromFileOpType", "DO_NOTHING"),
					resource.TestCheckResourceAttr(resourceName, "destination_job_config.outputFormat", "BINARY_FILE"),
					resource.TestCheckResourceAttr(resourceName, "destination_job_config.duplicateFileOpType", "REPLACE"),
				),
			},
			{
				Config: testAccCdmJob_update(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "job_type", "NORMAL_JOB"),
					resource.TestCheckResourceAttr(resourceName, "source_connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "destination_connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "config.0.retry_type", "RETRY_TRIPLE"),
					resource.TestCheckResourceAttr(resourceName, "config.0.scheduler_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_job_config.fromFileOpType", "DO_NOTHING"),
					resource.TestCheckResourceAttr(resourceName, "destination_job_config.outputFormat", "BINARY_FILE"),
					resource.TestCheckResourceAttr(resourceName, "destination_job_config.duplicateFileOpType", "REPLACE"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_job_config", "destination_job_config"},
			},
		},
	})
}

func testAccCdmLinkAndOBS(name, bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "input" {
  bucket        = "%[1]s-input"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "output" {
  bucket        = "%[1]s-output"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "dirty" {
  bucket        = "%[1]s-dirty"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cdm_link" "test" {
  name       = "%[2]s"
  connector  = "obs-connector"
  cluster_id = huaweicloud_cdm_cluster.test.id
  enabled    = true

  config = {
    "storageType" = "OBS"
    "server"      = trimprefix(huaweicloud_obs_bucket.output.bucket_domain_name, "${huaweicloud_obs_bucket.output.bucket}.")
    "port"        = "443"
  }

  access_key = "%[3]s"
  secret_key = "%[4]s"
}
`, bucketName, name, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testAccCdmJob_basic(name, bucketName string) string {
	clusterConfig := testAccCdmCluster_basic(name)
	linkAndOBSConfig := testAccCdmLinkAndOBS(name, bucketName)

	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cdm_job" "test" {
  name       = "%s"
  job_type   = "NORMAL_JOB"
  cluster_id = huaweicloud_cdm_cluster.test.id

  source_connector = "obs-connector"
  source_link_name = huaweicloud_cdm_link.test.name
  source_job_config = {
    "bucketName"               = huaweicloud_obs_bucket.input.bucket
    "inputDirectory"           = "/"
    "listTextFile"             = "false"
    "inputFormat"              = "BINARY_FILE"
    "fromCompression"          = "NONE"
    "fromFileOpType"           = "DO_NOTHING"
    "useMarkerFile"            = "false"
    "useTimeFilter"            = "false"
    "fileSeparator"            = "|"
    "filterType"               = "NONE"
    "useWildCard"              = "false"
    "decryption"               = "NONE"
    "nonexistentPathDisregard" = "false"
  }

  destination_connector = "obs-connector"
  destination_link_name = huaweicloud_cdm_link.test.name
  destination_job_config = {
    "bucketName"          = huaweicloud_obs_bucket.output.bucket
    "outputDirectory"     = "/"
    "outputFormat"        = "BINARY_FILE"
    "validateMD5"         = "true"
    "recordMD5Result"     = "false"
    "duplicateFileOpType" = "REPLACE"
    "useCustomDirectory"  = "false"
    "encryption"          = "NONE"
    "copyContentType"     = "false"
    "shouldClearTable"    = "false"
  }

  config {
    retry_type                          = "NONE"
    scheduler_enabled                   = true
    scheduler_cycle_type                = "month"
    scheduler_cycle                     = 1
    scheduler_run_at                    = 1
    scheduler_start_date                = "2024-01-04 12:00:00"
    scheduler_disposable_type           = "NONE"
    throttling_extractors_number        = 4
    throttling_record_dirty_data        = true
    throttling_dirty_write_to_link      = huaweicloud_cdm_link.test.name
    throttling_dirty_write_to_bucket    = huaweicloud_obs_bucket.dirty.bucket
    throttling_dirty_write_to_directory = "/"
    throttling_max_error_records        = 10
    throttling_loader_number            = 1
  }

  lifecycle {
    ignore_changes = [
      source_job_config, destination_job_config,
    ]
  }
}
`, clusterConfig, linkAndOBSConfig, name)
}

func testAccCdmJob_update(name, bucketName string) string {
	clusterConfig := testAccCdmCluster_basic(name)
	linkAndOBSConfig := testAccCdmLinkAndOBS(name, bucketName)

	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_cdm_job" "test" {
  name       = "%s"
  job_type   = "NORMAL_JOB"
  cluster_id = huaweicloud_cdm_cluster.test.id

  source_connector = "obs-connector"
  source_link_name = huaweicloud_cdm_link.test.name
  source_job_config = {
    "bucketName"               = huaweicloud_obs_bucket.input.bucket
    "inputDirectory"           = "/"
    "listTextFile"             = "false"
    "inputFormat"              = "BINARY_FILE"
    "fromCompression"          = "NONE"
    "fromFileOpType"           = "DO_NOTHING"
    "useMarkerFile"            = "false"
    "useTimeFilter"            = "false"
    "fileSeparator"            = "|"
    "filterType"               = "NONE"
    "useWildCard"              = "false"
    "decryption"               = "NONE"
    "nonexistentPathDisregard" = "false"
  }

  destination_connector = "obs-connector"
  destination_link_name = huaweicloud_cdm_link.test.name
  destination_job_config = {
    "bucketName"          = huaweicloud_obs_bucket.output.bucket
    "outputDirectory"     = "/"
    "outputFormat"        = "BINARY_FILE"
    "validateMD5"         = "true"
    "recordMD5Result"     = "false"
    "duplicateFileOpType" = "REPLACE"
    "useCustomDirectory"  = "false"
    "encryption"          = "NONE"
    "copyContentType"     = "false"
    "shouldClearTable"    = "false"
  }

  config {
    retry_type                          = "RETRY_TRIPLE"
    scheduler_enabled                   = true
    scheduler_cycle_type                = "month"
    scheduler_cycle                     = 1
    scheduler_run_at                    = 1
    scheduler_start_date                = "2024-01-04 12:00:00"
    scheduler_disposable_type           = "NONE"
    throttling_extractors_number        = 4
    throttling_record_dirty_data        = true
    throttling_dirty_write_to_link      = huaweicloud_cdm_link.test.name
    throttling_dirty_write_to_bucket    = huaweicloud_obs_bucket.dirty.bucket
    throttling_dirty_write_to_directory = "/"
    throttling_max_error_records        = 10
    throttling_loader_number            = 1
  }

  lifecycle {
    ignore_changes = [
      source_job_config, destination_job_config,
    ]
  }
}
`, clusterConfig, linkAndOBSConfig, name)
}
