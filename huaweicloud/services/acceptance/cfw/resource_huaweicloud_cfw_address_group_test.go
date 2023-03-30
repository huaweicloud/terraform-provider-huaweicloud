package cfw

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

func getAddressGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAddressGroup: Query the CFW IP address group detail
	var (
		getAddressGroupHttpUrl = "v1/{project_id}/address-sets/{id}"
		getAddressGroupProduct = "cfw"
	)
	getAddressGroupClient, err := cfg.NewServiceClient(getAddressGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	getAddressGroupPath := getAddressGroupClient.Endpoint + getAddressGroupHttpUrl
	getAddressGroupPath = strings.ReplaceAll(getAddressGroupPath, "{project_id}", getAddressGroupClient.ProjectID)
	getAddressGroupPath = strings.ReplaceAll(getAddressGroupPath, "{id}", state.Primary.ID)

	getAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAddressGroupResp, err := getAddressGroupClient.Request("GET", getAddressGroupPath, &getAddressGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AddressGroup: %s", err)
	}
	return utils.FlattenResponse(getAddressGroupResp)
}

func TestAccAddressGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_address_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAddressGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
				),
			},
			{
				Config: testAddressGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"object_id"},
			},
		},
	})
}

func testAddressGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_address_group" "test" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%s"
  description = "terraform test"
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testAddressGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_address_group" "test" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%s-update"
  description = "terraform test update"
}
`, testAccDatasourceFirewalls_basic(), name)
}
