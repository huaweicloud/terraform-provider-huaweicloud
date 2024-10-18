package dataarts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDataServiceCatalogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	httpUrl := "v1/{project_id}/service/servicecatalogs/{catalog_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{catalog_id}", state.Primary.ID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    state.Primary.Attributes["workspace_id"],
			"Dlm-Type":     state.Primary.Attributes["dlm_type"],
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", dataarts.CatalogResourceNotFoundCodes...)
	}

	return utils.FlattenResponse(requestResp)
}

func TestAccDataServiceCatalog_basic(t *testing.T) {
	var (
		obj interface{}

		rName       = "huaweicloud_dataarts_dataservice_catalog.test"
		name        = acceptance.RandomAccResourceName()
		updateName  = acceptance.RandomAccResourceName()
		basicConfig = testDataServiceCatalog_base()

		rc = acceptance.InitResourceCheck(rName, &obj, getDataServiceCatalogResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataServiceCatalog_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "parent_id", "0"),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "catalog_total", "0"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testDataServiceCatalog_basic_step2(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				// Refresh all resources and update the value of parameter 'catalog_total'.
				Config: testDataServiceCatalog_basic_step2(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_dataarts_dataservice_catalog.sub_path.0", "id"),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "catalog_total", "2"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testDataServiceCatalog_basic_step3(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				// Refresh all resources and update the value of parameter 'catalog_total'.
				Config: testDataServiceCatalog_basic_step3(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_dataarts_dataservice_catalog.sub_path.1", "id"),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(rName, "catalog_total", "3"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceCatalogImportState(rName),
				ImportStateVerifyIgnore: []string{
					"parent_id",
				},
			},
		},
	})
}

func testDataServiceCatalogImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, dlmType, resourceId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		dlmType = rs.Primary.Attributes["dlm_type"]
		resourceId = rs.Primary.ID
		if rs.Primary.Attributes["workspace_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute 'workspace_id' or resource ID is missing")
		}
		if rs.Primary.Attributes["dlm_type"] != "" {
			return fmt.Sprintf("%s/%s/%s", workspaceId, dlmType, resourceId), nil
		}
		return fmt.Sprintf("%s/%s", workspaceId, resourceId), nil
	}
}

func testDataServiceCatalog_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
// Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "sub_path" {
  count = 2

  workspace_id = "%[1]s"
  dlm_type     = "SHARED"
  name         = format("%[2]s_parent_%%d", count.index)
}

`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataServiceCatalog_basic_step1(name string) string {
	return fmt.Sprintf(`
// Under root path.
resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  workspace_id = "%[1]s"
  dlm_type     = "SHARED"
  name         = "%[2]s"
  description  = "Created by terraform script"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataServiceCatalog_basic_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  parent_id    = huaweicloud_dataarts_dataservice_catalog.sub_path[0].id
  workspace_id = "%[2]s"
  dlm_type     = "SHARED"
  name         = "%[3]s"
  description  = "Created by terraform script"
}

resource "huaweicloud_dataarts_dataservice_catalog" "child" {
  count = 2

  parent_id    = huaweicloud_dataarts_dataservice_catalog.test.id
  workspace_id = "%[2]s"
  dlm_type     = "SHARED"
  name         = format("%[3]s_child_%%d", count.index)
}
`, basicConfig, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataServiceCatalog_basic_step3(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_catalog" "test" {
  parent_id = huaweicloud_dataarts_dataservice_catalog.sub_path[1].id

  workspace_id = "%[2]s"
  dlm_type     = "SHARED"
  name         = "%[3]s"
  description  = "Updated by terraform script"
}

resource "huaweicloud_dataarts_dataservice_catalog" "child" {
  count = 3

  parent_id    = huaweicloud_dataarts_dataservice_catalog.test.id
  workspace_id = "%[2]s"
  dlm_type     = "SHARED"
  name         = format("%[3]s_child_%%d", count.index)
}
`, basicConfig, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
