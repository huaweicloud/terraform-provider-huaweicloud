package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEClusterPodIdentityAssociations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_cluster_pod_identity_associations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCCEClusterPodIdentityAssociations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "associations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "associations.0.uid"),
					resource.TestCheckResourceAttr(dataSource, "associations.0.cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(dataSource, "associations.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "associations.0.service_account"),
					resource.TestCheckResourceAttrSet(dataSource, "associations.0.agency_name"),
					resource.TestCheckResourceAttrSet(dataSource, "associations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "associations.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceCCEClusterPodIdentityAssociations_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_cluster_pod_identity_associations" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_CCE_CLUSTER_ID)
}
