package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEClusterPodIdentityAssumeAgency_basic(t *testing.T) {
	datasourceName := "data.huaweicloud_cce_cluster_pod_identity_assume_agency.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceClusterId(t)
			acceptance.TestAccPreCheckCceClusterServiceAccountToken(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClusterPodIdentityAssumeAgency_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "token", acceptance.HW_CCE_CLUSTER_SERVICE_ACCOUNT_TOKEN),
					resource.TestCheckResourceAttrSet(datasourceName, "assumed_agency.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "credentials.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "credentials.0.access_key_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "credentials.0.secret_access_key"),
					resource.TestCheckResourceAttrSet(datasourceName, "credentials.0.security_token"),
					resource.TestCheckResourceAttrSet(datasourceName, "credentials.0.expiration"),
					resource.TestCheckResourceAttrSet(datasourceName, "subject.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "subject.0.namespace"),
					resource.TestCheckResourceAttrSet(datasourceName, "subject.0.service_account"),
				),
			},
		},
	})
}

func testAccCCEClusterPodIdentityAssumeAgency_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cce_cluster_pod_identity_assume_agency" "test" {
  cluster_id = "%[1]s"
  token      = "%[2]s"
}`, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_SERVICE_ACCOUNT_TOKEN)
}
