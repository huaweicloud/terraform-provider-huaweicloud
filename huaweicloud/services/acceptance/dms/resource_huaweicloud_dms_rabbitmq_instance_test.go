package dms

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
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

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
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
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "product_id", "${data.huaweicloud_dms_product.test2.id}"),
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
					"password", "used_storage_space",
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
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_compatible(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccDmsRabbitmqInstance_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "used_storage_space",
				},
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_single(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccDmsRabbitmqInstance_single(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDmsRabbitmqInstance_Base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  description = "Test for DMS RabbitMQ"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[1]s"
  description = "secgroup for rabbitmq"
}

data "huaweicloud_dms_product" "test1" {
  engine        = "rabbitmq"
  instance_type = "cluster"
  version       = "3.7.17"
  node_num      = 3
}

data "huaweicloud_dms_product" "test2" {
  engine        = "rabbitmq"
  instance_type = "cluster"
  version       = "3.7.17"
  node_num      = 5
}

data "huaweicloud_availability_zones" "test" {}
`, rName)
}

func testAccDmsRabbitmqInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  product_id        = data.huaweicloud_dms_product.test1.id
  engine_version    = data.huaweicloud_dms_product.test1.version
  storage_spec_code = data.huaweicloud_dms_product.test1.storage_spec_code

  access_user = "user"
  password    = "Rabbitmqtest@123"

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
  name        = "%s"
  description = "rabbitmq test update"

  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  product_id        = data.huaweicloud_dms_product.test2.id
  engine_version    = data.huaweicloud_dms_product.test2.version
  storage_spec_code = data.huaweicloud_dms_product.test2.storage_spec_code

  access_user = "user"
  password    = "Rabbitmqtest@123"

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
  enterprise_project_id = "%s"

  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1]
  ]
  
  product_id        = data.huaweicloud_dms_product.test1.id
  engine_version    = data.huaweicloud_dms_product.test1.version
  storage_spec_code = data.huaweicloud_dms_product.test1.storage_spec_code

  access_user = "user"
  password    = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// After the 1.31.1 version, arguments storage_space and available_zones are deprecated.
func testAccDmsRabbitmqInstance_compatible(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_az" "test" {}
resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name        = "%s"
  description = "rabbitmq test"
  
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]

  product_id        = data.huaweicloud_dms_product.test1.id
  engine_version    = data.huaweicloud_dms_product.test1.version
  storage_space     = data.huaweicloud_dms_product.test1.storage
  storage_spec_code = data.huaweicloud_dms_product.test1.storage_spec_code


  access_user = "user"
  password    = "Rabbitmqtest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccDmsRabbitmqInstance_Base(rName), rName)
}

func testAccDmsRabbitmqInstance_single(rName string) string {
	randPwd := fmt.Sprintf("%s!#%d", acctest.RandString(5), acctest.RandIntRange(0, 999))
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_product" "single" {
  engine           = "rabbitmq"
  instance_type    = "single"
  version          = "3.7.17"
  node_num         = 1
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  product_id        = data.huaweicloud_dms_product.single.id
  engine_version    = data.huaweicloud_dms_product.single.version
  storage_spec_code = data.huaweicloud_dms_product.single.storage_spec_code
  storage_space     = data.huaweicloud_dms_product.single.storage

  access_user = "root"
  password    = "%[3]s"
}
`, testAccDmsRabbitmqInstance_Base(rName), rName, randPwd)
}
