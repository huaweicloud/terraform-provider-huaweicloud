package deprecated

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccCTSTrackerV1DataSource_basic(t *testing.T) {
	dataName := "data.huaweicloud_cts_tracker.tracker"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheckDeprecated(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCTSTrackerV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCTSTrackerV1DataSourceID(dataName),
					resource.TestCheckResourceAttr(dataName, "bucket_name", "tf-test-bucket"),
					resource.TestCheckResourceAttr(dataName, "status", "enabled"),
				),
			},
		},
	})
}

func testAccCheckCTSTrackerV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find cts tracker data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("tracker data source not set ")
		}

		return nil
	}
}

const testAccCTSTrackerV1DataSource_basic = `
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "tf-test-topic"
  display_name = "The display name of tf-test-topic"
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

data "huaweicloud_cts_tracker" "tracker" {
  tracker_name = huaweicloud_cts_tracker.tracker.id
}
`
