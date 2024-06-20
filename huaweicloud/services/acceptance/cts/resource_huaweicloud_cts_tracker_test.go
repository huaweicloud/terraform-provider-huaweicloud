package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCTSTracker_keepTracker(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cts_tracker.tracker"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTracker_keepTracker(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", "system"),
					resource.TestCheckResourceAttr(resourceName, "type", "system"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "compress_type", "gzip"),
					resource.TestCheckResourceAttr(resourceName, "is_sort_by_service", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(resourceName, "agency_name"),
				),
			},
			{
				Config: testAccCTSTracker_keepTracker_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts-updated"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "compress_type", "json"),
					resource.TestCheckResourceAttr(resourceName, "is_sort_by_service", "true"),
					resource.TestCheckResourceAttr(resourceName, "exclude_service.0", "KMS"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.newkey", "value"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCTSTrackerImportState(resourceName),
				ImportStateVerifyIgnore: []string{"tags"},
			},
		},
	})
}

func testAccCTSTrackerImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		tracker, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("CTS tracker not found")
		}

		name := tracker.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("CTS tracker resource %s not found", tracker.Primary.ID)
		}
		return name, nil
	}
}

func testAccCheckCTSTrackerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		ctsClient, err := cfg.HcCtsV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CTS client: %s", err)
		}

		name := rs.Primary.Attributes["name"]
		listOpts := &cts.ListTrackersRequest{
			TrackerName: &name,
		}

		response, err := ctsClient.ListTrackers(listOpts)
		if err != nil {
			return fmt.Errorf("error retrieving CTS tracker: %s", err)
		}

		if response.Trackers == nil || len(*response.Trackers) == 0 {
			return fmt.Errorf("can not find the CTS tracker %s", name)
		}

		return nil
	}
}

func testAccCheckCTSTrackerDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cts_tracker" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		listOpts := &cts.ListTrackersRequest{
			TrackerName: &name,
		}

		response, err := ctsClient.ListTrackers(listOpts)
		if err != nil {
			return fmt.Errorf("error retrieving CTS tracker: %s", err)
		}

		if response.Trackers == nil || len(*response.Trackers) == 0 {
			return fmt.Errorf("can not find the CTS tracker %s", name)
		}

		allTrackers := *response.Trackers
		ctsTracker := allTrackers[0]
		if ctsTracker.Status != nil && ctsTracker.Status.Value() != "disabled" {
			return fmt.Errorf("can not disable the CTS tracker %s", name)
		}
	}

	return nil
}

func getCTSTrackerResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcCtsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CTS client: %s", err)
	}

	name := state.Primary.Attributes["name"]
	trackerType := cts.GetListTrackersRequestTrackerTypeEnum().SYSTEM
	listOpts := &cts.ListTrackersRequest{
		TrackerName: &name,
		TrackerType: &trackerType,
	}

	response, err := client.ListTrackers(listOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CTS tracker: %s", err)
	}

	if response.Trackers == nil || len(*response.Trackers) == 0 {
		return nil, fmt.Errorf("can not find the CTS tracker %s", name)
	}

	allTrackers := *response.Trackers
	ctsTracker := allTrackers[0]

	return ctsTracker, nil
}

func TestAccCTSTracker_deleteTracker(t *testing.T) {
	var tracker cts.TrackerResponseBody
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cts_tracker.tracker"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&tracker,
		getCTSTrackerResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTracker_deleteTracker(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", "system"),
					resource.TestCheckResourceAttr(resourceName, "type", "system"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "compress_type", "json"),
					resource.TestCheckResourceAttr(resourceName, "is_sort_by_service", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "exclude_service.0", "KMS"),
					resource.TestCheckResourceAttrSet(resourceName, "agency_name"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCTSTrackerImportState(resourceName),
				ImportStateVerifyIgnore: []string{"tags", "delete_tracker"},
			},
		},
	})
}

func testAccCTSTracker_keepTracker(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name        = huaweicloud_obs_bucket.bucket.bucket
  file_prefix        = "cts"
  lts_enabled        = true
  compress_type      = "gzip"
  is_sort_by_service = false

  tags = {
    foo = "bar"
  }
}
`, rName)
}

func testAccCTSTracker_keepTracker_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name        = huaweicloud_obs_bucket.bucket.bucket
  file_prefix        = "cts-updated"
  lts_enabled        = false
  enabled            = false
  compress_type      = "json"
  is_sort_by_service = true
  exclude_service    = ["KMS"]

  tags = {
    foo    = "bar1"
    newkey = "value"
  }
}
`, rName)
}

func testAccCTSTracker_deleteTracker(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name        = huaweicloud_obs_bucket.bucket.bucket
  file_prefix        = "cts"
  lts_enabled        = true
  compress_type      = "json"
  is_sort_by_service = false
  delete_tracker     = true
  exclude_service    = ["KMS"]

  tags = {
    foo = "bar"
  }
}
`, rName)
}
