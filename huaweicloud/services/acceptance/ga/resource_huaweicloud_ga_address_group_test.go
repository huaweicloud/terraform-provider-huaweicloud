package ga

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

func getIpAddressGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	// getIpAddressGroup: Query the GA IP address group detail
	var (
		getIpAddressGroupHttpUrl = "v1/ip-groups/{ip_group_id}"
		getIpAddressGroupProduct = "ga"
	)

	getIpAddressGroupClient, err := cfg.NewServiceClient(getIpAddressGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IP address group client: %s", err)
	}

	getIpAddressGroupPath := getIpAddressGroupClient.Endpoint + getIpAddressGroupHttpUrl
	getIpAddressGroupPath = strings.ReplaceAll(getIpAddressGroupPath, "{ip_group_id}", state.Primary.ID)

	getIpAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getIpAddressGroupResp, err := getIpAddressGroupClient.Request("GET", getIpAddressGroupPath, &getIpAddressGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IP address group: %s", err)
	}
	return utils.FlattenResponse(getIpAddressGroupResp)
}

func TestAccIpAddressGroup_basic(t *testing.T) {
	var (
		obj         interface{}
		name        = acceptance.RandomAccResourceNameWithDash()
		updateName  = acceptance.RandomAccResourceNameWithDash()
		description = "Created by terraform"
		rName       = "huaweicloud_ga_address_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIpAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIpAddressGroup_basic(name, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0.cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0.description",
						"The IP addresses included in the address group"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testIpAddressGroup_basic(updateName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
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

func testIpAddressGroup_basic(name, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_address_group" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  ip_addresses {
    cidr        = "192.168.1.0/24"
    description = "The IP addresses included in the address group"
  }
}
`, name, description)
}
