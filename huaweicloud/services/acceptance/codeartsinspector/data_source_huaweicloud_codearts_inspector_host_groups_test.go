package codeartsinspector

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsInspectorHostGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_inspector_host_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsInspectorHostGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
				),
			},
		},
	})
}

func testDataSourceCodeartsInspectorHostGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_inspector_host_groups" "test" {
  depends_on = [huaweicloud_codearts_inspector_host_group.test]
}
`, testInspectorHostGroup_basic(name))
}
