package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getInstanceRoutesFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return apig.ListInstanceRoutes(client, state.Primary.ID)
}

func TestAccInstanceRoutes_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_apig_instance_routes.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceRoutesFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceRoutes_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "nexthops.#", "2"),
				),
			},
			{
				Config: testAccInstanceRoutes_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "nexthops.#", "2"),
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

func testAccInstanceRoutes_base() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccInstanceRoutes_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_instance_routes" "test" {
  instance_id = local.instance_id
  nexthops    = ["172.16.128.0/20","172.16.0.0/20"]
}
`, testAccInstanceRoutes_base())
}

func testAccInstanceRoutes_basic_step2() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_instance_routes" "test" {
  instance_id = local.instance_id
  nexthops    = ["172.16.64.0/20","172.16.192.0/20"]
}
`, testAccInstanceRoutes_base())
}
