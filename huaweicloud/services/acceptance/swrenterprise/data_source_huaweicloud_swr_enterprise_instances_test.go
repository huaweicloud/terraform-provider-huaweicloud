package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instances.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.spec"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.obs_bucket_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.access_address"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.expires_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.user_def_obs"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_cidr"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_instances" "test" {
  depends_on = [huaweicloud_swr_enterprise_instance.test]
}

data "huaweicloud_swr_enterprise_instances" "filter_by_status" {
  status = huaweicloud_swr_enterprise_instance.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instances.filter_by_status.instances) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_instances.filter_by_status.instances[*].status : 
	  v == huaweicloud_swr_enterprise_instance.test.status]
  )
}

data "huaweicloud_swr_enterprise_instances" "filter_by_enterprise_project_id" {
  enterprise_project_id = huaweicloud_swr_enterprise_instance.test.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_instances.filter_by_enterprise_project_id.instances) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_instances.filter_by_enterprise_project_id.instances[*].enterprise_project_id : 
	  v == huaweicloud_swr_enterprise_instance.test.enterprise_project_id]
  )
}
`, testAccSwrEnterpriseInstance_update(name))
}
