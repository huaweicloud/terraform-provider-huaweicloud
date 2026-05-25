package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDriverDelete_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_driver_delete.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDriverDelete_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

const testAccResourceDriverDelete_basic = `
resource "huaweicloud_drs_driver_delete" "test" {
  driver_type  = "db2"
  driver_names = ["xxx.jar"]
}
`
