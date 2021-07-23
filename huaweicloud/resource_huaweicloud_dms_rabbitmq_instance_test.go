package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v2/rabbitmq/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsRabbitmqInstances_basic(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsRabbitmqInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsRabbitmqInstanceExists(resourceName, instance),
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
					"password",
					"used_storage_space",
				},
			},
			{
				Config: testAccDmsRabbitmqInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsRabbitmqInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "rabbitmq test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
		},
	})
}

func TestAccDmsRabbitmqInstances_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_dms_rabbitmq_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsRabbitmqInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitmqInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsRabbitmqInstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDmsRabbitmqInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_rabbitmq_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("The Dms rabbitmq instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsRabbitmqInstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud dms instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmtp.Errorf("Error getting HuaweiCloud dms rabbitmq instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmtp.Errorf("The Dms rabbitmq instance not found.")
		}
		instance = *v
		return nil
	}
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
`, testAccDmsRabbitmqInstance_Base(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
