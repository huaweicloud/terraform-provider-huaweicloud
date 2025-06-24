package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImageBatchScan_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImageBatchScan_basic,
			},
		},
	})
}

const testAccImageBatchScan_basic string = `
resource "huaweicloud_hss_image_batch_scan" "test_basci" {
  image_type = "private_image"
  repo_type  = "SWR"

  image_info_list {
    namespace     = "test-namespace"
    image_name    = "test-image"
    image_version = "latest"
  }
}

resource "huaweicloud_hss_image_batch_scan" "test_operate_all" {
  image_type  = "private_image"
  repo_type   = "SWR"
  operate_all = true
}
`
