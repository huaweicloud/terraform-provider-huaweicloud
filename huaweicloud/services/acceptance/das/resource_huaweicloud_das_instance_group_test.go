package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/das"
)

func getInstanceGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}

	groupId := state.Primary.ID
	datastoreType := state.Primary.Attributes["datastore_type"]

	return das.GetInstanceGroupById(client, datastoreType, groupId)
}

func TestAccInstanceGroup_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_das_instance_group.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getInstanceGroupResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "datastore_type", "MySQL"),
					resource.TestCheckResourceAttr(rName, "group_name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
				),
			},
			{
				Config: testAccInstanceGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_name", name+"-updated"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceGroupImportIdFunc(rName),
			},
		},
	})
}

func testAccInstanceGroupImportIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["datastore_type"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["datastore_type"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["datastore_type"], rs.Primary.ID), nil
	}
}

func testAccInstanceGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_das_instance_group" "test" {
  datastore_type = "MySQL"
  group_name     = "%[1]s"
  description    = "Created by terraform script"
}
`, name)
}

func testAccInstanceGroup_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_das_instance_group" "test" {
  datastore_type = "MySQL"
  group_name     = "%[1]s-updated"
  description    = "Updated by terraform script"
}
`, name)
}
