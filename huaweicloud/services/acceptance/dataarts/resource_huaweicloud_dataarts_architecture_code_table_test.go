package dataarts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getArchitectureCodeTableFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/design/code-tables/{id}"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture code table: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccArchitectureCodeTable_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_architecture_code_table.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getArchitectureCodeTableFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsCodeTableDirectoryID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testArchitectureCodeTable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "code", fmt.Sprintf("%s_code", name)),
					resource.TestCheckResourceAttr(rName, "directory_id", acceptance.HW_DATAARTS_CODE_TABLE_DIRECTORY_ID),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "fields.0.name", "field"),
					resource.TestCheckResourceAttr(rName, "fields.0.code", "field_code"),
					resource.TestCheckResourceAttr(rName, "fields.0.type", "BIGINT"),
					resource.TestCheckResourceAttr(rName, "fields.0.description", "test field description"),
				),
			},
			{
				Config: testArchitectureCodeTable_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "code", fmt.Sprintf("%s_code_update", name)),
					resource.TestCheckResourceAttr(rName, "directory_id", acceptance.HW_DATAARTS_CODE_TABLE_DIRECTORY_ID),
					resource.TestCheckResourceAttr(rName, "description", "test update description"),
					resource.TestCheckResourceAttr(rName, "fields.0.name", "field_update"),
					resource.TestCheckResourceAttr(rName, "fields.0.code", "field_code_update"),
					resource.TestCheckResourceAttr(rName, "fields.0.type", "STRING"),
					resource.TestCheckResourceAttr(rName, "fields.0.description", "test field update description"),
					resource.TestCheckResourceAttr(rName, "fields.0.values.0", "value1"),
					resource.TestCheckResourceAttr(rName, "fields.0.values.1", "value2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testArchitectureCodeTableImportState(rName),
			},
		},
	})
}

func testArchitectureCodeTable_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_code_table" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  code         = "%[2]s_code"
  directory_id = "%[3]s"
  description  = "test description"

  fields {
    name        = "field"
    code        = "field_code"
    type        = "BIGINT"
    description = "test field description"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_CODE_TABLE_DIRECTORY_ID)
}

func testArchitectureCodeTable_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_code_table" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_update"
  code         = "%[2]s_code_update"
  directory_id = "%[3]s"
  description  = "test update description"

  fields {
    name        = "field_update"
    code        = "field_code_update"
    type        = "STRING"
    description = "test field update description"
    values      = ["value1", "value2"]
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_CODE_TABLE_DIRECTORY_ID)
}

func testArchitectureCodeTableImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		if workspaceID == "" {
			return "", fmt.Errorf("attribute (workspace_id) of Resource (%s) not found", name)
		}

		return workspaceID + "/" + rs.Primary.Attributes["name"], nil
	}
}
