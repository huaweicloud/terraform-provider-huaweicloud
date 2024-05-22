package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRAMOrganizationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/organization-share"
		product = "ram"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RAM client: %s", err)
	}

	getOrganizationPath := client.Endpoint + httpUrl
	getOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOrganizationResp, err := client.Request("GET", getOrganizationPath, &getOrganizationOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RAM organization: %s", err)
	}

	return utils.FlattenResponse(getOrganizationResp)
}

// Only organization administrators have permission to run this test case
func TestAccRAMOrganization_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ram_organization.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRAMOrganizationResourceFunc,
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testRAMOrganization_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
				),
			},
			{
				Config: testRAMOrganization_basic_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
				),
			},
		},
	})
}

const testRAMOrganization_basic = `
resource "huaweicloud_ram_organization" "test" {
  enabled = true
}
`

const testRAMOrganization_basic_update = `
resource "huaweicloud_ram_organization" "test" {
  enabled = false
}
`
