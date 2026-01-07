package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImageBaselineChangeEWP_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testImageBaselineChangeEWP_basic(),
			},
		},
	})
}

func testImageBaselineChangeEWP_basic() string {
	return `
resource "huaweicloud_hss_image_baseline_change_ewp" "test" {
  extended_weak_password = ["11"]
}
`
}

func TestAccImageBaselineChangeEWP_empty(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testImageBaselineChangeEWP_empty(),
			},
		},
	})
}

func testImageBaselineChangeEWP_empty() string {
	return `
resource "huaweicloud_hss_image_baseline_change_ewp" "test" {
  extended_weak_password = []
}
`
}
