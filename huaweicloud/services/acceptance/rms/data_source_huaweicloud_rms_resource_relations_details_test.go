package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsRelationsDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_relations_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsRelationsDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.relation_type"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.from_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.to_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.from_resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.to_resource_id"),
				),
			},
		},
	})
}

func testDataSourceRmsRelationsDetails_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_relations_details" "test" {
  resource_id = huaweicloud_compute_instance.test.id
  direction   = "in"
}
`, testDataSourceRmsHistories_base())
}
