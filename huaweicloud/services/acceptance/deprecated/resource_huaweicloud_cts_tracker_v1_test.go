package deprecated

import (
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cts/v1/tracker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccCTSTrackerV1_basic(t *testing.T) {
	var tracker tracker.Tracker
	resourceName := "huaweicloud_cts_tracker.tracker"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckCTSTrackerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(resourceName, &tracker),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "tf-test-bucket"),
					resource.TestCheckResourceAttr(resourceName, "file_prefix_name", "yO8Q"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCTSTrackerV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1Exists(resourceName, &tracker),
					resource.TestCheckResourceAttr(resourceName, "file_prefix_name", "yO8Q1"),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	ctsClient, err := config.CtsV1Client(config.GetRegion(nil))
	if err != nil {
		return fmtp.Errorf("Error creating cts client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cts_tracker" {
			continue
		}

		_, err := tracker.List(ctsClient, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return fmtp.Errorf("cts tracker still exists")
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}

	return nil
}

func testAccCheckCTSTrackerV1Exists(n string, trackers *tracker.Tracker) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		ctsClient, err := config.CtsV1Client(config.GetRegion(nil))
		if err != nil {
			return fmtp.Errorf("Error creating cts client: %s", err)
		}

		trackerList, err := tracker.List(ctsClient, tracker.ListOpts{TrackerName: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := trackerList[0]
		if found.TrackerName != rs.Primary.ID {
			return fmtp.Errorf("cts tracker not found")
		}

		*trackers = found

		return nil
	}
}

const testAccCTSTrackerV1_basic = `
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "topic_check"
  display_name = "The display name of topic_check"
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name               = huaweicloud_obs_bucket.bucket.bucket
  file_prefix_name          = "yO8Q"
  is_support_smn            = true
  topic_id                  = huaweicloud_smn_topic.topic_1.id
  is_send_all_key_operation = false
  operations                = ["login"]
  need_notify_user_list     = ["user1"]
}
`

const testAccCTSTrackerV1_update = `
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket"
  acl           = "public-read"
  force_destroy = true
}
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "topic_check1"
  display_name = "The display name of topic_check"
}
resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name               = huaweicloud_obs_bucket.bucket.bucket
  file_prefix_name          = "yO8Q1"
  is_support_smn            = true
  topic_id                  = huaweicloud_smn_topic.topic_1.id
  is_send_all_key_operation = false
  operations                = ["login"]
  need_notify_user_list     = ["user1"]
}
`
