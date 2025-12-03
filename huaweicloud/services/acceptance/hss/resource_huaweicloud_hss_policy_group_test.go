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

func getPolicyGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region   = acceptance.HW_REGION_NAME
		epsId    = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		policyId = state.Primary.ID
		product  = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.QueryPolicyGroupById(client, policyId, epsId)
}

func TestAccPolicyGroup_basic(t *testing.T) {
	var (
		policy interface{}
		name   = acceptance.RandomAccResourceName()
		rName  = "huaweicloud_hss_policy_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&policy,
		getPolicyGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSPolicyGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_id", acceptance.HW_HSS_POLICY_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "protect_mode", "high_detection"),
					resource.TestCheckResourceAttrSet(rName, "default_group"),
					resource.TestCheckResourceAttrSet(rName, "deletable"),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
					resource.TestCheckResourceAttrSet(rName, "support_os"),
					resource.TestCheckResourceAttrSet(rName, "support_version"),
				),
			},
			{
				Config: testAccPolicyGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "group_id", acceptance.HW_HSS_POLICY_GROUP_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "protect_mode", "equalization"),
					resource.TestCheckResourceAttrSet(rName, "default_group"),
					resource.TestCheckResourceAttrSet(rName, "deletable"),
					resource.TestCheckResourceAttrSet(rName, "host_num"),
					resource.TestCheckResourceAttrSet(rName, "support_os"),
					resource.TestCheckResourceAttrSet(rName, "support_version"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPolicyGroupImportStateIDFunc(rName),
				// The following fields will be ignored during import verification
				ImportStateVerifyIgnore: []string{
					"group_id",
				},
			},
		},
	})
}

func testAccPolicyGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_policy_group" "test" {
  group_id              = "%s"
  name                  = "%s"
  description           = "test description"
  enterprise_project_id = "all_granted_eps"
  protect_mode          = "high_detection"
}
`, acceptance.HW_HSS_POLICY_GROUP_ID, name)
}

func testAccPolicyGroup_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_policy_group" "test" {
  group_id              = "%s"
  name                  = "%s"
  description           = "test description"
  enterprise_project_id = "all_granted_eps"
  protect_mode          = "equalization"
}
`, acceptance.HW_HSS_POLICY_GROUP_ID, name)
}

func testAccPolicyGroupImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resourceName)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		id := rs.Primary.ID
		if epsId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<enterprise_project_id>/<id>', but got '%s/%s'", epsId, id)
		}
		return fmt.Sprintf("%s/%s", epsId, id), nil
	}
}
