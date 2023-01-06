package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/applications"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return applications.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccApplication_basic(t *testing.T) {
	var (
		app applications.Application

		rName = "huaweicloud_apig_application.test"
		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		description       = "Created by script"
		updateDescription = "Updated by script"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&app,
		getApplicationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(name, description),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", description),
					resource.TestCheckResourceAttrSet(rName, "app_key"),
					resource.TestCheckResourceAttrSet(rName, "app_secret"),
				),
			},
			{
				// update name, description and app_code.
				Config: testAccApplication_basic(updateName, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", updateDescription),
					resource.TestCheckResourceAttrSet(rName, "app_key"),
					resource.TestCheckResourceAttrSet(rName, "app_secret"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationImportIdFunc(),
			},
		},
	})
}

func testAccApplicationImportIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_application.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccApigApplication_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
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

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}
`, name)
}

func testAccApplication_basic(name, description string) string {
	code := utils.EncodeBase64String(acctest.RandString(64))
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "%[3]s"

  app_codes = ["%[4]s"]
}
`, testAccApigApplication_base(name), name, description, code)
}
