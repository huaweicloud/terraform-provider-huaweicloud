package aom

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAomPromInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_prom_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAomPromInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.prom_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.remote_write_url"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.remote_read_url"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.prom_http_api_endpoint"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.remote_write_url"),
					resource.TestMatchResourceAttr(dataSource, "instances.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "instances.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					resource.TestCheckOutput("prom_type_validation", "true"),
					resource.TestCheckOutput("prom_id_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceAomPromInstances_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "test" {
  prom_name             = "%[1]s"
  prom_type             = "VPC"
  enterprise_project_id = "%[2]s"
  prom_version          = "1.5"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceDataSourceAomPromInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_aom_prom_instances" "test" {
  enterprise_project_id = "all_granted_eps"
  prom_type             = huaweicloud_aom_prom_instance.test.prom_type
  cce_cluster_enable    = "false"
  prom_status           = "NORMAL"
}

data "huaweicloud_aom_prom_instances" "id_filter" {
  prom_id = huaweicloud_aom_prom_instance.test.id
}

locals {
  test_results           = data.huaweicloud_aom_prom_instances.test
  test_id_filter_results = data.huaweicloud_aom_prom_instances.id_filter
}

output "prom_type_validation" {
  value = alltrue([for v in local.test_results.instances[*].prom_type : v == huaweicloud_aom_prom_instance.test.prom_type])
}

output "prom_id_validation" {
  value = alltrue([for v in local.test_id_filter_results.instances[*].id : v == huaweicloud_aom_prom_instance.test.id])
}
`, testDataSourceDataSourceAomPromInstances_base(name))
}
