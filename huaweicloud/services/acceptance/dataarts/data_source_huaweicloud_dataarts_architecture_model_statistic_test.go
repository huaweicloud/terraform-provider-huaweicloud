package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before test, make sure at least one of logic model has been created.
func TestAccDatasourceModelStatistic_basic(t *testing.T) {
	rName := "data.huaweicloud_dataarts_architecture_model_statistic.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceModelStatistic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "logics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDatasourceModelStatistic_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_model_statistic" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
