package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getTranscodingResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcLiveV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Live v1 client: %s", err)
	}

	var (
		domain  = state.Primary.Attributes["domain_name"]
		appName = state.Primary.Attributes["app_name"]
		getOpts = &model.ShowTranscodingsTemplateRequest{
			Domain:  domain,
			AppName: &appName,
		}
	)
	resp, err := client.ShowTranscodingsTemplate(getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live transcoding: %s", err)
	}

	if resp.Templates == nil || len(*resp.Templates) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, nil
}

func TestAccTranscoding_basic(t *testing.T) {
	var (
		obj            model.ShowTranscodingsTemplateResponse
		pushDomainName = fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
		rName          = "huaweicloud_live_transcoding.test"
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTranscodingResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscoding_basic(pushDomainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", pushDomainName),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "video_encoding", "H264"),
					resource.TestCheckResourceAttr(rName, "trans_type", "publish"),
					resource.TestCheckResourceAttr(rName, "low_bitrate_hd", "false"),
					resource.TestCheckResourceAttr(rName, "templates.#", "1"),
					resource.TestCheckResourceAttr(rName, "templates.0.name", "t1"),
					resource.TestCheckResourceAttr(rName, "templates.0.width", "300"),
					resource.TestCheckResourceAttr(rName, "templates.0.height", "400"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate", "300"),
					resource.TestCheckResourceAttr(rName, "templates.0.protocol", "RTMP"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_interval", "500"),
					resource.TestCheckResourceAttr(rName, "templates.0.gop", "0"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate_adaptive", "adaptive"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_policy", "strictSync"),
				),
			},
			{
				Config: testTranscoding_update(pushDomainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "domain_name", pushDomainName),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "video_encoding", "H264"),
					resource.TestCheckResourceAttr(rName, "trans_type", "play"),
					resource.TestCheckResourceAttr(rName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(rName, "templates.#", "2"),
					resource.TestCheckResourceAttr(rName, "templates.0.protocol", "RTMP"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_interval", "0"),
					resource.TestCheckResourceAttr(rName, "templates.0.gop", "10"),
					resource.TestCheckResourceAttr(rName, "templates.0.bitrate_adaptive", "minimum"),
					resource.TestCheckResourceAttr(rName, "templates.0.i_frame_policy", "auto"),
					resource.TestCheckResourceAttr(rName, "templates.1.name", "t2"),
					resource.TestCheckResourceAttr(rName, "templates.1.width", "600"),
					resource.TestCheckResourceAttr(rName, "templates.1.height", "800"),
					resource.TestCheckResourceAttr(rName, "templates.1.bitrate", "300"),
					resource.TestCheckResourceAttr(rName, "templates.1.protocol", "RTMP"),
					resource.TestCheckResourceAttr(rName, "templates.1.gop", "10"),
					resource.TestCheckResourceAttr(rName, "templates.1.bitrate_adaptive", "off"),
					resource.TestCheckResourceAttr(rName, "templates.1.i_frame_policy", "auto"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s/live", pushDomainName),
				ImportStateVerifyIgnore: []string{
					"trans_type",
				},
			},
		},
	})
}

func testTranscoding_basic(pushDomainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name = "%s"
  type = "push"
}

resource "huaweicloud_live_transcoding" "test" {
  domain_name    = huaweicloud_live_domain.ingestDomain.name
  app_name       = "live"
  video_encoding = "H264"
  trans_type     = "publish"

  templates {
    name             = "t1"
    width            = 300
    height           = 400
    bitrate          = 300
    protocol         = "RTMP"
    i_frame_interval = 500
    gop              = 0
    bitrate_adaptive = "adaptive"
    i_frame_policy   = "strictSync"
  }
}
`, pushDomainName)
}

func testTranscoding_update(pushDomainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_live_domain" "ingestDomain" {
  name = "%s"
  type = "push"
}

resource "huaweicloud_live_transcoding" "test" {
  domain_name    = huaweicloud_live_domain.ingestDomain.name
  app_name       = "live"
  video_encoding = "H264"
  trans_type     = "play"
  low_bitrate_hd = true

  templates {
    name             = "t1"
    width            = 300
    height           = 400
    bitrate          = 300
    protocol         = "RTMP"
    i_frame_interval = 0
    gop              = 10
    bitrate_adaptive = "minimum"
    i_frame_policy   = "auto"
  }

  templates {
    name             = "t2"
    width            = 600
    height           = 800
    bitrate          = 300
    protocol         = "RTMP"
    gop              = 10
    bitrate_adaptive = "off"
    i_frame_policy   = "auto"
  }
}
`, pushDomainName)
}
