package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
)

func getHostGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		epsId   = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		groupId = state.Primary.ID
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.QueryHostGroupById(client, region, epsId, groupId)
}

func TestAccHostGroup_basic(t *testing.T) {
	var (
		hostGroup interface{}
		name      = acceptance.RandomAccResourceName()
		rName     = "huaweicloud_hss_host_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&hostGroup,
		getHostGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccHostGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
				),
			},
			{
				Config: testAccHostGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHostGroupImportStateIDFunc(rName),
				// The field `unprotect_host_ids` will be filled in during the creation and editing operations.
				// We only need to add ignore to the test case and do not need to make special instructions in the document.
				ImportStateVerifyIgnore: []string{
					"unprotect_host_ids",
				},
			},
		},
	})
}

func testAccHostGroup_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name                  = "%[2]s_${count.index}"
  image_id              = data.huaweicloud_images_image.test.id
  flavor_id             = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups       = [huaweicloud_networking_secgroup.test.name]
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = "%[3]s"

  key_pair   = huaweicloud_kps_keypair.test.name
  agent_list = "hss"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseComputeResources(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHostGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_host_group" "test" {
  name                  = "%[2]s"
  host_ids              = slice(huaweicloud_compute_instance.test[*].id, 0, 1)
  enterprise_project_id = "%[3]s"
}
`, testAccHostGroup_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHostGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_host_group" "test" {
  name                  = "%[2]s-update"
  host_ids              = huaweicloud_compute_instance.test[*].id
  enterprise_project_id = "%[3]s"
}
`, testAccHostGroup_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHostGroupImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		id := rs.Primary.ID
		if epsId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<enterprise_project_id>/<id>', but got '%s/%s'",
				epsId, id)
		}
		return fmt.Sprintf("%s/%s", epsId, id), nil
	}
}
