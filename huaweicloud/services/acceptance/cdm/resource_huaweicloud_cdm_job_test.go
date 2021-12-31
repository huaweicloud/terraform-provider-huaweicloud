package cdm

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/job"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getCdmJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.CdmV11Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating CDM v1.1 client, err=%s", err)
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
				Config: testAccCdmJob_basic(name, bucketName, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "job_type", "NORMAL_JOB"),
					resource.TestCheckResourceAttr(resourceName, "source_connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "destination_connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "config.0.retry_type", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "config.0.scheduler_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "source_job_config.fromFileOpType", "DO_NOTHING"),
					resource.TestCheckResourceAttr(resourceName, "destination_job_config.outputFormat", "BINARY_FILE"),
					resource.TestCheckResourceAttr(resourceName, "destination_job_config.duplicateFileOpType", "REPLACE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCdmJob_basic(name, bucketName, ak, sk string) string {
	clusterConfig := testAccCdmCluster_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "input" {
  bucket        = "%s-input"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "output" {
  bucket        = "%s-output"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cdm_link" "test" {
  name       = "%s"
  connector  = "obs-connector"
  cluster_id = huaweicloud_cdm_cluster.test.id
  enabled    = true

  config = {
    "storageType" = "OBS"
    "server"      = trimprefix(huaweicloud_obs_bucket.output.bucket_domain_name, "${huaweicloud_obs_bucket.output.bucket}.")
    "port"        = "443"
  }

  access_key = "%s"
  secret_key = "%s"
}

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
    retry_type                   = "NONE"
    scheduler_enabled            = false
    throttling_extractors_number = 4
    throttling_record_dirty_data = false
    throttling_max_error_records = 10
    throttling_loader_number     = 1
  }
}
`, clusterConfig, bucketName, bucketName, name, ak, sk, name)
}
