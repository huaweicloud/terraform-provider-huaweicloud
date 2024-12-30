package vpc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getNetworkAclResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	getNetworkAclHttpUrl := "v3/{project_id}/vpc/firewalls/" + state.Primary.ID
	getNetworkAclPath := client.Endpoint + getNetworkAclHttpUrl
	getNetworkAclPath = strings.ReplaceAll(getNetworkAclPath, "{project_id}", client.ProjectID)

	getNetworkAclOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getNetworkAclResp, err := client.Request("GET", getNetworkAclPath, &getNetworkAclOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving VPC network ACL: %s", err)
	}

	return utils.FlattenResponse(getNetworkAclResp)
}

func TestAccNetworkAcl_basic(t *testing.T) {
	var (
		networkAcl   interface{}
		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_vpc_network_acl.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&networkAcl,
			getNetworkAclResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAcl_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),

					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.source_ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.source_port", "22-30,33"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.destination_ip_address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.destination_port", "8001-8010"),

					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.source_ip_address", "172.16.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.source_port", "22-30,33"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.destination_ip_address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.destination_port", "8001-8010"),

					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),

					resource.TestCheckOutput("is_associated_subnets_different", "false"),

					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccNetworkAcl_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.source_ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.source_port", "22-30,33"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.destination_ip_address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.0.destination_port", "8001-8010"),

					resource.TestCheckResourceAttr(resourceName, "ingress_rules.1.action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.1.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.1.protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.1.source_ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "ingress_rules.1.destination_ip_address", "0.0.0.0/0"),

					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.source_ip_address", "172.16.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.source_port", "22-30,33"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.destination_ip_address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.0.destination_port", "8001-8010"),

					resource.TestCheckResourceAttr(resourceName, "egress_rules.1.action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.1.ip_version", "4"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.1.protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.1.source_ip_address", "172.16.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "egress_rules.1.destination_ip_address", "0.0.0.0/0"),

					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),

					resource.TestCheckOutput("is_associated_subnets_different", "false"),

					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config:            testAccNetworkAcl_import(name),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkAcl_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}
`, name)
}

func testAccNetworkAcl_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test" {
  count      = 1
  name       = "%s-${count.index}"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index), 1)
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_network_acl" "test" {
  name        = "%s"
  description = "created by terraform"

  ingress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "192.168.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  egress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "172.16.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }
  
  associated_subnets {
    subnet_id = huaweicloud_vpc_subnet.test[0].id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

output "is_associated_subnets_different" {
  value = length(setsubtract(huaweicloud_vpc_network_acl.test.associated_subnets[*].subnet_id,
  huaweicloud_vpc_subnet.test[*].id)) != 0
}
`, testAccNetworkAcl_base(name), name, name)
}

func testAccNetworkAcl_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test" {
  count      = 2
  name       = "%s-${count.index}"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index), 1)
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_network_acl" "test" {
  name        = "%s-update"
  description = "created by terraform update"
  enabled     = false

  ingress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "192.168.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  ingress_rules {
    action                 = "deny"
    ip_version             = 4
    protocol               = "icmp"
    source_ip_address      = "192.168.0.0/24"
    destination_ip_address = "0.0.0.0/0"
  }

  egress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "172.16.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  egress_rules {
    action                 = "deny"
    ip_version             = 4
    protocol               = "icmp"
    source_ip_address      = "172.16.0.0/24"
    destination_ip_address = "0.0.0.0/0"
  }
  
  associated_subnets {
    subnet_id = huaweicloud_vpc_subnet.test[0].id
  }

  associated_subnets {
    subnet_id = huaweicloud_vpc_subnet.test[1].id
  }

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}

output "is_associated_subnets_different" {
  value = length(setsubtract(huaweicloud_vpc_network_acl.test.associated_subnets[*].subnet_id,
  huaweicloud_vpc_subnet.test[*].id)) != 0
}
`, testAccNetworkAcl_base(name), name, name)
}

func testAccNetworkAcl_import(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test" {
  count      = 2
  name       = "%s-${count.index}"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, count.index), 1)
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_network_acl" "test" {
  name        = "%s-update"
  description = "created by terraform update"
  enabled     = false

  ingress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "192.168.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  ingress_rules {
    action                 = "deny"
    ip_version             = 4
    protocol               = "icmp"
    source_ip_address      = "192.168.0.0/24"
    destination_ip_address = "0.0.0.0/0"
  }

  egress_rules {
    action                 = "allow"
    ip_version             = 4
    protocol               = "tcp"
    source_ip_address      = "172.16.0.0/24"
    source_port            = "22-30,33"
    destination_ip_address = "0.0.0.0/0"
    destination_port       = "8001-8010"
  }

  egress_rules {
    action                 = "deny"
    ip_version             = 4
    protocol               = "icmp"
    source_ip_address      = "172.16.0.0/24"
    destination_ip_address = "0.0.0.0/0"
  }
  
  associated_subnets {
    subnet_id = huaweicloud_vpc_subnet.test[0].id
  }

  associated_subnets {
    subnet_id = huaweicloud_vpc_subnet.test[1].id
  }
}
`, testAccNetworkAcl_base(name), name, name)
}
