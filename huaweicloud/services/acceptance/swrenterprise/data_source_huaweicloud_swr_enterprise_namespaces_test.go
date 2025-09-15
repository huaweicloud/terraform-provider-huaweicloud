package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrNamespaces_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrNamespaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.#"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.metadata.0.public"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.repo_count"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "namespaces.0.updated_at"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("public_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrNamespaces_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_namespaces" "test" {
  depends_on = [huaweicloud_swr_enterprise_namespace.test]

  instance_id = huaweicloud_swr_enterprise_instance.test.id
}

data "huaweicloud_swr_enterprise_namespaces" "filter_by_name" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = huaweicloud_swr_enterprise_namespace.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_namespaces.filter_by_name.namespaces) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_namespaces.filter_by_name.namespaces[*].name : v == huaweicloud_swr_enterprise_namespace.test.name]
  )
}

data "huaweicloud_swr_enterprise_namespaces" "filter_by_public" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  public      = huaweicloud_swr_enterprise_namespace.test.metadata[0].public
}

output "public_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_namespaces.filter_by_public.namespaces) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_namespaces.filter_by_public.namespaces[*].metadata[0].public : 
	  v == huaweicloud_swr_enterprise_namespace.test.metadata[0].public]
  )
}
`, testAccSwrEnterpriseNamespace_basic(name))
}
