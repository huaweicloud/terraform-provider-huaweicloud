package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRecordingResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcLiveV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Live v1 client: %s", err)
	}
	return client.ShowRecordRule(&model.ShowRecordRuleRequest{Id: state.Primary.ID})
}

func TestAccRecording_basic(t *testing.T) {
	var obj model.ShowRecordRuleRequest

	rName := "huaweicloud_live_recording.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRecordingResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveIngestDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRecording_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "stream_name", "streamname"),
					resource.TestCheckResourceAttr(rName, "type", "CONTINUOUS_RECORD"),
					resource.TestCheckResourceAttr(rName, "obs.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "mp4.0.recording_length", "120"),
					resource.TestCheckResourceAttrSet(rName, "mp4.0.file_naming"),
				),
			},
			{
				Config: testRecording_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "stream_name", "streamname"),
					resource.TestCheckResourceAttr(rName, "type", "CONTINUOUS_RECORD"),
					resource.TestCheckResourceAttr(rName, "obs.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "hls.0.recording_length", "120"),
					resource.TestCheckResourceAttr(rName, "hls.0.record_slice_duration", "5"),
					resource.TestCheckResourceAttrSet(rName, "hls.0.file_naming"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testRecording_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true

  lifecycle {
    ignore_changes = [
      cors_rule,
    ]
  }
}

resource "huaweicloud_live_bucket_authorization" "test" {
  depends_on = [huaweicloud_obs_bucket.bucket]

  bucket = "%[1]s"
}
`, acceptance.RandomAccResourceNameWithDash())
}

func testRecording_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_recording" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test]

  domain_name = "%[2]s"
  app_name    = "live"
  stream_name = "streamname"
  type        = "CONTINUOUS_RECORD"

  obs {
    region = huaweicloud_obs_bucket.bucket.region
    bucket = huaweicloud_obs_bucket.bucket.bucket
  }

  mp4 {
    recording_length = 120
  }
}
`, testRecording_base(), acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testRecording_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_recording" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test]

  domain_name = "%[2]s"
  app_name    = "live"
  stream_name = "streamname"
  type        = "CONTINUOUS_RECORD"

  obs {
    region = huaweicloud_obs_bucket.bucket.region
    bucket = huaweicloud_obs_bucket.bucket.bucket
  }

  hls {
    recording_length      = 120
    record_slice_duration = 5
  }
}
`, testRecording_base(), acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}
