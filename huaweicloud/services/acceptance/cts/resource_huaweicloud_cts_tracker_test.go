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

func TestAccCTSTracker_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cts_tracker.tracker"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCTSTrackerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTracker_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", "system"),
					resource.TestCheckResourceAttr(resourceName, "type", "system"),
					resource.TestCheckResourceAttr(resourceName, "status", "enabled"),
				),
			},
			{
				Config: testAccCTSTracker_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "file_prefix", "cts-updated"),
					resource.TestCheckResourceAttr(resourceName, "lts_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCTSTrackerImportState(resourceName),
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

func testAccCheckCTSTrackerDestroy(s *terraform.State) error {
	// the system tracker always exists
	return nil
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

func testAccCTSTracker_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name = huaweicloud_obs_bucket.bucket.bucket
  file_prefix = "cts"
  lts_enabled = true
}
`, rName)
}

func testAccCTSTracker_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name = huaweicloud_obs_bucket.bucket.bucket
  file_prefix = "cts-updated"
  lts_enabled = false
  enabled     = false
}
`, rName)
}
