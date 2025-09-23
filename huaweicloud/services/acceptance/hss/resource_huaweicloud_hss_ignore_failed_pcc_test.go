package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIgnoreFailedPCC_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIgnoreFailedPCC_basic,
			},
		},
	})
}

const testAccIgnoreFailedPCC_basic string = `
resource "huaweicloud_hss_ignore_failed_pcc" "test" {
  action                = "ignore"
  operate_all           = true
  enterprise_project_id = "0"
}
`
