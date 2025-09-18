package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterprisePrivateNetworkAccessControls_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_private_network_access_controls.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterprisePrivateNetworkAccessControls_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.vpcep_endpoint_id"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.endpoint_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.vpc_name"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.vpc_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.subnet_name"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.subnet_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "internal_endpoints.0.created_at"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterprisePrivateNetworkAccessControls_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_private_network_access_controls" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
}
`, testAccSwrEnterpriseInstance_update(name))
}
