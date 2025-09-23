package live

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRecordingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/record/rules/{id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live recording: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccRecording_basic(t *testing.T) {
	var (
		randInt      = acctest.RandInt()
		rName        = "huaweicloud_live_recording.test"
		recordingObj interface{}
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&recordingObj,
		getRecordingFunc,
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
				Config: testRecording_basic(randInt),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "stream_name", "streamname"),
					resource.TestCheckResourceAttr(rName, "type", "CONTINUOUS_RECORD"),
					resource.TestCheckResourceAttr(rName, "obs.0.region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(rName, "hls.0.recording_length", "60"),
					resource.TestCheckResourceAttr(rName, "hls.0.record_slice_duration", "3"),
					resource.TestCheckResourceAttr(rName, "hls.0.max_stream_pause_length", "-1"),
					resource.TestCheckResourceAttrSet(rName, "hls.0.file_naming"),
					resource.TestCheckResourceAttrSet(rName, "hls.0.ts_file_naming"),
					resource.TestCheckResourceAttr(rName, "flv.0.recording_length", "80"),
					resource.TestCheckResourceAttr(rName, "flv.0.max_stream_pause_length", "0"),
					resource.TestCheckResourceAttrSet(rName, "flv.0.file_naming"),
					resource.TestCheckResourceAttr(rName, "mp4.0.recording_length", "100"),
					resource.TestCheckResourceAttr(rName, "mp4.0.max_stream_pause_length", "10"),
					resource.TestCheckResourceAttrSet(rName, "mp4.0.file_naming"),
				),
			},
			{
				Config: testRecording_update(randInt),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "hls.0.recording_length", "120"),
					resource.TestCheckResourceAttr(rName, "hls.0.record_slice_duration", "15"),
					resource.TestCheckResourceAttr(rName, "hls.0.max_stream_pause_length", "8"),
					resource.TestCheckResourceAttrSet(rName, "hls.0.file_naming"),
					resource.TestCheckResourceAttr(rName, "mp4.0.recording_length", "90"),
					resource.TestCheckResourceAttr(rName, "mp4.0.max_stream_pause_length", "0"),
					resource.TestCheckResourceAttrSet(rName, "mp4.0.file_naming"),
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

func testRecording_base(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "tf-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_live_bucket_authorization" "test" {
  bucket = huaweicloud_obs_bucket.test.bucket
}
`, randInt)
}

func testRecording_basic(randInt int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_recording" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test]

  domain_name = "%[2]s"
  app_name    = "live"
  stream_name = "streamname"
  type        = "CONTINUOUS_RECORD"

  obs {
    region = huaweicloud_obs_bucket.test.region
    bucket = huaweicloud_obs_bucket.test.bucket
  }

  hls {
    recording_length        = 60
    record_slice_duration   = 3
    max_stream_pause_length = -1
  }

  flv {
    recording_length = 80
  }

  mp4 {
    recording_length        = 100
    max_stream_pause_length = 10
  }
}
`, testRecording_base(randInt), acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testRecording_update(randInt int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_recording" "test" {
  depends_on = [huaweicloud_live_bucket_authorization.test]

  domain_name = "%[2]s"
  app_name    = "live"
  stream_name = "streamname"
  type        = "CONTINUOUS_RECORD"

  obs {
    region = huaweicloud_obs_bucket.test.region
    bucket = huaweicloud_obs_bucket.test.bucket
  }

  hls {
    recording_length        = 120
    record_slice_duration   = 15
    max_stream_pause_length = 8
  }

  mp4 {
    recording_length        = 90
    max_stream_pause_length = 0
  }
}
`, testRecording_base(randInt), acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}
