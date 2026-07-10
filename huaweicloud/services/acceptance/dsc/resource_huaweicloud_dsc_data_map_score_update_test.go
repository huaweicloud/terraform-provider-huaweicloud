package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDscDataMapScoreUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDscDataMapScoreUpdate_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("huaweicloud_dsc_data_map_score_update.test", "status"),
				),
			},
		},
	})
}

const testAccResourceDscDataMapScoreUpdate_basic = `
resource "huaweicloud_dsc_data_map_score_update" "test" {}
`
