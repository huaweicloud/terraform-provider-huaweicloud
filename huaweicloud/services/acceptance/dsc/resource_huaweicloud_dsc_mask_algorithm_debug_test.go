package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test, please ensure that you have created a DSC instance.
func TestAccMaskAlgorithmDebug_basic(t *testing.T) {
	var (
		rName              = "huaweicloud_dsc_mask_algorithm_debug.test"
		rNameWithParameter = "huaweicloud_dsc_mask_algorithm_debug.with_parameter"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMaskAlgorithmDebug_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "algorithm", "SHA256"),
					resource.TestCheckResourceAttr(rName, "algorithm_type", "MASK_BY_HASH"),
					resource.TestCheckResourceAttrSet(rName, "processed_data"),
					// Test the resource with parameter.
					resource.TestCheckResourceAttr(rNameWithParameter, "data", "1234567890123456789"),
					resource.TestCheckResourceAttr(rNameWithParameter, "algorithm", "PRESNM"),
					resource.TestCheckResourceAttr(rNameWithParameter, "algorithm_type", "MASK_BY_OVERWRITE"),
					resource.TestCheckResourceAttr(rNameWithParameter, "processed_data", "1234***********6789"),
				),
			},
		},
	})
}

const testAccMaskAlgorithmDebug_basic = `
resource "huaweicloud_dsc_mask_algorithm_debug" "test" {
  data           = "test"
  algorithm      = "SHA256"
  algorithm_type = "MASK_BY_HASH"
}

resource "huaweicloud_dsc_mask_algorithm_debug" "with_parameter" {
  data           = "1234567890123456789"
  algorithm      = "PRESNM"
  algorithm_type = "MASK_BY_OVERWRITE"
  parameter      = jsonencode({
    type   = "CHAR"
    first  = 4
    second = 4
    method = "*"
  })
}
`
