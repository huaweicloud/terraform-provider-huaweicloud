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

func getArchitectureCodeTableValuesFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/design/code-tables/{id}/values"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.Attributes["table_id"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture code table values: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	fieldName := state.Primary.Attributes["field_name"]
	findFieldValuesExpr := fmt.Sprintf("data.value.records[?name_ch=='%s'] | [0].code_table_field_values", fieldName)
	curJson := utils.PathSearch(findFieldValuesExpr, getRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	// After the test case is executed, the values of the all fields are cleared.
	// There is no need to check if id of the value is nil.
	if len(curArray) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccArchitectureCodeTableValues_basic(t *testing.T) {
	var obj interface{}

	baseConfig := testArchitectureCodeTableValuesConfig_base()
	rName := "huaweicloud_dataarts_architecture_code_table_values.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getArchitectureCodeTableValuesFunc,
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
				Config: testArchitectureCodeTableValues_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "values.0.value", "1"),
					resource.TestCheckResourceAttrSet(rName, "table_id"),
					resource.TestCheckResourceAttrSet(rName, "field_id"),
					resource.TestCheckResourceAttrSet(rName, "field_ordinal"),
					resource.TestCheckResourceAttrSet(rName, "field_name"),
					resource.TestCheckResourceAttrSet(rName, "field_code"),
					resource.TestCheckResourceAttrSet(rName, "field_type"),
				),
			},
			{
				Config: testArchitectureCodeTableValues_basic_update(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "values.0.value", "one"),
					resource.TestCheckResourceAttr(rName, "values.1.value", "two"),
					resource.TestCheckResourceAttrSet(rName, "table_id"),
					resource.TestCheckResourceAttrSet(rName, "field_id"),
					resource.TestCheckResourceAttrSet(rName, "field_ordinal"),
					resource.TestCheckResourceAttrSet(rName, "field_name"),
					resource.TestCheckResourceAttrSet(rName, "field_code"),
					resource.TestCheckResourceAttrSet(rName, "field_type"),
				),
			},
		},
	})
}

func testArchitectureCodeTableValuesConfig_base() string {
	name := acceptance.RandomAccResourceName()
	dirName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_directory" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "CODE"
}

resource "huaweicloud_dataarts_architecture_code_table" "test" {
  workspace_id = "%[1]s"
  name         = "%[3]s"
  code         = "%[3]s_code"
  directory_id = huaweicloud_dataarts_architecture_directory.test.id
  description  = "test description"
  
  fields {
    name        = "field"
    code        = "field_code"
    type        = "BIGINT"
    description = "test field description"
  }
  
  fields {
    name        = "field1"
    code        = "field_code1"
    type        = "STRING"
    description = "test field1 description"
  }
  
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, dirName, name)
}

func testArchitectureCodeTableValues_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_dataarts_architecture_code_table_values" "test" {
  workspace_id  = "%[2]s"
  table_id      = huaweicloud_dataarts_architecture_code_table.test.id
  field_id      = huaweicloud_dataarts_architecture_code_table.test.fields[1].id
  field_ordinal = huaweicloud_dataarts_architecture_code_table.test.fields[1].ordinal
  field_name    = huaweicloud_dataarts_architecture_code_table.test.fields[1].name
  field_code    = huaweicloud_dataarts_architecture_code_table.test.fields[1].code
  field_type    = huaweicloud_dataarts_architecture_code_table.test.fields[1].type
  
  values {
    value = "1"
  }
}
`, baseConfig, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testArchitectureCodeTableValues_basic_update(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_dataarts_architecture_code_table_values" "test" {
  workspace_id  = "%[2]s"
  table_id      = huaweicloud_dataarts_architecture_code_table.test.id
  field_id      = huaweicloud_dataarts_architecture_code_table.test.fields[1].id
  field_ordinal = huaweicloud_dataarts_architecture_code_table.test.fields[1].ordinal
  field_name    = huaweicloud_dataarts_architecture_code_table.test.fields[1].name
  field_code    = huaweicloud_dataarts_architecture_code_table.test.fields[1].code
  field_type    = huaweicloud_dataarts_architecture_code_table.test.fields[1].type
  
  values {
    value = "one"
  }
  values {
    value = "two"
  }
}
`, baseConfig, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
