package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceBatchUpdateVault_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBatchUpdateVault,
			},
		},
	})
}

const testBatchUpdateVault = `
resource "huaweicloud_cbr_batch_update_vault" "test" {
  smn_notify = true
  threshold  = 80
}
`
