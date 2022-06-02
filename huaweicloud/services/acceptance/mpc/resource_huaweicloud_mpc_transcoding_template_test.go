package mpc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcMpcV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MPC client: %s", err)
	}

	id, err := strconv.ParseInt(state.Primary.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	resp, err := c.ListTemplate(&mpc.ListTemplateRequest{TemplateId: &[]int32{int32(id)}})
	if err != nil {
		return nil, fmt.Errorf("error retrieving MPC transcoding template: %d", err)
	}

	templateList := *resp.TemplateArray
	template := templateList[0].Template
	if template == nil {
		return nil, fmt.Errorf("unable to retrieve MPC transcoding template: %d", id)
	}

	return template, nil
}

func TestAccTranscodingTemplate_basic(t *testing.T) {
	var template mpc.QueryTransTemplate
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "huaweicloud_mpc_transcoding_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
    max_reference_frames    = 4
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
    max_reference_frames    = 4
    height                  = 0
    width                   = 0
  }
}
`, rName)
}
