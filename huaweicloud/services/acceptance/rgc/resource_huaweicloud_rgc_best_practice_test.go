package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBestPractice_basic(t *testing.T) {
	rName := "huaweicloud_rgc_best_practice.best_practice"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBestPractice_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "total_score"),
					resource.TestCheckResourceAttrSet(rName, "detect_time"),
				),
			},
		},
	})
}

func testBestPractice_basic() string {
	return `resource "huaweicloud_rgc_best_practice" "best_practice" {}`
}
