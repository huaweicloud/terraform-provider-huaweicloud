package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/live"
)

func getTranscodingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	getRespBody, err := live.GetTranscodingTemplates(client, state.Primary.Attributes["domain_name"],
		state.Primary.Attributes["app_name"])
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live transcoding: %s", err)
	}

	return getRespBody, nil
}

func TestAccTranscoding_basic(t *testing.T) {
	var (
		transcodingObj interface{}
		rName          = "huaweicloud_live_transcoding.test"
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&transcodingObj,
		getTranscodingFunc,
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
				Config: testAccTranscoding_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "video_encoding", "H264"),
					resource.TestCheckResourceAttr(rName, "trans_type", "publish"),
					resource.TestCheckResourceAttr(rName, "low_bitrate_hd", "false"),
					resource.TestCheckResourceAttr(rName, "templates.#", "1"),
					resource.TestCheckResourceAttr(rName, "templates.0.name", "t1"),
					resource.TestCheckResourceAttr(rName, "templates.0.width", "300"),
					resource.TestCheckResourceAttr(rName, "templates.0.height", "400"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate", "300"),
					resource.TestCheckResourceAttr(rName, "templates.0.frame_rate", "60"),
					resource.TestCheckResourceAttr(rName, "templates.0.protocol", "RTMP"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_interval", "500"),
					resource.TestCheckResourceAttr(rName, "templates.0.gop", "0"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate_adaptive", "adaptive"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_policy", "strictSync"),
				),
			},
			{
				Config: testAccTranscoding_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "video_encoding", "H265"),
					resource.TestCheckResourceAttr(rName, "trans_type", "play"),
					resource.TestCheckResourceAttr(rName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(rName, "templates.#", "2"),
					resource.TestCheckResourceAttr(rName, "templates.0.width", "480"),
					resource.TestCheckResourceAttr(rName, "templates.0.height", "520"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate", "600"),
					resource.TestCheckResourceAttr(rName, "templates.0.frame_rate", "0"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_interval", "0"),
					resource.TestCheckResourceAttr(rName, "templates.0.gop", "10"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate_adaptive", "minimum"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_policy", "auto"),
					resource.TestCheckResourceAttr(rName, "templates.1.name", "t2"),
					resource.TestCheckResourceAttr(rName, "templates.1.width", "600"),
					resource.TestCheckResourceAttr(rName, "templates.1.height", "800"),
					resource.TestCheckResourceAttr(rName, "templates.1.bitrate", "300"),
					resource.TestCheckResourceAttr(rName, "templates.1.frame_rate", "30"),
					resource.TestCheckResourceAttr(rName, "templates.1.protocol", "RTMP"),
					resource.TestCheckResourceAttr(rName, "templates.1.i_frame_interval", "50"),
					resource.TestCheckResourceAttr(rName, "templates.1.gop", "2"),
					resource.TestCheckResourceAttr(rName, "templates.1.bitrate_adaptive", "off"),
					resource.TestCheckResourceAttr(rName, "templates.1.i_frame_policy", "auto"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trans_type",
				},
			},
		},
	})
}

func testAccTranscoding_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_transcoding" "test" {
  domain_name    = "%s"
  app_name       = "live"
  video_encoding = "H264"
  trans_type     = "publish"

  templates {
    name             = "t1"
    width            = 300
    height           = 400
    bitrate          = 300
    frame_rate       = 60
    protocol         = "RTMP"
    i_frame_interval = "500"
    gop              = "0"
    bitrate_adaptive = "adaptive"
    i_frame_policy   = "strictSync"
  }
}
`, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}

func testAccTranscoding_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_transcoding" "test" {
  domain_name    = "%s"
  app_name       = "live"
  video_encoding = "H265"
  trans_type     = "play"
  low_bitrate_hd = true

  templates {
    name             = "t1"
    width            = 480
    height           = 520
    bitrate          = 600
    frame_rate       = 0
    protocol         = "RTMP"
    i_frame_interval = "0"
    gop              = "10"
    bitrate_adaptive = "minimum"
    i_frame_policy   = "auto"
  }

  templates {
    name             = "t2"
    width            = 600
    height           = 800
    bitrate          = 300
    frame_rate       = 30
    protocol         = "RTMP"
    bitrate_adaptive = "off"
    i_frame_policy   = "auto"
  }
}
`, acceptance.HW_LIVE_INGEST_DOMAIN_NAME)
}
