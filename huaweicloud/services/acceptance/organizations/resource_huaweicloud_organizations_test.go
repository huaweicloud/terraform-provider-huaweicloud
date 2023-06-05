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

func getOrganizationsResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getOrganizations: Query Organizations
	var (
		getOrganizationsHttpUrl = "v1/organizations"
		getOrganizationsProduct = "organizations"
	)
	getOrganizationsClient, err := cfg.NewServiceClient(getOrganizationsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationsPath := getOrganizationsClient.Endpoint + getOrganizationsHttpUrl

	getOrganizationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationsResp, err := getOrganizationsClient.Request("GET", getOrganizationsPath, &getOrganizationsOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations: %s", err)
	}
	return utils.FlattenResponse(getOrganizationsResp)
}

func TestAccOrganizations_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_organizations.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrganizationsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOrganizations_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "root_tags.key2", "value2"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "account_id"),
					resource.TestCheckResourceAttrSet(rName, "account_name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "root_id"),
					resource.TestCheckResourceAttrSet(rName, "root_name"),
					resource.TestCheckResourceAttrSet(rName, "root_urn"),
				),
			},
			{
				Config: testOrganizations_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "root_tags.key4", "value4"),
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

func testOrganizations_basic() string {
	return `
resource "huaweicloud_organizations" "test" {
  root_tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}
`
}

func testOrganizations_basic_update() string {
	return `
resource "huaweicloud_organizations" "test" {
  root_tags = {
    "key3" = "value3"
    "key4" = "value4"
  }
}
`
}
