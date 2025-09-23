package mpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mpc"
)

func getTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return nil, fmt.Errorf("error creating MPC client: %s", err)
	}

	return mpc.GetTranscodingTemplate(client, state.Primary.ID)
}

func TestAccTranscodingTemplate_basic(t *testing.T) {
	var (
		templateObj  interface{}
		rName        = acceptance.RandomAccResourceNameWithDash()
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_mpc_transcoding_template.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&templateObj,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "video.0.profile", "4"),
					resource.TestCheckResourceAttr(resourceName, "video.0.output_policy", "transcode"),
				),
			},
			{
				Config: testTranscodingTemplate_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "false"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.profile", "3"),
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

func testTranscodingTemplate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 1

  audio {
    bitrate       = 0
    channels      = 2
    codec         = 2
    output_policy = "transcode"
    sample_rate   = 1
  }

  video {
    max_consecutive_bframes = 7
    bitrate                 = 0
    black_bar_removal       = 0
    codec                   = 2
    fps                     = 0
    level                   = 15
    max_iframes_interval    = 5
    output_policy           = "transcode"
    quality                 = 1
    profile                 = 4
    height                  = 0
    width                   = 0
  }
}
`, rName)
}

func testTranscodingTemplate_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = false
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 2

  audio {
    bitrate       = 0
    channels      = 2
    codec         = 1
    output_policy = "transcode"
    sample_rate   = 1
  }

  video {
    max_consecutive_bframes = 7
    bitrate                 = 0
    black_bar_removal       = 0
    codec                   = 1
    fps                     = 0
    level                   = 15
    max_iframes_interval    = 5
    output_policy           = "transcode"
    quality                 = 1
    profile                 = 3
    height                  = 0
    width                   = 0
  }
}
`, rName)
}
