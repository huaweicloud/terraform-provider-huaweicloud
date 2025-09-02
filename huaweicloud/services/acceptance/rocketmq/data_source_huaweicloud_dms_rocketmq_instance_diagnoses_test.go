package rocketmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDiagnoses_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dms_rocketmq_instance_diagnoses.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		rName          = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDiagnoses_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "reports.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "reports.0.abnormal_item_sum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reports.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reports.0.faulted_node_sum"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reports.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reports.0.report_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reports.0.status"),
				),
			},
		},
	})
}

func testAccDiagnoses_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rocketmq_instance_diagnoses" "test" {
  instance_id = "%s"

  depends_on = [
    huaweicloud_dms_rocketmq_instance_diagnosis.test
  ]
}
`, testAccDiagnosis_basic(name), acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID)
}
