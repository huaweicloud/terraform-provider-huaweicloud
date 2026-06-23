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

func getResourceMetricFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetMetricInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceMetric_basic(t *testing.T) {
	var (
		rName      = "huaweicloud_secmaster_metric.test"
		metricName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceMetricFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMetric_basic(metricName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", metricName),
					resource.TestCheckResourceAttr(rName, "metric_type", "LOGGING"),
					resource.TestCheckResourceAttr(rName, "data_type", "STATISTICS"),
					resource.TestCheckResourceAttr(rName, "cache_ttl", "10"),
					resource.TestCheckResourceAttr(rName, "metric_dimension", "1"),
					resource.TestCheckResourceAttrSet(rName, "is_built_in"),
				),
			},
			{
				Config: testAccResourceMetric_update(metricName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", metricName)),
					resource.TestCheckResourceAttr(rName, "cache_ttl", "20"),
					resource.TestCheckResourceAttr(rName, "metric_dimension", "2"),
					resource.TestCheckResourceAttr(rName, "report_period", "1"),
					resource.TestCheckResourceAttr(rName, "max_query_range", "6"),
					resource.TestCheckResourceAttr(rName, "effective_column", "test_update"),
					resource.TestCheckResourceAttr(rName, "compound_expression", "($1+$3)/$2"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"version"},
				ImportStateIdFunc:       testAccMetricImportStateFunc(rName),
			},
		},
	})
}

func testAccResourceMetric_basic(metricName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_metric" "test" {
  workspace_id     = "%[1]s"
  name             = "%[2]s"
  metric_type      = "LOGGING"
  data_type        = "STATISTICS"
  metric_dimension = 1
  cache_ttl        = 10
  report_period    = 0
  is_built_in      = false
  max_query_range  = 5
  version          = "0.1"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, metricName)
}

func testAccResourceMetric_update(metricName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_metric" "test" {
  workspace_id        = "%[1]s"
  name                = "%[2]s_update"
  metric_type         = "LOGGING"
  data_type           = "STATISTICS"
  metric_dimension    = 2
  cache_ttl           = 20
  report_period       = 1
  is_built_in         = false
  max_query_range     = 6
  effective_column    = "test_update"
  compound_expression = "($1+$3)/$2"
  version             = "0.1"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, metricName)
}

func testAccMetricImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, metricId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		metricId = rs.Primary.ID

		if workspaceId == "" || metricId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<metric_id>', but got '%s/%s'",
				workspaceId, metricId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, metricId), nil
	}
}
