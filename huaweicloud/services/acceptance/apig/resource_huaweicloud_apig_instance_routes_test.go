package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getInstanceRoutesFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := instances.ListFeaturesOpts{
		// Default value of parameter 'limit' is 20, parameter 'offset' is an invalid parameter.
		// If we omit it, we can only obtain 20 features, other features will be lost.
		Limit: 500,
	}
	resp, err := instances.ListFeatures(client, state.Primary.ID, opts)
	if err != nil {
		return nil, fmt.Errorf("error querying feature list: %s", err)
	}

	for _, val := range resp {
		if val.Name == "route" {
			return val, nil
		}
	}
	return nil, fmt.Errorf("error querying feature: route")
}

func TestAccInstanceRoutes_basic(t *testing.T) {
	var (
		feature instances.Feature

		rName = "huaweicloud_apig_instance_routes.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&feature,
		getInstanceRoutesFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceRoutes_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "nexthops.#", "2"),
				),
			},
			{
				Config: testAccInstanceRoutes_basic_step2(name),
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

func testAccInstanceRoutes_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "0"
  maintain_begin        = "14:00:00"
}
`, common.TestBaseNetwork(name), name)
}

func testAccInstanceRoutes_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_instance_routes" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  nexthops    = ["172.16.128.0/20","172.16.0.0/20"]
}
`, testAccInstanceRoutes_base(name))
}

func testAccInstanceRoutes_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_instance_routes" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  nexthops    = ["172.16.64.0/20","172.16.192.0/20"]
}
`, testAccInstanceRoutes_base(name))
}
