package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/elb/v3/ipgroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccElbV3IpGroup_basic(t *testing.T) {
	var c ipgroups.IpGroup
	name := fmtp.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "huaweicloud_elb_ipgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElbV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3IpGroupConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3IpGroupExists(resourceName, &c),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "1"),
				),
			},
			{
				Config: testAccElbV3IpGroupConfig_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmtp.Sprintf("%s_updated", name)),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test updated"),
					resource.TestCheckResourceAttr(resourceName, "ip_list.#", "2"),
				),
			},
		},
	})
}

func TestAccElbV3IpGroup_withEpsId(t *testing.T) {
	var c ipgroups.IpGroup
	name := fmtp.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "huaweicloud_elb_ipgroup.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElbV3IpGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3IpGroupConfig_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3IpGroupExists(resourceName, &c),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckElbV3IpGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_elb_ipgroup" {
			continue
		}

		_, err := ipgroups.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IpGroup still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3IpGroupExists(
	n string, c *ipgroups.IpGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
		}

		found, err := ipgroups.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("IpGroup not found")
		}

		*c = *found

		return nil
	}
}

func testAccElbV3IpGroupConfig_basic(name string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_elb_ipgroup" "test"{
  name        = "%s"
  description = "terraform test"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }
}
`, name)
}

func testAccElbV3IpGroupConfig_update(name string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_elb_ipgroup" "test"{
  name        = "%s_updated"
  description = "terraform test updated"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }

  ip_list {
    ip = "192.168.10.11"
    description = "ECS02"
  }
}
`, name)
}

func testAccElbV3IpGroupConfig_withEpsId(name string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_elb_ipgroup" "test"{
  name        = "%s"
  description = "terraform test"

  ip_list {
    ip = "192.168.10.10"
    description = "ECS01"
  }

  enterprise_project_id = "%s"
}
`, name, HW_ENTERPRISE_PROJECT_ID_TEST)
}
