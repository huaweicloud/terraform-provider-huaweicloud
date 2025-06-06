package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSnapshotMetadata_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_evs_snapshot_metadata.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSnapshotMetadata_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.key", "value"),
				),
			},
		},
	})
}

func testDataSourceSnapshotMetadata_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}

resource "huaweicloud_evsv3_snapshot" "test" {
  volume_id = huaweicloud_evs_volume.test.id
  name      = "%[1]s"
}

resource "huaweicloud_evs_snapshot_metadata" "test" {
  snapshot_id = huaweicloud_evsv3_snapshot.test.id

  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testDataSourceSnapshotMetadata_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evs_snapshot_metadata" "test" {
  depends_on = [huaweicloud_evs_snapshot_metadata.test]

  snapshot_id = huaweicloud_evsv3_snapshot.test.id
}
`, testDataSourceSnapshotMetadata_base(name))
}
