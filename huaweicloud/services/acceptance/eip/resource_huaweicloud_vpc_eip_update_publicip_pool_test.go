package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eip"
)

func getEipUpdatePublicipPoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("vpc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC client: %s", err)
	}

	return eip.ReadEipUpdatePublicipPool(client, state.Primary.ID)
}

func TestAccEipUpdatePublicipPool_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpc_eip_update_publicip_pool.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEipUpdatePublicipPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVpcEipPoolId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEipUpdatePublicipPool_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "publicip_pool_id", acceptance.HW_VPC_EIP_PUBLICIP_POOL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttrSet(rName, "size"),
					resource.TestCheckResourceAttrSet(rName, "used"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "public_border_group"),
					resource.TestCheckResourceAttrSet(rName, "shared"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
				),
			},
			{
				Config: testAccEipUpdatePublicipPool_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "publicip_pool_id", acceptance.HW_VPC_EIP_PUBLICIP_POOL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				Config: testAccEipUpdatePublicipPool_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "publicip_pool_id", acceptance.HW_VPC_EIP_PUBLICIP_POOL_ID),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
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

func testAccEipUpdatePublicipPool_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip_update_publicip_pool" "test" {
  publicip_pool_id = "%[1]s"
  name             = "%[2]s"
  description      = "test description"
}
`, acceptance.HW_VPC_EIP_PUBLICIP_POOL_ID, name)
}

func testAccEipUpdatePublicipPool_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip_update_publicip_pool" "test" {
  publicip_pool_id = "%[1]s"
  name             = "%[2]s"
  description      = "test description update"
}
`, acceptance.HW_VPC_EIP_PUBLICIP_POOL_ID, name)
}

func testAccEipUpdatePublicipPool_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip_update_publicip_pool" "test" {
  publicip_pool_id = "%[1]s"
  name             = "%[2]s_update"
  description      = "test description update"
}
`, acceptance.HW_VPC_EIP_PUBLICIP_POOL_ID, name)
}
