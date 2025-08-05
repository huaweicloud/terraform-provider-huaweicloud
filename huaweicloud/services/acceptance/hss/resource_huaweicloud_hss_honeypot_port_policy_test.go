package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
)

func getResourceHoneypotPortPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("hss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.GetHoneypotPortPolicy(client, state.Primary.ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccResourceHoneypotPortPolicy_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_hss_honeypot_port_policy.test"
		name         = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			resourceName,
			&object,
			getResourceHoneypotPortPolicyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccHoneypotPortPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_name", name),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "ports_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ports_list.0.port", "8002"),
					resource.TestCheckResourceAttr(resourceName, "ports_list.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "host_id.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "white_ip.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "group_list.#", "1"),
				),
			},
			{
				Config: testAccHoneypotPortPolicy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_name", name+"_update"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "ports_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ports_list.0.port", "8008"),
					resource.TestCheckResourceAttr(resourceName, "ports_list.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "host_id.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "white_ip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "group_list.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ports_list",
					"host_id",
					"group_list",
					"enterprise_project_id",
				},
			},
		},
	})
}

func testAccHoneypotPortPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_host_group" "test" {
  name                  = "%[1]s"
  host_ids              = ["%[2]s"]
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_hss_honeypot_port_policy" "test" {
  policy_name = "%[1]s" 
  os_type     = "Linux"

  ports_list {
    port     = 8002
    protocol = "tcp"
  }

  ports_list {
    port     = 8006
    protocol = "tcp"
  }

  host_id               = ["%[2]s"]
  white_ip              = ["192.168.1.24","192.168.1.26"]
  group_list            = [huaweicloud_hss_host_group.test.id]
  enterprise_project_id = "%[3]s"
}
`, name, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHoneypotPortPolicy_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_host_group" "test" {
  name                  = "%[1]s"
  host_ids              = ["%[2]s"]
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_hss_honeypot_port_policy" "test" {
  policy_name = "%[1]s_update" 
  os_type     = "Linux"

  ports_list {
    port     = 8008
    protocol = "tcp"
  }

  host_id               = ["%[2]s"]
  white_ip              = ["192.168.1.28"]
  group_list            = [huaweicloud_hss_host_group.test.id]
  enterprise_project_id = "%[3]s"
}
`, name, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
