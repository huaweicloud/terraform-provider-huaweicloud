package dsc

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMultiOrganizationTree_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_multi_organization_tree.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMultiOrganizationTree_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ou_list.#"),
				),
			},
		},
	})
}

func testDataSourceMultiOrganizationTree_basic() string {
	randomUUID, _ := uuid.NewRandom()
	return fmt.Sprintf(`
data "huaweicloud_dsc_multi_organization_tree" "test" {
  entity_id = "%s"
}
`, randomUUID.String())
}
