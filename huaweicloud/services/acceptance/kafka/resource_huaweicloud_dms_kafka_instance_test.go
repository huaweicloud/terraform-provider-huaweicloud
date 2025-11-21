package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getKafkaInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccKafkaInstance_prePaid(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "huaweicloud_dms_kafka_instance.test"
	baseNetwork := common.TestBaseNetwork(rName)

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
				Config: testAccKafkaInstance_newFormat_prePaid(baseNetwork, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "10:00:00"),
				),
			},
			{
				Config: testAccKafkaInstance_newFormat_prePaid_update(baseNetwork, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "06:00:00"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_user",
					"manager_password",
					"used_storage_space",
					"cross_vpc_accesses",
					"auto_renew",
					"period",
					"period_unit",
				},
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_newFormat(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "2.7"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_spec_code",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.ios.0.storage_spec_code"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestMatchResourceAttr(resourceName, "availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceName, "arch_type", "X86"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "access_user", "user"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "log.retention.hours"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "48"),
					resource.TestCheckResourceAttr(resourceName, "security_protocol", "SASL_PLAINTEXT"),
					resource.TestCheckResourceAttr(resourceName, "enabled_mechanisms.0", "SCRAM-SHA-512"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.1.advertised_ip", "www.terraform-test.com"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.2.advertised_ip", "192.168.0.53"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_space"),
					resource.TestCheckResourceAttrSet(resourceName, "maintain_begin"),
					resource.TestCheckResourceAttrSet(resourceName, "maintain_end"),
					// Check attributes.
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestMatchResourceAttr(resourceName, "partition_num", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(resourceName, "connect_address"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "is_logical_volume", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "node_num"),
					resource.TestCheckResourceAttrSet(resourceName, "pod_connect_address"),
					resource.TestCheckResourceAttr(resourceName, "type", "cluster"),
				),
			},
			{
				Config: testAccKafkaInstance_newFormatUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "4"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor_id",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_spec_code",
						"data.huaweicloud_dms_kafka_flavors.test", "flavors.0.ios.0.storage_spec_code"),
					resource.TestCheckResourceAttr(resourceName, "storage_space", "600"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "auto.create.groups.enable"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "false"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.0.advertised_ip", "192.168.0.61"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.1.advertised_ip", "test.terraform.com"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.2.advertised_ip", "192.168.0.62"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_accesses.3.advertised_ip", "192.168.0.63"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "10:00:00"),
				),
			},
		},
	})
}

func TestAccKafkaInstance_publicIp(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_publicIp(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_ids.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "dumping", "true"),
					resource.TestCheckResourceAttr(resourceName, "new_tenant_ips.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "new_tenant_ips.0", "192.168.0.20"),
					resource.TestCheckResourceAttr(resourceName, "new_tenant_ips.1", "192.168.0.18"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.private_plain_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.public_plain_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.private_sasl_ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.public_sasl_ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.private_sasl_plaintext_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.public_sasl_plaintext_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_auto_topic", "true"),
					resource.TestCheckResourceAttr(resourceName, "vpc_client_plain", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by Terraform script"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "retention_policy", "time_base"),
					// Check attributes.
					resource.TestCheckResourceAttr(resourceName, "enable_public_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_address.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "connector_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connector_node_num"),
				),
			},
			{
				Config: testAccKafkaInstance_publicIp_update(rName, password, 5),
				ExpectError: regexp.MustCompile("error resizing instance: the old EIP ID should not be changed, and the adding nums of " +
					"EIP ID should be same as the adding broker nums"),
			},
			{
				Config: testAccKafkaInstance_publicIp_update(rName, password, 4),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broker_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_ids.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "new_tenant_ips.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "new_tenant_ips.0", "192.168.0.79"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.private_plain_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.private_sasl_ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.public_plain_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "port_protocol.0.public_sasl_ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_auto_topic", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "retention_policy", "produce_reject"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"enabled_mechanisms",
					"arch_type",
					"new_tenant_ips",
				},
			},
		},
	})
}

func testAccKafkaInstance_newFormat(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2]
  ]
  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3
  arch_type      = "X86"

  ssl_enable         = true
  access_user        = "user"
  password           = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = ""
  }
  cross_vpc_accesses {
    advertised_ip = "www.terraform-test.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.53"
  }

  parameters {
    name  = "log.retention.hours"
    value = "48"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccKafkaInstance_newFormatUpdate(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.4u8g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2]
  ]
  engine_version = "2.7"
  storage_space  = 600
  broker_num     = 4
  arch_type      = "X86"

  ssl_enable         = true
  access_user        = "user"
  password           = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]
  maintain_begin     = "06:00:00"
  maintain_end       = "10:00:00"

  cross_vpc_accesses {
    advertised_ip = "192.168.0.61"
  }
  cross_vpc_accesses {
    advertised_ip = "test.terraform.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.62"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.63"
  }

  parameters {
    name  = "auto.create.groups.enable"
    value = "false"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccKafkaInstance_newFormat_prePaid(baseNetwork, rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3

  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"
  maintain_begin   = "06:00"
  maintain_end     = "10:00"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, baseNetwork, rName)
}

func testAccKafkaInstance_newFormat_prePaid_update(baseNetwork, updateName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  description       = "kafka test update"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  flavor_id         = local.flavor.id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3

  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"
  maintain_begin   = "02:00"
  maintain_end     = "06:00"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}`, baseNetwork, updateName)
}

func testAccKafkaInstance_publicIpBase(count int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  count = %d

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "test_eip_${count.index}"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, count)
}

func testAccKafkaInstance_publicIp(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

%s

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 3)

  engine_version        = "2.7"
  storage_space         = 300
  broker_num            = 3
  arch_type             = "X86"
  public_ip_ids         = huaweicloud_vpc_eip.test[*].id
  new_tenant_ips        = ["192.168.0.20", "192.168.0.18"]
  dumping               = true
  description           = "Created by Terraform script"
  enable_auto_topic     = true
  vpc_client_plain      = true
  enterprise_project_id = "%[4]s"
  retention_policy      = "time_base"

  port_protocol {
    private_plain_enable = true
    public_plain_enable  = true
  }
}`, common.TestBaseNetwork(rName), testAccKafkaInstance_publicIpBase(3), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccKafkaInstance_publicIp_update(rName, password string, brokerNum int) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

%s

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 3)

  engine_version        = "2.7"
  storage_space         = 600
  broker_num            = %d
  arch_type             = "X86"
  new_tenant_ips        = ["192.168.0.79"]
  dumping               = true
  public_ip_ids         = huaweicloud_vpc_eip.test[*].id
  access_user           = "test"
  password              = "%[5]s"
  enabled_mechanisms    = ["SCRAM-SHA-512"]
  enable_auto_topic     = false
  enterprise_project_id = "%[6]s"
  retention_policy      = "produce_reject"

  port_protocol {
    private_plain_enable    = false
    private_sasl_ssl_enable = true
    public_plain_enable     = false
    public_sasl_ssl_enable  = true
  }
}`, common.TestBaseNetwork(rName), testAccKafkaInstance_publicIpBase(4), rName, brokerNum, password, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}
