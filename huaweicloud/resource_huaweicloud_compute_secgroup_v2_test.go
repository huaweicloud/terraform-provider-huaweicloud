package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccComputeV2SecGroup_basic(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroup_basic_orig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
			{
				ResourceName:      "huaweicloud_compute_secgroup_v2.sg_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComputeV2SecGroup_update(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroup_basic_orig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
			{
				Config: testAccComputeV2SecGroup_basic_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup),
					testAccCheckComputeV2SecGroupRuleCount(&secgroup, 2),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_groupID(t *testing.T) {
	var secgroup1, secgroup2, secgroup3 secgroups.SecurityGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroup_groupID_orig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_2", &secgroup2),
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_3", &secgroup3),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup1, &secgroup3),
				),
			},
			{
				Config: testAccComputeV2SecGroup_groupID_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_2", &secgroup2),
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_3", &secgroup3),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup2, &secgroup3),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_self(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroup_self,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup, &secgroup),
					resource.TestCheckResourceAttr(
						"huaweicloud_compute_secgroup_v2.sg_1", "rule.3170486100.self", "true"),
					resource.TestCheckResourceAttr(
						"huaweicloud_compute_secgroup_v2.sg_1", "rule.3170486100.from_group_id", ""),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_icmpZero(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroup_icmpZero,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_timeout(t *testing.T) {
	var secgroup secgroups.SecurityGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroup_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("huaweicloud_compute_secgroup_v2.sg_1", &secgroup),
				),
			},
		},
	})
}

func testAccCheckComputeV2SecGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_secgroup_v2" {
			continue
		}

		_, err := secgroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2SecGroupExists(n string, secgroup *secgroups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := secgroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Security group not found")
		}

		*secgroup = *found

		return nil
	}
}

func testAccCheckComputeV2SecGroupRuleCount(secgroup *secgroups.SecurityGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(secgroup.Rules) != count {
			return fmtp.Errorf("Security group rule count does not match. Expected %d, got %d", count, len(secgroup.Rules))
		}

		return nil
	}
}

func testAccCheckComputeV2SecGroupGroupIDMatch(sg1, sg2 *secgroups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg2.Rules) == 1 {
			if sg1.Name != sg2.Rules[0].Group.Name || sg1.TenantID != sg2.Rules[0].Group.TenantID {
				return fmtp.Errorf("%s was not correctly applied to %s", sg1.Name, sg2.Name)
			}
		} else {
			return fmtp.Errorf("%s rule count is incorrect", sg2.Name)
		}

		return nil
	}
}

const testAccComputeV2SecGroup_basic_orig = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = 1
    to_port = 65535
    ip_protocol = "udp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}
`

const testAccComputeV2SecGroup_basic_update = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 2200
    to_port = 2200
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}
`

const testAccComputeV2SecGroup_groupID_orig = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "huaweicloud_compute_secgroup_v2" "sg_2" {
  name = "sg_2"
  description = "second test security group"
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}

resource "huaweicloud_compute_secgroup_v2" "sg_3" {
  name = "sg_3"
  description = "third test security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    from_group_id = "${huaweicloud_compute_secgroup_v2.sg_1.id}"
  }
}
`

const testAccComputeV2SecGroup_groupID_update = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "huaweicloud_compute_secgroup_v2" "sg_2" {
  name = "sg_2"
  description = "second test security group"
  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}

resource "huaweicloud_compute_secgroup_v2" "sg_3" {
  name = "sg_3"
  description = "third test security group"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    from_group_id = "${huaweicloud_compute_secgroup_v2.sg_2.id}"
  }
}
`

const testAccComputeV2SecGroup_self = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    self = true
  }
}
`

const testAccComputeV2SecGroup_icmpZero = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 0
    to_port = 0
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }
}
`

const testAccComputeV2SecGroup_timeout = `
resource "huaweicloud_compute_secgroup_v2" "sg_1" {
  name = "sg_1"
  description = "first test security group"
  rule {
    from_port = 0
    to_port = 0
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }

  timeouts {
    delete = "5m"
  }
}
`
