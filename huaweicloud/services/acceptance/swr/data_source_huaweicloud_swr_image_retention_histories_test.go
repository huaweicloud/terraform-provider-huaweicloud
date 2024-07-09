package swr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrImageRetentionHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_image_retention_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrRepository(t)
			acceptance.TestAccPreCheckSwrOrigination(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrImageRetentionHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.organization"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.retention_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.rule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.tag"),
					resource.TestMatchResourceAttr(dataSource, "records.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceSwrImageRetentionHistories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_swr_image_retention_histories" "test" {
  organization = "%[1]s"
  repository   = "%[2]s"
}
`, acceptance.HW_SWR_ORGANIZATION, acceptance.HW_SWR_REPOSITORY)
}
