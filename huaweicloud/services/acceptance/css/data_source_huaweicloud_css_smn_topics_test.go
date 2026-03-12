package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssSmnTopics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_smn_topics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// The user account ID
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssSmnTopics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "topics_name.#"),
				),
			},
		},
	})
}

func testAccDataSourceCssSmnTopics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_smn_topics" "test" {
  domain_id = "%s"
}
`, acceptance.HW_DOMAIN_ID)
}
