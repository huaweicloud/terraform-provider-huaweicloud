package modelarts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceServiceFlavors_basic(t *testing.T) {
	rName := "data.huaweicloud_modelarts_service_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceServiceFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.is_open"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.status"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.billing_spec"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.is_free"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.over_quota"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.extend_params"),

					resource.TestCheckOutput("is_open_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("is_free_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceServiceFlavors_basic() string {
	return `
data "huaweicloud_modelarts_service_flavors" "test" {
}

data "huaweicloud_modelarts_service_flavors" "is_open_filter" {
  is_open = true
}
output "is_open_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_modelarts_service_flavors.is_open_filter.flavors[*].is_open : v == true])
}

data "huaweicloud_modelarts_service_flavors" "status_filter" {
  status = "normal"
}
output "status_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_modelarts_service_flavors.status_filter.flavors[*].status : v == "normal"])
}

data "huaweicloud_modelarts_service_flavors" "is_free_filter" {
  is_free = true
}
output "is_free_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_modelarts_service_flavors.is_free_filter.flavors[*].is_free : v == true])
}
`
}
