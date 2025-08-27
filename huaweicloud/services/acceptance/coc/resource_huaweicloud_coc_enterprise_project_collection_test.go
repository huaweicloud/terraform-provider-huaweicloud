package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocEnterpriseProjectCollection_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocEnterpriseProjectCollection_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocEnterpriseProjectCollection_basic() string {
	return `
resource "huaweicloud_coc_enterprise_project_collection" "test" {
  ep_id_list = ["0"]
}
`
}
