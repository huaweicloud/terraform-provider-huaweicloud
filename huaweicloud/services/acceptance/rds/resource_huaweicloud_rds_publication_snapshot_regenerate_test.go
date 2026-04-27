package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRdsPublicationSnapshotRegenerate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRdsPublicationSnapshotRegenerate_basic(),
			},
		},
	})
}

func testAccResourceRdsPublicationSnapshotRegenerate_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publications" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_publication_snapshot_regenerate" "test" {
  instance_id    = "%[1]s"
  publication_id = data.huaweicloud_rds_publications.test.publications[0].id
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
