package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDscBatchAddDataMask_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDscBatchAddDataMask_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"huaweicloud_dsc_batch_add_data_mask.test", "masked_data"),
				),
			},
		},
	})
}

const testAccResourceDscBatchAddDataMask_basic = `
resource "huaweicloud_dsc_batch_add_data_mask" "test" {
  mask_strategies {
    name      = "col"
    algorithm = "SHA256"
  }

  data = [
    {
      col = "test1111"
    },
    {
      col = "test22222"
    }
  ]
}
`
