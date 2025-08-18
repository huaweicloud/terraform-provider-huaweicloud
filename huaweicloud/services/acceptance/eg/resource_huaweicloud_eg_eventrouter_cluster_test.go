package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eg"
)

func getEventRouterClusterFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("eg", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG client: %s", err)
	}

	return eg.GetEventRouterClusterById(client, state.Primary.ID)
}

func TestAccEventRouterCluster_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_eg_eventrouter_cluster.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getEventRouterClusterFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEventRouterCluster_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "source_type", "KAFKA"),
					resource.TestCheckResourceAttr(rName, "sink_type", "KAFKA"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccEventRouterCluster_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "source_type", "KAFKA"),
					resource.TestCheckResourceAttr(rName, "sink_type", "KAFKA"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"flavor",
				},
			},
		},
	})
}

func testAccEventRouterCluster_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 1), 1)
  vpc_id     = huaweicloud_vpc.test.id
}
`, name)
}

func testAccEventRouterCluster_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_eventrouter_cluster" "test" {
  name               = "%[2]s"
  source_type        = "KAFKA"
  sink_type          = "KAFKA"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  description        = "Created by terraform script"
  availability_zones = join(",", slice(data.huaweicloud_availability_zones.test.names, 0, 2))
  flavor             = "small"
}
`, testAccEventRouterCluster_base(name), name)
}

func testAccEventRouterCluster_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_eventrouter_cluster" "test" {
  name               = "%[2]s"
  source_type        = "KAFKA"
  sink_type          = "KAFKA"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  availability_zones = join(",", slice(data.huaweicloud_availability_zones.test.names, 0, 2))
  flavor             = "small"
}
`, testAccEventRouterCluster_base(name), name)
}
