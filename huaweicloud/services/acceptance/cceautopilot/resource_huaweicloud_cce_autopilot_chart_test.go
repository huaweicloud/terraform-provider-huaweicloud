package cceautopilot

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

func getAutopilotChartFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getAutopilotChartHttpUrl = "autopilot/v2/charts/{chart_id}"
		getAutopilotChartProduct = "cce"
	)
	getAutopilotChartClient, err := cfg.NewServiceClient(getAutopilotChartProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE client: %s", err)
	}

	getAutopilotChartHttpPath := getAutopilotChartClient.Endpoint + getAutopilotChartHttpUrl
	getAutopilotChartHttpPath = strings.ReplaceAll(getAutopilotChartHttpPath, "{chart_id}", state.Primary.ID)

	getAutopilotChartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAutopilotChartResp, err := getAutopilotChartClient.Request("GET", getAutopilotChartHttpPath, &getAutopilotChartOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE autopilot chart: %s", err)
	}

	return utils.FlattenResponse(getAutopilotChartResp)
}

func TestAccAutopilotChart_basic(t *testing.T) {
	var (
		chart        interface{}
		resourceName = "huaweicloud_cce_autopilot_chart.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&chart,
			getAutopilotChartFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckCceChartPath(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAutopilotChart_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "value"),
					resource.TestCheckResourceAttrSet(resourceName, "instruction"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttrSet(resourceName, "public"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_url"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccAutopilotChart_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "value"),
					resource.TestCheckResourceAttrSet(resourceName, "instruction"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttrSet(resourceName, "public"),
					resource.TestCheckResourceAttrSet(resourceName, "chart_url"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"content", "parameters",
				},
			},
		},
	})
}

func testAccAutopilotChart_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_autopilot_chart" "test" {
  content    = "%s"
  parameters = "{\"override\":true,\"skip_lint\":true,\"source\":\"package\"}"
}
`, acceptance.HW_CCE_CHART_PATH)
}

func testAccAutopilotChart_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_autopilot_chart" "test" {
  content    = "%s"
  parameters = "{\"override\":false,\"skip_lint\":false,\"source\":\"package\"}"
}
`, acceptance.HW_CCE_CHART_PATH)
}
