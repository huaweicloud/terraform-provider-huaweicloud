package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to limitations in the testing environment, this data source has no return value
// and can only be checked for successful execution.
func TestAccDataSourceKubernetesContainerDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_kubernetes_container_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKubernetesContainerDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

// The `container_id` used is dummy data for testing.
func testDataSourceKubernetesContainerDetail_basic() string {
	return `
data "huaweicloud_hss_kubernetes_container_detail" "test" {
  container_id          = "19031ec457570cce64f789a57ce6551509"
  enterprise_project_id = "0"
}
`
}
