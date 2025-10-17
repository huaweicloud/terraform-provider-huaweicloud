package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKubernetesPodDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_pod_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubernetesPodDetail_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "containers.#"),
				),
			},
		},
	})
}

const testAccDataSourceKubernetesPodDetail_basic = `
data "huaweicloud_hss_kubernetes_pod_detail" "test" {
  pod_name              = "non-exist"
  enterprise_project_id = "0"
}
`
