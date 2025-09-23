package secmaster

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

func getDataObjectRelationsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getDataObjectRelationsHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/{dataclass_type}/{data_object_id}/{related_dataclass_type}/search"
	getDataObjectRelationsPath := client.Endpoint + getDataObjectRelationsHttpUrl
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath, "{project_id}", client.ProjectID)
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{workspace_id}", state.Primary.Attributes["workspace_id"])
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{dataclass_type}", state.Primary.Attributes["data_class"])
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{data_object_id}", state.Primary.Attributes["data_object_id"])
	getDataObjectRelationsPath = strings.ReplaceAll(getDataObjectRelationsPath,
		"{related_dataclass_type}", state.Primary.Attributes["related_data_class"])

	getDataObjectRelationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getDataObjectRelationsOpt.JSONBody = map[string]interface{}{
		"limit":  1000,
		"offset": 0,
	}

	getDataObjectRelationsResp, err := client.Request("POST", getDataObjectRelationsPath, &getDataObjectRelationsOpt)
	if err != nil {
		return nil, err
	}

	getDataObjectRelationsRespBody, err := utils.FlattenResponse(getDataObjectRelationsResp)
	if err != nil {
		return nil, err
	}

	relatedDataObjectIds := utils.PathSearch("data[*].data_object.id",
		getDataObjectRelationsRespBody, make([]interface{}, 0)).([]interface{})
	if len(relatedDataObjectIds) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return relatedDataObjectIds, nil
}

func TestAccDataObjectRelations_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_data_object_relations.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataObjectRelationsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataObjectRelations_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "related_data_object_ids.0",
						"huaweicloud_secmaster_indicator.test_1", "id"),
				),
			},
			{
				Config: testAccDataObjectRelations_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "related_data_object_ids.0",
						"huaweicloud_secmaster_indicator.test_2", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataObjectRelationsImportState(rName),
			},
		},
	})
}

func testDataObjectRelationsImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		workspaceID := rs.Primary.Attributes["workspace_id"]
		dataClass := rs.Primary.Attributes["data_class"]
		dataObjectID := rs.Primary.Attributes["data_object_id"]
		relatedDataClass := rs.Primary.Attributes["related_data_class"]
		if workspaceID == "" || dataClass == "" || dataObjectID == "" || relatedDataClass == "" {
			return "", fmt.Errorf("invalid format specified for import ID, " +
				"want '<workspace_id>/<data_class>/<data_object_id>/<related_data_class>'")
		}

		return fmt.Sprintf("%s/%s/%s/%s", workspaceID, dataClass, dataObjectID, relatedDataClass), nil
	}
}

func testAccDataObjectRelations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_data_object_relations" "test" {
  workspace_id            = "%[2]s"
  data_class              = "incidents"
  data_object_id          = huaweicloud_secmaster_incident.test.id
  related_data_class      = "indicators"
  related_data_object_ids = [huaweicloud_secmaster_indicator.test_1.id]
}
`, testAccDataObjectRelations_build(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testAccDataObjectRelations_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_data_object_relations" "test" {
  workspace_id            = "%[2]s"
  data_class              = "incidents"
  data_object_id          = huaweicloud_secmaster_incident.test.id
  related_data_class      = "indicators"
  related_data_object_ids = [huaweicloud_secmaster_indicator.test_2.id]
}
`, testAccDataObjectRelations_build(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testAccDataObjectRelations_build(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_indicator" "test_1" {
  workspace_id = "%[2]s"
  name         = "%[3]s_1"
  
  type {
    category       = "Domain"
    indicator_type = "Domain"
    id             = "%[4]s"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Closed"
  confidence            = "90"
  first_occurrence_time = "2024-07-26T09:33:55.000+08:00"
  last_occurrence_time  = "2023-07-27T21:15:30.000+08:00"
  threat_degree         = "Gray"
  granularity           = "1"
  value                 = "test.terraform.com"
}

resource "huaweicloud_secmaster_indicator" "test_2" {
  workspace_id = "%[2]s"
  name         = "%[3]s_2"
  
  type {
    category       = "Domain"
    indicator_type = "Domain"
    id             = "%[4]s"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Closed"
  confidence            = "90"
  first_occurrence_time = "2024-07-26T09:33:55.000+08:00"
  last_occurrence_time  = "2023-07-27T21:15:30.000+08:00"
  threat_degree         = "Gray"
  granularity           = "1"
  value                 = "test2.terraform.com"
}
`, testIncident_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID, name, acceptance.HW_SECMASTER_INDICATOR_TYPE_ID)
}
