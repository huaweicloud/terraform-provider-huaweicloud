package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccUniAgent_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_uniagent.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAOMUniAgentAgentID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testUniAgent_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version", "1.1.6"),
				),
			},
			{
				Config: testUniAgent_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version", "1.1.5"),
				),
			},
		},
	})
}

func testUniAgent_ECS_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name                        = "%[2]s"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  admin_pass                  = "Terraform@123"
  delete_disks_on_termination = true
  security_group_ids          = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}`, common.TestBaseComputeResources(name), name)
}

func testUniAgent_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_uniagent" "test" {
  installer_agent_id = "%[2]s"
  version            = "1.1.6"
  public_net_flag    = false
  proxy_region_id    = 0
  inner_ip           = huaweicloud_compute_instance.test.access_ip_v4
  port               = 22
  account            = "root"
  password           = "Terraform@123"
  os_type            = "LINUX"
}`, testUniAgent_ECS_base(name), acceptance.HW_AOM_UNIAGENT_AGENT_ID)
}

func testUniAgent_update(name string) string {
	return fmt.Sprintf(`
%[1]s

# enter installer_agent_id to agent_id, just for test
# when update version, API do not check for agent_id, and response will be success, but actually do nothing
resource "huaweicloud_aom_uniagent" "test" {
  installer_agent_id = "%[2]s"
  version            = "1.1.5"
  public_net_flag    = false
  proxy_region_id    = 0
  agent_id           = "%[2]s"
  inner_ip           = huaweicloud_compute_instance.test.access_ip_v4
  port               = 22
  account            = "root"
  password           = "Terraform@123"
  os_type            = "LINUX"
}`, testUniAgent_ECS_base(name), acceptance.HW_AOM_UNIAGENT_AGENT_ID)
}
