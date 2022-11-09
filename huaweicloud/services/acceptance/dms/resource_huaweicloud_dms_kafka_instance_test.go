package dms

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getKafkaInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccKafkaInstance_basic(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_kafka_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestMatchResourceAttr(resourceName, "cross_vpc_accesses.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
			{
				Config: testAccKafkaInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
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
					"manager_password",
					"used_storage_space",
					"cross_vpc_accesses",
				},
			},
		},
	})
}

func TestAccKafkaInstance_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccKafkaInstance_compatible(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_compatible(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_space", "data.huaweicloud_dms_product.test", "storage"),
				),
			},
		},
	})
}

func TestAccKafkaInstance_newFormat(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_newFormat(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttrPair(resourceName, "broker_num",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.properties.0.min_broker"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_spec_code",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.ios.0.storage_spec_code"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.1.advertised_ip", "www.terraform-test.com"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.2.advertised_ip", "192.168.0.53"),
				),
			},
			{
				Config: testAccKafkaInstance_newFormatUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttrPair(resourceName, "broker_num",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.properties.0.min_broker"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_spec_code",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.ios.1.storage_spec_code"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.0.advertised_ip", "172.16.35.62"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.1.advertised_ip", "www.terraform-test-1.com"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.2.advertised_ip", "192.168.0.53"),
				),
			},
		},
	})
}

func testAccKafkaInstance_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name        = "%s"
  cidr        = "192.168.11.0/24"
  description = "test for kafka"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.11.0/24"
  gateway_ip = "192.168.11.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for kafka"
}
`, rName, rName, rName)
}

func testAccKafkaInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  bandwidth         = "100MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test"
  access_user        = "user"
  password           = "Kafkatest@123"
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  manager_user      = "kafka-user"
  manager_password  = "Kafkatest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccKafkaInstance_base(rName), rName)
}

func testAccKafkaInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test update"
  access_user        = "user"
  password           = "Kafkatest@123"
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code

  manager_user      = "kafka-user"
  manager_password  = "Kafkatest@123"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
`, testAccKafkaInstance_base(rName), updateName)
}

func testAccKafkaInstance_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name                  = "%s"
  description           = "kafka test"
  access_user           = "user"
  password              = "Kafkatest@123"
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
  product_id            = data.huaweicloud_dms_product.test.id
  engine_version        = data.huaweicloud_dms_product.test.version
  storage_spec_code     = data.huaweicloud_dms_product.test.storage_spec_code

  manager_user          = "kafka-user"
  manager_password      = "Kafkatest@123"
  enterprise_project_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccKafkaInstance_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKafkaInstance_compatible(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_az" "test" {}

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  instance_type     = "cluster"
  version           = "2.3.0"
  bandwidth         = "300MB"
  storage_spec_code = "dms.physical.storage.ultra"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name        = "%s"
  description = "kafka test"

  # use deprecated argument "available_zones"
  available_zones   = [data.huaweicloud_dms_az.test.id]
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code
  storage_space     = data.huaweicloud_dms_product.test.storage
  # use deprecated argument "bandwidth"
  bandwidth         = data.huaweicloud_dms_product.test.bandwidth

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, testAccKafkaInstance_base(rName), rName)
}

func testAccKafkaInstance_newFormat(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_kafka_flavors.test

  advertised_ips = ["", "www.terraform-test.com", "192.168.0.53"]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.query_results.flavors[0].id
  storage_spec_code  = local.query_results.flavors[0].ios[0].storage_spec_code
  availability_zones = local.query_results.flavors[0].ios[0].availability_zones
  engine_version     = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space      = local.query_results.flavors[0].properties[0].min_broker * local.query_results.flavors[0].properties[0].min_storage_per_node
  broker_num         = local.query_results.flavors[0].properties[0].min_broker

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  dynamic "cross_vpc_accesses" {
    for_each = local.advertised_ips
    content {
      advertised_ip = cross_vpc_accesses.value
    }
  }
}`, testAccKafkaInstance_base(rName), rName)
}

func testAccKafkaInstance_newFormatUpdate(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_kafka_flavors.test

  advertised_ips = ["172.16.35.62", "www.terraform-test-1.com", "192.168.0.53"]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.query_results.flavors[0].id
  storage_spec_code  = local.query_results.flavors[0].ios[1].storage_spec_code
  availability_zones = local.query_results.flavors[0].ios[0].availability_zones
  engine_version     = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space      = local.query_results.flavors[0].properties[0].min_broker * local.query_results.flavors[0].properties[0].max_storage_per_node
  broker_num         = local.query_results.flavors[0].properties[0].min_broker

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  dynamic "cross_vpc_accesses" {
    for_each = local.advertised_ips
    content {
      advertised_ip = cross_vpc_accesses.value
    }
  }
}`, testAccKafkaInstance_base(rName), rName)
}
