package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCertificateRevoke_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateRevoke_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAccCertificateRevoke_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_cluster_certificate_revoke" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  user_id    = "%[2]s"
}
`, testAccCluster_basic(name), acceptance.HW_USER_ID)
}
