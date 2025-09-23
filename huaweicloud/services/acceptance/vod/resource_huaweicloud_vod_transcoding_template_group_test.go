package vod

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "vod"
		httpUrl = "v1.0/{project_id}/asset/template_group/transcodings"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?group_id=%s", state.Primary.ID)
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	templateGroup := utils.PathSearch("template_group_list|[0]", respBody, nil)
	if templateGroup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return templateGroup, nil
}

func TestAccTranscodingTemplateGroup_basic(t *testing.T) {
	var template interface{}
	rName := acceptance.RandomAccResourceName()
	rNameUpdate := rName + "_update"
	resourceName := "huaweicloud_vod_transcoding_template_group.test"

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
				Config: testTranscodingTemplateGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test group"),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "false"),
					resource.TestCheckResourceAttr(resourceName, "audio_codec", "HEAAC1"),
					resource.TestCheckResourceAttr(resourceName, "video_codec", "H265"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.output_format", "MP4"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.audio.0.channels", "1"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.audio.0.sample_rate", "2"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.audio.0.bitrate", "0"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.video.0.height", "1080"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.video.0.width", "1920"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.video.0.quality", "FHD"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testTranscodingTemplateGroup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "test group update"),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "audio_codec", "AAC"),
					resource.TestCheckResourceAttr(resourceName, "video_codec", "H264"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.output_format", "HLS"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.audio.0.channels", "2"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.audio.0.bitrate", "8"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.video.0.height", "720"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.video.0.width", "1280"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.0.video.0.quality", "HD"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.1.video.0.height", "1080"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.1.video.0.width", "1920"),
					resource.TestCheckResourceAttr(resourceName, "quality_info.1.video.0.quality", "FHD"),
				),
			},
		},
	})
}

func testTranscodingTemplateGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_transcoding_template_group" "test" {
  name        = "%s"
  description = "test group"
  audio_codec = "HEAAC1"
  video_codec = "H265"

  quality_info {
    output_format = "MP4"

    audio {
      channels    = 1
      sample_rate = 2
    }

    video {
      bitrate    = 1000
      frame_rate = 1
      height     = 1080
      quality    = "FHD"
      width      = 1920
    }
  }
}
`, rName)
}

func testTranscodingTemplateGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_transcoding_template_group" "test" {
  name                 = "%s"
  description          = "test group update"
  audio_codec          = "AAC"
  hls_segment_duration = 5
  low_bitrate_hd       = true
  video_codec          = "H264"

  quality_info {
    output_format = "HLS"

    audio {
      bitrate     = 8
      channels    = 2
      sample_rate = 2
    }

    video {
      bitrate    = 1000
      frame_rate = 1
      height     = 720
      quality    = "HD"
      width      = 1280
    }
  }

  quality_info {
    output_format = "HLS"

    audio {
      channels    = 2
      sample_rate = 2
    }

    video {
      bitrate    = 1000
      frame_rate = 1
      height     = 1080
      quality    = "FHD"
      width      = 1920
    }
  }
}
`, rName)
}
