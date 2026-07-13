package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getAssetDomainLabelResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetAssetDomainLabelByName(client, state.Primary.Attributes["name"], state.Primary.Attributes["parent_id"])
}

// Before running the test, please ensure that you have created a DSC instance.
func TestAccAssetDomainLabel_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dsc_asset_domain_label.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getAssetDomainLabelResourceFunc)

		rNameWithChild = "huaweicloud_dsc_asset_domain_label.test2"
		rcWithChild    = acceptance.InitResourceCheck(rNameWithChild, &obj, getAssetDomainLabelResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithChild.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAssetDomainLabel_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "parent_id", "top label"),
					resource.TestCheckResourceAttr(rName, "depth", "1"),
					resource.TestCheckResourceAttr(rName, "is_leaf", "0"),
					rcWithChild.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNameWithChild, "parent_id", rName, "id"),
					resource.TestCheckResourceAttr(rNameWithChild, "depth", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAssetDomainLabelImportStateIdFunc(rName),
			},
			{
				ResourceName:      rNameWithChild,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAssetDomainLabelImportStateIdFunc(rNameWithChild),
			},
		},
	})
}

func testAccAssetDomainLabelImportStateIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", rName)
		}

		name := rs.Primary.Attributes["name"]
		parentId := rs.Primary.Attributes["parent_id"]
		if name == "" || parentId == "" {
			return "", fmt.Errorf("missing some attributes, want '<name>/<parent_id>', but got '%s/%s'", name, parentId)
		}

		return fmt.Sprintf("%s/%s", name, parentId), nil
	}
}

func testAssetDomainLabel_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_asset_domain_label" "test" {
  name      = "%[1]s"
  parent_id = "top label"
}

resource "huaweicloud_dsc_asset_domain_label" "test2" {
  name      = "%[1]s_child"
  parent_id = huaweicloud_dsc_asset_domain_label.test.id
}
`, name)
}
