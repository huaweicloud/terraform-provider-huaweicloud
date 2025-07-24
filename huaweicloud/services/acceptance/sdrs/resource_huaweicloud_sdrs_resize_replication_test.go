package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccResizeReplication_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResizeReplication_basic(),
			},
		},
	})
}

func testResizeReplication_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}
data "huaweicloud_sdrs_domain" "test" {}

resource "huaweicloud_sdrs_protection_group" "test" {
  name                     = "%[2]s"
  source_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  target_availability_zone = data.huaweicloud_availability_zones.test.names[1]
  domain_id                = data.huaweicloud_sdrs_domain.test.id
  source_vpc_id            = huaweicloud_vpc.test.id
  description              = "test description"
}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  description       = "test volume for sdrs replication pair"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SSD"
  size              = 100

  lifecycle {
    ignore_changes = [size]
  }
}

resource "huaweicloud_sdrs_replication_pair" "test" {
  name                 = "%[2]s"
  group_id             = huaweicloud_sdrs_protection_group.test.id
  volume_id            = huaweicloud_evs_volume.test.id
  description          = "test description"
  delete_target_volume = true
}

resource "huaweicloud_sdrs_resize_replication" "test" {
  replication_id = huaweicloud_sdrs_replication_pair.test.id
  new_size       = 150
}
`, common.TestVpc(name), name)
}
