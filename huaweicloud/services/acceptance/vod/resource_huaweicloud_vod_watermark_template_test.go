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

func getResourceWatermarkTemplate(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "vod"
		httpUrl = "v1.0/{project_id}/template/watermark"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?id=%s", state.Primary.ID)
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VOD watermark template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	template := utils.PathSearch("templates|[0]", respBody, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return template, nil
}

func TestAccWatermarkTemplate_basic(t *testing.T) {
	var template interface{}
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
