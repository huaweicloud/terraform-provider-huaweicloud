package vod

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceWatermarkTemplate(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcVodV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VOD client: %s", err)
	}

	resp, err := client.ListWatermarkTemplate(&vod.ListWatermarkTemplateRequest{Id: &[]string{state.Primary.ID}})
	if err != nil {
		return nil, fmt.Errorf("error retrieving VOD watermark template: %s", err)
	}

	if resp.Templates == nil || len(*resp.Templates) == 0 {
		return nil, fmt.Errorf("unable to retrieve VOD watermark template: %s", state.Primary.ID)
	}

	templateList := *resp.Templates
	return templateList[0], nil
}

func TestAccWatermarkTemplate_basic(t *testing.T) {
	var template vod.WatermarkTemplate
	rName := acceptance.RandomAccResourceName()
	updateName := rName + "-update"
	resourceName := "huaweicloud_vod_watermark_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getResourceWatermarkTemplate,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckVODWatermark(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWatermarkTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "image_type", "PNG"),
					resource.TestCheckResourceAttr(resourceName, "position", "TOPRIGHT"),
					resource.TestCheckResourceAttr(resourceName, "image_process", "TRANSPARENT"),
					resource.TestCheckResourceAttr(resourceName, "horizontal_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "vertical_offset", "0"),
					resource.TestCheckResourceAttr(resourceName, "width", "0.01"),
					resource.TestCheckResourceAttr(resourceName, "height", "0.01"),
				),
			},
			{
				Config: testAccWatermarkTemplate_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "image_type", "PNG"),
					resource.TestCheckResourceAttr(resourceName, "position", "TOPLEFT"),
					resource.TestCheckResourceAttr(resourceName, "image_process", "ORIGINAL"),
					resource.TestCheckResourceAttr(resourceName, "horizontal_offset", "0.05"),
					resource.TestCheckResourceAttr(resourceName, "vertical_offset", "0.05"),
					resource.TestCheckResourceAttr(resourceName, "width", "0.1"),
					resource.TestCheckResourceAttr(resourceName, "height", "0.1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_file",
				},
			},
		},
	})
}

func testAccWatermarkTemplate_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_watermark_template" "test" {
  name       = "%s"
  image_file = "%s"
  image_type = "PNG"
}
`, rName, acceptance.HW_VOD_WATERMARK_FILE)
}

func testAccWatermarkTemplate_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vod_watermark_template" "test" {
  name              = "%s"
  image_file        = "%s"
  image_type        = "PNG"
  position          = "TOPLEFT"
  image_process     = "ORIGINAL"
  horizontal_offset = "0.05"
  vertical_offset   = "0.05"
  width             = "0.1"
  height            = "0.1"
}
`, rName, acceptance.HW_VOD_WATERMARK_FILE)
}
