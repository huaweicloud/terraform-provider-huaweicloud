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

func getBusinessMetricResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/design/biz-metrics/{id}"
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
		return nil, fmt.Errorf("error retrieving DataArts Architecture business metric: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccBusinessMetric_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_architecture_business_metric.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBusinessMetricResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsBizCatalogID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBusinessMetric_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "biz_catalog_id", acceptance.HW_DATAARTS_BIZ_CATALOG_ID),
					resource.TestCheckResourceAttr(rName, "owner", "test-owner"),
					resource.TestCheckResourceAttr(rName, "owner_department", "test-department"),
					resource.TestCheckResourceAttr(rName, "time_filters", "双周"),
					resource.TestCheckResourceAttr(rName, "interval_type", "HOUR"),
					resource.TestCheckResourceAttr(rName, "destination", "test destination"),
					resource.TestCheckResourceAttr(rName, "definition", "test definition"),
					resource.TestCheckResourceAttr(rName, "expression", "a+b+c"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "apply_scenario", "test apply scenario"),
					resource.TestCheckResourceAttr(rName, "technical_metric", "100"),
					resource.TestCheckResourceAttr(rName, "measure", "test measure"),
					resource.TestCheckResourceAttr(rName, "general_filters", "test general filters"),
					resource.TestCheckResourceAttr(rName, "data_origin", "test data origin"),
					resource.TestCheckResourceAttr(rName, "unit", "test unit"),
					resource.TestCheckResourceAttr(rName, "name_alias", "alias-name"),
					resource.TestCheckResourceAttr(rName, "code", "testCode"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "biz_catalog_path"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "l1"),
				),
			},
			{
				Config: testBusinessMetric_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "biz_catalog_id", acceptance.HW_DATAARTS_BIZ_CATALOG_ID),
					resource.TestCheckResourceAttr(rName, "owner", "test-owner-update"),
					resource.TestCheckResourceAttr(rName, "owner_department", "test-department-update"),
					resource.TestCheckResourceAttr(rName, "time_filters", "周,季度,其他"),
					resource.TestCheckResourceAttr(rName, "interval_type", "WEEK"),
					resource.TestCheckResourceAttr(rName, "destination", "test destination update"),
					resource.TestCheckResourceAttr(rName, "definition", "test definition update"),
					resource.TestCheckResourceAttr(rName, "expression", "a+b+c+d"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "apply_scenario", "test apply scenario update"),
					resource.TestCheckResourceAttr(rName, "technical_metric", "200"),
					resource.TestCheckResourceAttr(rName, "measure", "test measure update"),
					resource.TestCheckResourceAttr(rName, "general_filters", "test general filters update"),
					resource.TestCheckResourceAttr(rName, "data_origin", "test data origin update"),
					resource.TestCheckResourceAttr(rName, "unit", "test unit update"),
					resource.TestCheckResourceAttr(rName, "name_alias", "alias-name-update"),
					resource.TestCheckResourceAttr(rName, "code", "testCodeUpdate"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "biz_catalog_path"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "l1"),
				),
			},
			{
				Config: testBusinessMetric_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "biz_catalog_id", acceptance.HW_DATAARTS_BIZ_CATALOG_ID),
					resource.TestCheckResourceAttr(rName, "owner", "test-owner-update"),
					resource.TestCheckResourceAttr(rName, "owner_department", "test-department-update"),
					resource.TestCheckResourceAttr(rName, "time_filters", "周,季度,其他"),
					resource.TestCheckResourceAttr(rName, "interval_type", "WEEK"),
					resource.TestCheckResourceAttr(rName, "destination", "test destination update"),
					resource.TestCheckResourceAttr(rName, "definition", "test definition update"),
					resource.TestCheckResourceAttr(rName, "expression", "a+b+c+d"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "apply_scenario", ""),
					resource.TestCheckResourceAttr(rName, "technical_metric", "0"),
					resource.TestCheckResourceAttr(rName, "measure", ""),
					resource.TestCheckResourceAttr(rName, "general_filters", ""),
					resource.TestCheckResourceAttr(rName, "data_origin", ""),
					resource.TestCheckResourceAttr(rName, "unit", ""),
					resource.TestCheckResourceAttr(rName, "name_alias", ""),
					resource.TestCheckResourceAttr(rName, "code", "testCodeUpdate"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "biz_catalog_path"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "l1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataArtsStudioImportState(rName),
			},
		},
	})
}

func testBusinessMetric_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  name             = "%[1]s"
  workspace_id     = "%[2]s"
  biz_catalog_id   = "%[3]s"
  owner            = "test-owner"
  owner_department = "test-department"
  time_filters     = "双周"
  interval_type    = "HOUR"
  destination      = "test destination"
  definition       = "test definition"
  expression       = "a+b+c"
  description      = "test description"
  apply_scenario   = "test apply scenario"
  technical_metric = 100
  measure          = "test measure"
  general_filters  = "test general filters"
  data_origin      = "test data origin"
  unit             = "test unit"
  name_alias       = "alias-name"
  code             = "testCode"
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_BIZ_CATALOG_ID)
}

func testBusinessMetric_basic_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  name             = "%[1]s"
  workspace_id     = "%[2]s"
  biz_catalog_id   = "%[3]s"
  owner            = "test-owner-update"
  owner_department = "test-department-update"
  time_filters     = "周,季度,其他"
  interval_type    = "WEEK"
  destination      = "test destination update"
  definition       = "test definition update"
  expression       = "a+b+c+d"
  description      = "test description update"
  apply_scenario   = "test apply scenario update"
  technical_metric = 200
  measure          = "test measure update"
  general_filters  = "test general filters update"
  data_origin      = "test data origin update"
  unit             = "test unit update"
  name_alias       = "alias-name-update"
  code             = "testCodeUpdate"
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_BIZ_CATALOG_ID)
}

func testBusinessMetric_basic_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  name             = "%[1]s"
  workspace_id     = "%[2]s"
  biz_catalog_id   = "%[3]s"
  owner            = "test-owner-update"
  owner_department = "test-department-update"
  time_filters     = "周,季度,其他"
  interval_type    = "WEEK"
  destination      = "test destination update"
  definition       = "test definition update"
  expression       = "a+b+c+d"
  code             = "testCodeUpdate"
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_BIZ_CATALOG_ID)
}

// testDataArtsStudioImportState use to return an id with format <workspace_id>/<id>
func testDataArtsStudioImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		if workspaceID == "" {
			return "", fmt.Errorf("attribute (workspace_id) of Resource (%s) not found", name)
		}

		return workspaceID + "/" + rs.Primary.ID, nil
	}
}
