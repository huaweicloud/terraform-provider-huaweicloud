package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterAction_basic(t *testing.T) {
	randUUID, _ := uuid.NewRandom()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterAction_nonExistentCluster(randUUID.String()),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`error operating cluster \(%s\) with action \(restart\)`, randUUID.String())),
			},
			{
				Config: testAccClusterAction_basic("stop"),
			},
			{
				Config: testAccClusterAction_basic("start"),
			},
			{
				Config: testAccClusterAction_basic("restart"),
			},
		},
	})
}

func testAccClusterAction_nonExistentCluster(randUUID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_action" "test" {
  cluster_id = "%[1]s"
  action     = "restart"
}
`, randUUID)
}

func testAccClusterAction_basic(action string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_action" "test" {
  cluster_id = "%[1]s"
  action     = "%[2]s"

  enable_force_new  = "true"
}
`, acceptance.HW_DWS_CLUSTER_ID, action)
}
