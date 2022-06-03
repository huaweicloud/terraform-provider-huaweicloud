package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sms/v3/templates"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getServerTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SmsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return templates.Get(client, state.Primary.ID)
}

func TestAccServerTemplate_basic(t *testing.T) {
	var temp templates.TemplateResponse
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sms_server_template.test"
	azDataName := "data.huaweicloud_availability_zones.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&temp,
		getServerTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServerTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "target_server_name", name),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_name", "autoCreate"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.0", "autoCreate"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "autoCreate"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", azDataName, "names.0"),
				),
			},
			{
				Config: testAccServerTemplate_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "target_server_name", name),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", azDataName, "names.1"),
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

func TestAccServerTemplate_existing(t *testing.T) {
	var temp templates.TemplateResponse
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sms_server_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&temp,
		getServerTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServerTemplate_existing(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "target_server_name", name),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_name", name),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "autoCreate"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.0", "huaweicloud_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testAccServerTemplate_existing_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.0", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0", "huaweicloud_networking_secgroup.test", "id"),
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

func testAccServerTemplate_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sms_server_template" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
`, name)
}

func testAccServerTemplate_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sms_server_template" "test" {
  name               = "%s-update"
  target_server_name = "%s"
  availability_zone  = data.huaweicloud_availability_zones.test.names[1]
  volume_type        = "GPSSD"
}
`, name, name)
}

func testAccServerTemplate_network(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  description = "created by acc test"
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[1]s"
  description = "created by acc test"
}
`, name)
}

func testAccServerTemplate_existing(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sms_server_template" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_ids        = [ huaweicloud_vpc_subnet.test.id ]
}
`, testAccServerTemplate_network(name), name)
}

func testAccServerTemplate_existing_update(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sms_server_template" "test" {
  name               = "%s"
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  vpc_id             = huaweicloud_vpc.test.id
  subnet_ids         = [ huaweicloud_vpc_subnet.test.id ]
  security_group_ids = [ huaweicloud_networking_secgroup.test.id ]
}
`, testAccServerTemplate_network(name), name)
}
