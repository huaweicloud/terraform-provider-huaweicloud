package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getWorkloadConfigurationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	return dws.GetWorkloadConfiguration(client, state.Primary.ID)
}

// lintignore:AT001
func TestAccWorkloadConfiguration_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_dws_workload_configuration.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getWorkloadConfigurationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testWorkloadConfiguration_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workload_switch", "off"),
				),
			},
			{
				Config: testWorkloadConfiguration_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workload_switch", "on"),
					resource.TestCheckResourceAttr(rName, "max_concurrency_num", "0"),
				),
			},
			{
				Config: testWorkloadConfiguration_basic_step3(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workload_switch", "on"),
					resource.TestCheckResourceAttr(rName, "max_concurrency_num", "-1"),
				),
			},
			{
				Config: testWorkloadConfiguration_basic_step4(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workload_switch", "on"),
					resource.TestCheckResourceAttr(rName, "max_concurrency_num", "100"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testWorkloadConfiguration_basic(concurrencyNum string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_workload_configuration" "test" {
  cluster_id          = "%s"
  workload_switch     = "on"
  max_concurrency_num = "%s"
}
`, acceptance.HW_DWS_CLUSTER_ID, concurrencyNum)
}

func testWorkloadConfiguration_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_workload_configuration" "test" {
  cluster_id      = "%s"
  workload_switch = "off"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}

func testWorkloadConfiguration_basic_step2() string {
	return testWorkloadConfiguration_basic("0")
}

func testWorkloadConfiguration_basic_step3() string {
	return testWorkloadConfiguration_basic("-1")
}

func testWorkloadConfiguration_basic_step4() string {
	return testWorkloadConfiguration_basic("100")
}
