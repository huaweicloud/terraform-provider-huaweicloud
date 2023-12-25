package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	cce "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getChartFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.HcCceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}

	req := cce.ShowChartRequest{
		ChartId: state.Primary.ID,
	}

	return client.ShowChart(&req)
}

func TestAccChart_basic(t *testing.T) {
	var (
		chart        cce.ShowChartResponse
		resourceName = "huaweicloud_cce_chart.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&chart,
			getChartFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckCceChartPath(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccChart_basic(),
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
				Config: testAccChart_update(),
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

func testAccChart_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_chart" "test" {
  content    = "%s"
  parameters = "{\"override\":true,\"skip_lint\":true,\"source\":\"package\"}"
}
`, acceptance.HW_CCE_CHART_PATH)
}

func testAccChart_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_cce_chart" "test" {
  content    = "%s"
  parameters = "{\"override\":false,\"skip_lint\":false,\"source\":\"package\"}"
}
`, acceptance.HW_CCE_CHART_PATH)
}
