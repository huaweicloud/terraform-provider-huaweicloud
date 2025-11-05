package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceRegistries_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_registries.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceRegistries_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "registries.#"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.url"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.insecure"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.credential.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.credential.0.access_secret"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.credential.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "registries.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "total"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceRegistries_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_instance_registries" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_instance_registries" "filter_by_name" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name        = data.huaweicloud_swr_enterprise_instance_registries.test.registries[0].name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instance_registries.filter_by_name.registries) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_instance_registries.filter_by_name.registries[*].name :
	  strcontains(v, data.huaweicloud_swr_enterprise_instance_registries.test.registries[0].name)]
  )
}

data "huaweicloud_swr_enterprise_instance_registries" "filter_by_type" {
  instance_id  = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  type = data.huaweicloud_swr_enterprise_instance_registries.test.registries[0].type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instance_registries.filter_by_type.registries) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_instance_registries.filter_by_type.registries[*].type :
	  v == data.huaweicloud_swr_enterprise_instance_registries.test.registries[0].type]
  )
}
`
}
