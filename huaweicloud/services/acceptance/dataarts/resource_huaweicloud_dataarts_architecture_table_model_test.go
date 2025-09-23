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

func getArchitectureTableModelResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getTableModelClient, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio Client: %s", err)
	}
	getTableModelHttpUrl := "v2/{project_id}/design/table-model/{id}?latest=true"
	getTableModelPath := getTableModelClient.Endpoint + getTableModelHttpUrl
	getTableModelPath = strings.ReplaceAll(getTableModelPath, "{project_id}", getTableModelClient.ProjectID)
	getTableModelPath = strings.ReplaceAll(getTableModelPath, "{id}", state.Primary.ID)
	getTableModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}
	getTableModelResp, err := getTableModelClient.Request("GET", getTableModelPath, &getTableModelOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture table model: %s", err)
	}
	return utils.FlattenResponse(getTableModelResp)
}

func TestAccResourceArchitectureTableModel_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_architecture_table_model.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getArchitectureTableModelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsSubjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureTableModel_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "physical_table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "create"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "key"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name_en", "key_en"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.data_type", "INTEGER"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.ordinal", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.name", "key2"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.name_en", "key2_en"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.data_type", "INTEGER"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.ordinal", "2"),
					resource.TestCheckResourceAttr(resourceName, "has_related_logic_table", "false"),
					resource.TestCheckResourceAttr(resourceName, "has_related_physical_table", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_partition", "false"),
					resource.TestCheckResourceAttr(resourceName, "compression", "NO"),
					resource.TestCheckResourceAttr(resourceName, "table_type", "DWS_ROW"),
					resource.TestCheckResourceAttr(resourceName, "reversed", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "catalog_path"),
					resource.TestCheckResourceAttrSet(resourceName, "env_type"),
					resource.TestCheckResourceAttrSet(resourceName, "extend_info"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "table_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccArchitectureTableModel_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "physical_table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "table_name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "update"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name", "key"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.name_en", "key_en"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.data_type", "INTEGER"),
					resource.TestCheckResourceAttr(resourceName, "attributes.0.ordinal", "1"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.name", "key2"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.name_en", "key2_en"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.data_type", "INTEGER"),
					resource.TestCheckResourceAttr(resourceName, "attributes.1.ordinal", "2"),
					resource.TestCheckResourceAttr(resourceName, "has_related_logic_table", "false"),
					resource.TestCheckResourceAttr(resourceName, "has_related_physical_table", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_partition", "false"),
					resource.TestCheckResourceAttr(resourceName, "compression", "NO"),
					resource.TestCheckResourceAttr(resourceName, "table_type", "DWS_ROW"),
					resource.TestCheckResourceAttrSet(resourceName, "catalog_path"),
					resource.TestCheckResourceAttrSet(resourceName, "env_type"),
					resource.TestCheckResourceAttrSet(resourceName, "extend_info"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "table_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataArtsStudioImportState(resourceName),
			},
		},
	})
}

func testAccArchitectureTableModel_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_table_model" "test" {
  workspace_id        = "%[2]s"
  model_id            = huaweicloud_dataarts_architecture_model.test.id
  subject_id          = "%[3]s"
  physical_table_name = "%[4]s"
  table_name          = "%[4]s"
  description         = "create"
  dw_type             = huaweicloud_dataarts_architecture_model.test.dw_type
  compression         = "NO"
  table_type          = "DWS_ROW"
  reversed            = false

  attributes {
    name      = "key"
    name_en   = "key_en"
    data_type = "INTEGER"
    ordinal   = "1"
  }
  attributes {
    name      = "key2"
    name_en   = "key2_en"
    data_type = "INTEGER"
    ordinal   = "2"
  }
}
`, testAccModel_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_SUBJECT_ID, name)
}

func testAccArchitectureTableModel_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_table_model" "test" {
  workspace_id        = "%[2]s"
  model_id            = huaweicloud_dataarts_architecture_model.test.id
  subject_id          = "%[3]s"
  physical_table_name = "%[4]s"
  table_name          = "%[4]s"
  description         = "update"
  dw_type             = huaweicloud_dataarts_architecture_model.test.dw_type
  compression         = "NO"
  table_type          = "DWS_ROW"
  reversed            = false

  attributes {
    name      = "key"
    name_en   = "key_en"
    data_type = "INTEGER"
    ordinal   = "1"
  }
  attributes {
    name      = "key2"
    name_en   = "key2_en"
    data_type = "INTEGER"
    ordinal   = "2"
  }
}
`, testAccModel_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_SUBJECT_ID, name)
}
