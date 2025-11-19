package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_kafka_instances.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byInstanceId   = "data.huaweicloud_dms_kafka_instances.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byName   = "data.huaweicloud_dms_kafka_instances.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byFuzzyName   = "data.huaweicloud_dms_kafka_instances.filter_by_name_fuzzy"
		dcByFuzzyName = acceptance.InitDataSourceCheck(byFuzzyName)

		byStatus   = "data.huaweicloud_dms_kafka_instances.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEpsId   = "data.huaweicloud_dms_kafka_instances.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_id_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_name_useful", "true"),
					dcByFuzzyName.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_name_fuzzy_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.type"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.description"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.product_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_spec_code"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_space"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_type"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.storage_resource_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.vpc_name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.vpc_client_plain"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.network_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.security_group_name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.maintain_begin"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.maintain_end"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.enable_public_ip", "true"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.public_ip_ids.#", "3"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.public_conn_addresses"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.retention_policy"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.dumping", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.connector_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.connector_node_num"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.enable_auto_topic", "true"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.partition_num", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.ssl_enable", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.security_protocol"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.enabled_mechanisms.#"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.used_storage_space"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.connect_address"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.port"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.status"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.specification"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.user_id"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.user_name"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.tags.%"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.broker_num"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.ces_version"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.extend_times"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.is_logical_volume"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.message_query_inst_enable", "true"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.new_spec_billing_enable", "true"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.node_num", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.order_id", ""),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.pod_connect_address"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.port_protocol.#", "1"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.port_protocol.0.private_plain_enable", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.port_protocol.0.private_plain_address"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.port_protocol.0.private_sasl_ssl_enable", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.port_protocol.0.private_sasl_ssl_address"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.port_protocol.0.public_plain_enable", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.port_protocol.0.public_plain_address"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.port_protocol.0.public_sasl_plaintext_enable", "true"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.port_protocol.0.public_sasl_plaintext_address"),
					resource.TestCheckResourceAttr(byInstanceId, "instances.0.support_features.#", "0"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.cross_vpc_accesses.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.cross_vpc_accesses.0.listener_ip"),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.cross_vpc_accesses.0.advertised_ip"),
					resource.TestMatchResourceAttr(byInstanceId, "instances.0.cross_vpc_accesses.0.port", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.cross_vpc_accesses.0.port_id"),
					// deprecated attributes.
					resource.TestCheckResourceAttrSet(byInstanceId, "instances.0.cross_vpc_accesses.0.lisenter_ip"),
					// management_connect_address and resource_spec_code are not returned.
					// Only some region supports IPV6, so we don't check the ipv6_enable and ipv6_connect_addresses.
				),
			},
		},
	})
}

func testAccDataInstances_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = try(data.huaweicloud_dms_kafka_flavors.test.flavors[0], {})
}

# The number of EIPs to be created is equal to the number of brokers for the instance.
resource "huaweicloud_vpc_eip" "test" {
  count = 3

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s${count.index}"
    size        = 1
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%[2]s"
  flavor_id         = local.flavor.id
  engine_version    = "3.x"
  storage_space     = try(local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node, null)
  storage_spec_code = try(local.flavor.ios[0].storage_spec_code, null)
  broker_num        = 3

  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  access_user        = "%[2]s"
  password           = "%[3]s"
  enabled_mechanisms = ["SCRAM-SHA-512"]
  public_ip_ids      = huaweicloud_vpc_eip.test[*].id
  enable_auto_topic  = true
  description        = "Create by Terraform script"
  dumping            = true

  port_protocol {
    private_plain_enable         = true
    private_sasl_ssl_enable      = true
    public_plain_enable          = true
    public_sasl_plaintext_enable = true
  }

  tags = {
    foo = "bar"
  }
}
`, common.TestBaseNetwork(name), name, acceptance.RandomPassword())
}

func testAccDataInstances_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_instances" "all" {
  depends_on = [huaweicloud_dms_kafka_instance.test]
}

# Filter by instance ID.
locals {
  instance_id = huaweicloud_dms_kafka_instance.test.id
}

data "huaweicloud_dms_kafka_instances" "filter_by_instance_id" {
  instance_id = local.instance_id

  depends_on = [huaweicloud_dms_kafka_instance.test]
}

locals {
  instance_id_filter_result = [for v in data.huaweicloud_dms_kafka_instances.filter_by_instance_id.instances[*].id :
  v == local.instance_id]
}

output "is_instance_id_useful" {
  value = length(local.instance_id_filter_result) == 1 && alltrue(local.instance_id_filter_result)
}

# Filter by instance name with exact match.
locals {
  instance_name       = try(data.huaweicloud_dms_kafka_instances.all.instances[0].name, "")
  instance_name_fuzzy = try(substr(local.instance_name, 0, length(local.instance_name) - 1), local.instance_name)
}

# fuzzy_match default value is false (exact match).
data "huaweicloud_dms_kafka_instances" "filter_by_name" {
  name = local.instance_name

  depends_on = [huaweicloud_dms_kafka_instance.test]
}

locals {
  instance_name_filter_result = [for v in data.huaweicloud_dms_kafka_instances.filter_by_name.instances[*].name :
  v == local.instance_name]
}

output "is_instance_name_useful" {
  value = length(local.instance_name_filter_result) == 1 && alltrue(local.instance_name_filter_result)
}

# Filter by instance name with fuzzy match.
data "huaweicloud_dms_kafka_instances" "filter_by_name_fuzzy" {
  fuzzy_match = true
  name        = local.instance_name_fuzzy

  depends_on = [huaweicloud_dms_kafka_instance.test]
}

locals {
  instance_name_fuzzy_filter_result = [for v in data.huaweicloud_dms_kafka_instances.filter_by_name_fuzzy.instances[*].name :
  strcontains(v, local.instance_name_fuzzy)]
}

output "is_instance_name_fuzzy_useful" {
  value = length(local.instance_name_fuzzy_filter_result) > 0 && alltrue(local.instance_name_fuzzy_filter_result)
}

# Filter by instance status.
locals {
  status = try(data.huaweicloud_dms_kafka_instances.all.instances[0].status, "")
}

data "huaweicloud_dms_kafka_instances" "filter_by_status" {
  status = local.status

  depends_on = [huaweicloud_dms_kafka_instance.test]
}

locals {
  status_filter_result = [for v in data.huaweicloud_dms_kafka_instances.filter_by_status.instances[*].status :
  v == local.status]
}

output "is_status_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by enterprise project ID to which the instance belongs.
locals {
  eps_id = try(data.huaweicloud_dms_kafka_instances.all.instances[0].enterprise_project_id, "")
}

data "huaweicloud_dms_kafka_instances" "filter_by_eps_id" {
  enterprise_project_id = local.eps_id

  depends_on = [huaweicloud_dms_kafka_instance.test]
}

locals {
  eps_id_filter_result = [for v in data.huaweicloud_dms_kafka_instances.filter_by_eps_id.instances[*].enterprise_project_id :
  v == local.eps_id]
}

output "is_enterprise_project_id_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}
`, testAccDataInstances_base())
}
