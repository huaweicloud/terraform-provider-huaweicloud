package cts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCTSDataTrackerResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CTS client: %s", err)
	}

	getTrackerHttpUrl := "v3/{project_id}/trackers?tracker_name={name}&tracker_type={tracker_type}"
	name := state.Primary.Attributes["name"]
	trackerType := "data"
	getDataTrackerOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTrackerPath := client.Endpoint + getTrackerHttpUrl
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{project_id}", client.ProjectID)
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{name}", name)
	getTrackerPath = strings.ReplaceAll(getTrackerPath, "{tracker_type}", trackerType)
	response, err := client.Request("GET", getTrackerPath, &getDataTrackerOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the CTS data tracker: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, err
	}

	tracker := utils.PathSearch("trackers|[0]", respBody, nil)
	if tracker == nil {
		return nil, fmt.Errorf("error retrieving the CTS data tracker %s", name)
	}

	return tracker, nil
}

func TestAccCTSDataTracker_basic(t *testing.T) {
	var dataTracker interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cts_data_tracker.tracker"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&dataTracker,
		getCTSDataTrackerResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCTSDataTracker_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "transfer_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "type", "data"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "data_operation.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(resourceName, "data_bucket",
						"huaweicloud_obs_bucket.data_bucket", "bucket"),
					resource.TestCheckResourceAttrSet(resourceName, "agency_name"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "stream_id"),
					resource.TestCheckResourceAttrSet(resourceName, "log_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "log_topic_name"),
				),
			},
			{
				Config: testAccCTSDataTracker_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "transfer_enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket_name",
						"huaweicloud_obs_bucket.trans_bucket", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts"),
					resource.TestCheckResourceAttr(resourceName, "obs_retention_period", "30"),
					resource.TestCheckResourceAttr(resourceName, "validate_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.newkey", "value"),
					resource.TestCheckResourceAttr(resourceName, "compress_type", "json"),
					resource.TestCheckResourceAttr(resourceName, "is_sort_by_service", "false"),
				),
			},
			{
				Config: testAccCTSDataTracker_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "transfer_enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "bucket_name",
						"huaweicloud_obs_bucket.trans_bucket", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts"),
					resource.TestCheckResourceAttr(resourceName, "obs_retention_period", "30"),
					resource.TestCheckResourceAttr(resourceName, "validate_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.newkey", "value"),
					resource.TestCheckResourceAttr(resourceName, "compress_type", "json"),
					resource.TestCheckResourceAttr(resourceName, "is_sort_by_service", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testCTSDataTrackerImportState(resourceName),
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccCTSDataTracker_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "data_bucket" {
  bucket = "%[1]s-data"
  acl    = "public-read"
}

resource "huaweicloud_cts_data_tracker" "tracker" {
  name        = "%[1]s"
  data_bucket = huaweicloud_obs_bucket.data_bucket.bucket
  lts_enabled = true
  
  tags = {
    foo = "bar"
  }
}
`, rName)
}

func testAccCTSDataTracker_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "data_bucket" {
  bucket = "%[1]s-data"
  acl    = "public-read"
}

resource "huaweicloud_obs_bucket" "trans_bucket" {
  bucket        = "%[1]s-log"
  acl           = "private"
  force_destroy = true

  lifecycle {
    ignore_changes = [lifecycle_rule]
  }
}

resource "huaweicloud_cts_data_tracker" "tracker" {
  name                 = "%[1]s"
  data_bucket          = huaweicloud_obs_bucket.data_bucket.bucket
  bucket_name          = huaweicloud_obs_bucket.trans_bucket.bucket
  obs_retention_period = 30
  file_prefix          = "cts"
  validate_file        = false
  lts_enabled          = false
  compress_type        = "json"
  is_sort_by_service   = false

  tags = {
    foo    = "bar1"
    newkey = "value"
  }
}
`, rName)
}

func testAccCTSDataTracker_disable(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "data_bucket" {
  bucket = "%[1]s-data"
  acl    = "public-read"
}

resource "huaweicloud_obs_bucket" "trans_bucket" {
  bucket        = "%[1]s-log"
  acl           = "private"
  force_destroy = true

  lifecycle {
    ignore_changes = [lifecycle_rule]
  }
}

resource "huaweicloud_cts_data_tracker" "tracker" {
  name                 = "%[1]s"
  data_bucket          = huaweicloud_obs_bucket.data_bucket.bucket
  bucket_name          = huaweicloud_obs_bucket.trans_bucket.bucket
  obs_retention_period = 30
  file_prefix          = "cts"
  validate_file        = false
  lts_enabled          = false
  compress_type        = "json"
  is_sort_by_service   = false
  enabled              = false

  tags = {
    foo    = "bar1"
    newkey = "value"
  }
}
`, rName)
}

func testCTSDataTrackerImportState(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		return rs.Primary.Attributes["name"], nil
	}
}
