package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetGroup(client, state.Primary.Attributes["component_id"], state.Primary.ID)
}

func TestAccResourceGroup_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	nameUpdate := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "sync_mode", "AUTO"),
					resource.TestCheckResourceAttr(resourceName, "region_id", "cn-north-4"),
					resource.TestCheckResourceAttrSet(resourceName, "code"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
				ImportStateIdFunc:       testGroupImportState(resourceName),
			},
			{
				Config: testAccGroup_updated(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdate),
					resource.TestCheckResourceAttr(resourceName, "sync_mode", "MANUAL"),
					resource.TestCheckResourceAttr(resourceName, "region_id", "cn-north-4"),
					resource.TestCheckResourceAttrSet(resourceName, "code"),
				),
			},
		},
	})
}

func testAccGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_group" "test" {
  name           = "%[2]s"
  component_id   = huaweicloud_coc_component.test.id
  region_id      = "cn-north-4"
  sync_mode      = "AUTO"
  vendor         = "RMS"
  application_id = huaweicloud_coc_application.test.id
  sync_rules {
    enterprise_project_id = "0"
    rule_tags             = jsonencode([{
        "key": "key1",
        "value": "value1"
    }])
  }
  force_delete = true
}
`, testAccComponent_basic(name), name)
}

func testAccGroup_updated(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_group" "test" {
  name           = "%[2]s"
  component_id   = huaweicloud_coc_component.test.id
  region_id      = "cn-north-4"
  sync_mode      = "MANUAL"
  vendor         = "RMS"
  application_id = huaweicloud_coc_application.test.id
  force_delete   = false
}
`, testAccComponent_basic(name), name)
}

func testGroupImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		componentID := rs.Primary.Attributes["component_id"]
		if componentID == "" {
			return "", fmt.Errorf("attribute (component_id) of resource (%s) not found", name)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", componentID, rs.Primary.ID), nil
	}
}
