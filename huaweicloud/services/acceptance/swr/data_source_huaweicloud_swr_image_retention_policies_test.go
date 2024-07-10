package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageRetentionPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_image_retention_policies.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImageRetentionPolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.rules.0.template"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.rules.0.params"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.rules.0.tag_selectors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "retention_policies.0.rules.0.tag_selectors.0.kind"),
					resource.TestCheckResourceAttrSet(dataSource,
						"retention_policies.0.rules.0.tag_selectors.0.pattern"),
				),
			},
		},
	})
}

func testAccDataSourceImageRetentionPolicies_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_swr_image_retention_policies" "test" {
  depends_on = [huaweicloud_swr_image_retention_policy.test]
  
  organization = huaweicloud_swr_organization.test.name
  repository   = huaweicloud_swr_repository.test.name
}
`, testSwrImageRetentionPolicy_basic(rName))
}
