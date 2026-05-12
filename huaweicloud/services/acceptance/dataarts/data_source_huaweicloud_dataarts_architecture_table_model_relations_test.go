package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureTableModelRelations_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_dataarts_architecture_table_model_relations.test"
		all    = acceptance.InitDataSourceCheck(dcName)

		byBizTypeNotFound   = "data.huaweicloud_dataarts_architecture_table_model_relations.filter_by_biz_not_found"
		dcbyBizTypeNotFound = acceptance.InitDataSourceCheck(byBizTypeNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsModelID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchitectureTableModelRelations_basic(),
				Check: resource.ComposeTestCheckFunc(
					all.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "relations.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "relations.#"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.source_table_id"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.source_table_name"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.target_table_id"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.target_table_name"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.created_at"),
					resource.TestCheckResourceAttrSet(dcName, "relations.0.updated_at"),
					dcbyBizTypeNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_biz_type_filter_useful_not_found", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureTableModelRelations_basic() string {
	return fmt.Sprintf(`
# Query all relations and without any filter
data "huaweicloud_dataarts_architecture_table_model_relations" "test" {
  workspace_id = "%[1]s"
  model_id     = "%[2]s"
}

# Test filter by biz type not found
data "huaweicloud_dataarts_architecture_table_model_relations" "filter_by_biz_not_found" {
  workspace_id = "%[1]s"
  model_id     = "%[2]s"
  biz_type     = "FACT_LOGIC_TABLE"
}

output "is_biz_type_filter_useful_not_found" {
  value = length(data.huaweicloud_dataarts_architecture_table_model_relations.filter_by_biz_not_found.relations) == 0
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_MODEL_ID)
}
