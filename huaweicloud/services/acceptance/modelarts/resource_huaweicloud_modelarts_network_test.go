package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getNetworkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("modelarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetNetworkById(client, state.Primary.ID)
}

func TestAccNetwork_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelarts_network.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getNetworkResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetwork_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "10.168.0.0/16"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "0"),
				),
			},
			{
				Config: testAccNetwork_basic_step2(name), // add a connection
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "10.168.0.0/16"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testAccNetwork_basic_step3(name), // remove a connection
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "10.168.0.0/16"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "0"),
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

func testAccNetwork_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_network" "test" {
  # Network deletion will delete the peer connection, so the subnet can be deleted after the network is deleted.
  depends_on = [huaweicloud_vpc_subnet.test]

  name = "%[2]s"
  cidr = "10.168.0.0/16" # The recommended connecting CIDR about SFS Turbo.
}
`, common.TestVpc(name), name)
}

func testAccNetwork_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_network" "test" {
  depends_on = [huaweicloud_vpc_subnet.test]

  name = "%[2]s"
  cidr = "10.168.0.0/16"

  peer_connections {
    vpc_id    = huaweicloud_vpc.test.id
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestVpc(name), name)
}

func testAccNetwork_basic_step3(name string) string {
	return testAccNetwork_basic_step1(name)
}
