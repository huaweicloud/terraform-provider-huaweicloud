package organizations

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

func resourceOrganizationRead(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	// getOrganization: Query Organizations organization
	var (
		region                 = acceptance.HW_REGION_NAME
		getOrganizationHttpUrl = "v1/organizations"
		getOrganizationProduct = "organizations"
	)
	getOrganizationClient, err := cfg.NewServiceClient(getOrganizationProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationPath := getOrganizationClient.Endpoint + getOrganizationHttpUrl

	getOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationResp, err := getOrganizationClient.Request("GET", getOrganizationPath, &getOrganizationOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations organization: %s", err)
	}
	return utils.FlattenResponse(getOrganizationResp)
}

func TestAccOrganization_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_organizations_organization.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		resourceOrganizationRead,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOrganization_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "root_tags.key2", "value2"),
					resource.TestCheckResourceAttr(rName, "enabled_policy_types.0", "service_control_policy"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "master_account_id"),
					resource.TestCheckResourceAttrSet(rName, "master_account_name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "root_id"),
					resource.TestCheckResourceAttrSet(rName, "root_name"),
					resource.TestCheckResourceAttrSet(rName, "root_urn"),
				),
			},
			{
				Config: testOrganization_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "root_tags.key4", "value4"),
					resource.TestCheckResourceAttr(rName, "enabled_policy_types.0", "tag_policy"),
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

func testOrganization_basic() string {
	return `
resource "huaweicloud_organizations_organization" "test" {
  enabled_policy_types = ["service_control_policy"]

  root_tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}
`
}

func testOrganization_basic_update() string {
	return `
resource "huaweicloud_organizations_organization" "test" {
  enabled_policy_types = ["tag_policy"]

  root_tags = {
    "key3" = "value3"
    "key4" = "value4"
  }
}
`
}
