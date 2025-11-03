package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUniAgentBatchInstall_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAomInstallerAgentId(t)
			acceptance.TestAccPreCheckAomUniAgentVersion(t)
			acceptance.TestAccPreCheckAomTargetMachineAccount(t)
			acceptance.TestAccPreCheckAomTargetMachinePassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testUniAgentBatchInstall_basic(rName),
			},
		},
	})
}

func testUniAgentBatchInstall_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_vpc" "test" {
  id = "%[2]s"
}

data "huaweicloud_vpc_subnet" "test" {
  vpc_id = data.huaweicloud_vpc.test.id
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%[1]s"
  description = "Security group for Terraform resources"
}

variable "security_group_rules_configuration" {
  type = list(object({
    direction = string
    protocol  = string
  }))

  default = [
    {direction = "ingress", protocol = "tcp"},
    {direction = "ingress", protocol = "udp"},
    {direction = "ingress", protocol = "icmp"},
    {direction = "egress", protocol = "tcp"},
    {direction = "egress", protocol = "udp"},
    {direction = "egress", protocol = "icmp"},
  ]
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = length(var.security_group_rules_configuration)

  direction         = var.security_group_rules_configuration[count.index].direction
  ethertype         = "IPv4"
  protocol          = var.security_group_rules_configuration[count.index].protocol
  ports             = var.security_group_rules_configuration[count.index].protocol == "icmp" ? "0" : "1-65535"
  remote_ip_prefix  = data.huaweicloud_vpc.test.cidr
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name                = "%[1]s"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = try(data.huaweicloud_compute_flavors.test.ids[0], "NOT_FOUND")
  system_disk_type    = "SAS"
  system_disk_size    = 40
  security_group_ids  = [huaweicloud_networking_secgroup.test.id] 
  admin_pass          = "%[3]s"
  stop_before_destroy = true

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    test = "uniagent_install_target"
  }
}
`, rName, acceptance.HW_VPC_ID, acceptance.HW_AOM_TARGET_MACHINE_PASSWORD)
}

func testUniAgentBatchInstall_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_uniagent_batch_install" "test" {
  region             = "%[2]s"
  proxy_region_id    = 0
  installer_agent_id = "%[3]s"
  version            = "%[4]s"
  public_net_flag    = false

  dynamic "agent_import_param_list" {
    for_each = huaweicloud_compute_instance.test

    content {
      account  = "%[5]s"
      password = "%[6]s"
      inner_ip = agent_import_param_list.value.access_ip_v4
      port     = 22
      os_type  = "LINUX"
    }
  }

  depends_on = [huaweicloud_compute_instance.test]
}
`, testUniAgentBatchInstall_base(rName), acceptance.HW_REGION_NAME, acceptance.HW_AOM_INSTALLER_AGENT_ID,
		acceptance.HW_AOM_UNIAGENT_VERSION, acceptance.HW_AOM_TARGET_MACHINE_ACCOUNT, acceptance.HW_AOM_TARGET_MACHINE_PASSWORD)
}
