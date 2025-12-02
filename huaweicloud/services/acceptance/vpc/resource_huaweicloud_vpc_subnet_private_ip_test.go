package vpc

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

func getSubnetPrivateIPResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpc", region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
	}

	getSubnetPrivateIPHttpUrl := "v1/{project_id}/privateips/{privateip_id}"
	getSubnetPrivateIPPath := client.Endpoint + getSubnetPrivateIPHttpUrl
	getSubnetPrivateIPPath = strings.ReplaceAll(getSubnetPrivateIPPath, "{project_id}", client.ProjectID)
	getSubnetPrivateIPPath = strings.ReplaceAll(getSubnetPrivateIPPath, "{privateip_id}", state.Primary.ID)
	getSubnetPrivateIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getSubnetPrivateIPResp, err := client.Request("GET", getSubnetPrivateIPPath, &getSubnetPrivateIPOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VPC subnet private IP: %s", err)
	}

	return utils.FlattenResponse(getSubnetPrivateIPResp)
}

func TestAccSubnetPrivateIP_basic(t *testing.T) {
	var (
		privateIP    interface{}
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_vpc_subnet_private_ip.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&privateIP,
			getSubnetPrivateIPResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetPrivateIP_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.0.111"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "DOWN"),
					resource.TestCheckResourceAttr(resourceName, "device_owner", "neutron:VIP_PORT"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSubnetPrivateIP_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet_private_ip" "test" {
  subnet_id    = huaweicloud_vpc_subnet.test.id
  ip_address   = "192.168.0.111"
  device_owner = "neutron:VIP_PORT"
}
`, testAccVpcSubnetV1_basic(name))
}
