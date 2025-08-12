package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getCatalogueResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getCatalogueClient, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	workspaceID := state.Primary.Attributes["workspace_id"]
	return secmaster.ReadCatalogueDetail(getCatalogueClient, workspaceID, state.Primary.ID)
}

func TestAccCatalogue_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_secmaster_catalogue.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCatalogueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterLayoutID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCatalogue_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "layout_id", acceptance.HW_SECMASTER_LAYOUT_ID),
					resource.TestCheckResourceAttr(rName, "catalogue_address", "/ssa/workspace/soc/situation/report/detail1"),
					resource.TestCheckResourceAttr(rName, "parent_alias_en", "test-first-block1"),
					resource.TestCheckResourceAttr(rName, "parent_alias_zh", "测试一级目录1"),
					resource.TestCheckResourceAttr(rName, "parent_catalogue", "demo-first-name1"),
					resource.TestCheckResourceAttr(rName, "second_alias_en", "test-second-block1"),
					resource.TestCheckResourceAttr(rName, "second_alias_zh", "测试二级目录1"),
					resource.TestCheckResourceAttr(rName, "second_catalogue", "demo-second-name1"),
					resource.TestCheckResourceAttr(rName, "second_catalogue_code", "SECURITY_REPORT"),
					resource.TestCheckResourceAttrSet(rName, "catalogue_status"),
					resource.TestCheckResourceAttrSet(rName, "is_card_area"),
					resource.TestCheckResourceAttrSet(rName, "is_display"),
					resource.TestCheckResourceAttrSet(rName, "is_landing_page"),
					resource.TestCheckResourceAttrSet(rName, "is_navigation"),
					resource.TestCheckResourceAttrSet(rName, "layout_name"),
					resource.TestCheckResourceAttrSet(rName, "publisher_name"),
				),
			},
			{
				Config: testCatalogue_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "layout_id", acceptance.HW_SECMASTER_LAYOUT_ID),
					resource.TestCheckResourceAttr(rName, "catalogue_address", "/ssa/workspace/soc/situation/report/detail1"),
					resource.TestCheckResourceAttr(rName, "parent_alias_en", "test-first-block1"),
					resource.TestCheckResourceAttr(rName, "parent_alias_zh", "测试一级目录1"),
					resource.TestCheckResourceAttr(rName, "parent_catalogue", "demo-first-name1"),
					resource.TestCheckResourceAttr(rName, "second_alias_en", "test-second-block2"),
					resource.TestCheckResourceAttr(rName, "second_alias_zh", "测试二级目录2"),
					resource.TestCheckResourceAttr(rName, "second_catalogue", "demo-second-name2"),
					resource.TestCheckResourceAttr(rName, "second_catalogue_code", "SECURITY_REPORT"),
					resource.TestCheckResourceAttrSet(rName, "catalogue_status"),
					resource.TestCheckResourceAttrSet(rName, "is_card_area"),
					resource.TestCheckResourceAttrSet(rName, "is_display"),
					resource.TestCheckResourceAttrSet(rName, "is_landing_page"),
					resource.TestCheckResourceAttrSet(rName, "is_navigation"),
					resource.TestCheckResourceAttrSet(rName, "layout_name"),
					resource.TestCheckResourceAttrSet(rName, "publisher_name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testCatalogueImportState(rName),
				ImportStateVerifyIgnore: []string{"second_catalogue_code"},
			},
		},
	})
}

func testCatalogue_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_catalogue" "test" {
  workspace_id          = "%s"
  layout_id             = "%s"
  parent_catalogue      = "demo-first-name1"
  parent_alias_zh       = "测试一级目录1"
  parent_alias_en       = "test-first-block1"
  second_catalogue      = "demo-second-name1"
  second_alias_zh       = "测试二级目录1"
  second_alias_en       = "test-second-block1"
  second_catalogue_code = "SECURITY_REPORT"
  catalogue_address     = "/ssa/workspace/soc/situation/report/detail1"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_LAYOUT_ID)
}

func testCatalogue_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_catalogue" "test" {
  workspace_id          = "%s"
  layout_id             = "%s"
  parent_catalogue      = "demo-first-name1"
  parent_alias_zh       = "测试一级目录1"
  parent_alias_en       = "test-first-block1"
  second_catalogue      = "demo-second-name2"
  second_alias_zh       = "测试二级目录2"
  second_alias_en       = "test-second-block2"
  second_catalogue_code = "SECURITY_REPORT"
  catalogue_address     = "/ssa/workspace/soc/situation/report/detail1"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_LAYOUT_ID)
}

func testCatalogueImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}
