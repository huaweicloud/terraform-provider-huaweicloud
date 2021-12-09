package dms

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDmsRabitMqInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccDmsRabbitmqInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccDmsRabbitmqInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDmsRabitMqInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID),
				),
			},
		},
	})
}

func testAccDmsRabbitmqInstance_Base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_az" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dms_product" "test" {
  engine        = "rabbitmq"
  instance_type = "cluster"
  version       = "3.7.17"
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for rabbitmq"
}
`, rName)
}

func testAccDmsRabbitmqInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name              = "%s"
  description       = "rabbitmq test"
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  vpc_id            = data.huaweicloud_vpc.test.id
  network_id        = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_space     = data.huaweicloud_dms_product.test.storage
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName)
}

func testAccDmsRabbitmqInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name              = "%s"
  description       = "rabbitmq test update"
  access_user       = "user"
  password          = "Rabbitmqtest@123"
  vpc_id            = data.huaweicloud_vpc.test.id
  network_id        = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_space     = data.huaweicloud_dms_product.test.storage
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), updateName)
}

func testAccDmsRabbitmqInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name                  = "%s"
  description           = "rabbitmq test"
  access_user           = "user"
  password              = "Rabbitmqtest@123"
  vpc_id                = data.huaweicloud_vpc.test.id
  network_id            = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  available_zones       = [data.huaweicloud_dms_az.test.id]
  product_id            = data.huaweicloud_dms_product.test.id
  engine_version        = data.huaweicloud_dms_product.test.version
  storage_space         = data.huaweicloud_dms_product.test.storage
  storage_spec_code     = data.huaweicloud_dms_product.test.storage_spec_code
  enterprise_project_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
