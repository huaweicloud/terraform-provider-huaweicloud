package hss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCiCdConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
		id      = state.Primary.ID
		epsId   = state.Primary.Attributes["enterprise_project_id"]
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/image/cicd/configurations/{cicd_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cicd_id}", id)
	if epsId != "" {
		requestPath = fmt.Sprintf("%s?enterprise_project_id=%s", requestPath, epsId)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, the error information returned by the get API is inaccurate and cannot
		// determine whether the resource exists.
		// At this point, it is necessary to call the query list API to determine whether the resource exists.
		listRespBody, _ := hss.ListCiCdConfiguration(client, id, epsId)
		if listRespBody == nil {
			return nil, golangsdk.ErrDefault404{}
		}

		return nil, fmt.Errorf("error retrieving HSS CiCd configuration: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func TestAccCiCdConfiguration_basic(t *testing.T) {
	var (
		obj          interface{}
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_hss_cicd_configuration.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCiCdConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCiCdConfiguration_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cicd_name", name),
					resource.TestCheckResourceAttr(resourceName, "vulnerability_whitelist.0", "vulnerability_whitelist_test"),
					resource.TestCheckResourceAttr(resourceName, "vulnerability_blocklist.0", "vulnerability_blocklist_test"),
					resource.TestCheckResourceAttr(resourceName, "image_whitelist.0", "image_whitelist_test"),
				),
			},
			{
				Config: testCiCdConfiguration_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "vulnerability_whitelist.0", "vulnerability_whitelist_update"),
					resource.TestCheckResourceAttr(resourceName, "vulnerability_blocklist.0", "vulnerability_blocklist_update"),
					resource.TestCheckResourceAttr(resourceName, "image_whitelist.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
				},
			},
		},
	})
}

func testCiCdConfiguration_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_cicd_configuration" "test" {
  cicd_name               = "%[1]s"
  vulnerability_whitelist = ["vulnerability_whitelist_test"]
  vulnerability_blocklist = ["vulnerability_blocklist_test"]
  image_whitelist         = ["image_whitelist_test"]
  enterprise_project_id   = "0"
}
`, name)
}

func testCiCdConfiguration_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_cicd_configuration" "test" {
  cicd_name               = "%[1]s"
  vulnerability_whitelist = ["vulnerability_whitelist_update"]
  vulnerability_blocklist = ["vulnerability_blocklist_update"]
  image_whitelist         = []
  enterprise_project_id   = "0"
}
`, name)
}
