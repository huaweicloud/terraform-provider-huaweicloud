package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapStarrocksReplicationPause_basic(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_taurusdb_htap_starrocks_replication.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getHtapInstanceReplicationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksReplicationPause_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksReplicationPause_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_taurusdb_htap_starrocks_replication_pause" "test" {
  instance_id = huaweicloud_taurusdb_htap_starrocks_replication.test.instance_id
  task_name   = huaweicloud_taurusdb_htap_starrocks_replication.test.task_name

  depends_on = [huaweicloud_taurusdb_htap_starrocks_replication.test]
}
`, testAccTaurusDBHtapStarrocksReplication_base(name))
}
