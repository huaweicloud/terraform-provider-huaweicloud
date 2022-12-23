package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getEnvironmentFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return apig.GetEnvironmentFormServer(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccEnvironment_basic(t *testing.T) {
	var (
		env environments.Environment

		rName = "huaweicloud_apig_environment.test"
		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&env,
		getEnvironmentFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccEnvironment_update(updateName),
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
				ImportStateIdFunc: testAccEnvironmentImportStateFunc(),
			},
		},
	})
}

func testAccEnvironmentImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_environment.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccEnvironment_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[1]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, name)
}

func testAccEnvironment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_environment" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, testAccEnvironment_base(name), name)
}

func testAccEnvironment_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_environment" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
}
`, testAccEnvironment_base(name), name)
}
