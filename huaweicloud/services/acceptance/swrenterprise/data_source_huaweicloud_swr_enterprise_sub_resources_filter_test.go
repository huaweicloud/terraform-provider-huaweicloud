package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseSubResourcesFilter_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_sub_resources_filter.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseSubResourcesFilter_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.%"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseSubResourcesFilter_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_sub_resources_filter" "test" {
  depends_on = [huaweicloud_swr_enterprise_namespace.test]

  resource_type     = "instances"
  resource_id       = huaweicloud_swr_enterprise_instance.test.id
  sub_resource_type = "namespaces"
}

data "huaweicloud_swr_enterprise_sub_resources_filter" "filter_by_tags" {
  depends_on = [huaweicloud_swr_enterprise_namespace.test]

  resource_type     = "instances"
  resource_id       = huaweicloud_swr_enterprise_instance.test.id
  sub_resource_type = "namespaces"

  tags {
    key    = "value1"
    values = ["bar1"]
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_sub_resources_filter.filter_by_tags.resources) > 0
}
`, testAccSwrEnterpriseNamespace_update(name))
}
