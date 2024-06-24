package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcSearchCriteria_basic(t *testing.T) {
	dataSource := "data.huaweicloud_lts_search_criteria.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSearchCriteria_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "search_criteria.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "search_criteria.0.log_stream_id"),
					resource.TestCheckResourceAttrSet(dataSource, "search_criteria.0.log_stream_name"),
					resource.TestMatchResourceAttr(dataSource, "search_criteria.0.criteria.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "search_criteria.0.criteria.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "search_criteria.0.criteria.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "search_criteria.0.criteria.0.criteria"),
					resource.TestCheckResourceAttrSet(dataSource, "search_criteria.0.criteria.0.type"),
				),
			},
			{
				Config:      testDataSourceSearchCriteria_expectError(),
				ExpectError: regexp.MustCompile("The log group does not existed"),
			},
		},
	})
}

func testDataSourceSearchCriteria_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lts_search_criteria" "test" {
  depends_on = [
    huaweicloud_lts_search_criteria.test
  ]

  log_group_id = huaweicloud_lts_group.test.id
}
`, testSearchCriteria_basic(name))
}

func testDataSourceSearchCriteria_expectError() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_lts_search_criteria" "test" {
  log_group_id = "%s"
}
`, randUUID)
}
