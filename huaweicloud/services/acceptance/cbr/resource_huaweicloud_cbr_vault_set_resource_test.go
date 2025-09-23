package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceVaultSetResource_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceVaultSetResource_basic(name),
			},
		},
	})
}

func testResourceVaultSetResource_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[1]s"
  size              = 20
  description       = "test description"
  volume_type       = "GPSSD"
  device_type       = "SCSI"
  multiattach       = false
}

resource "huaweicloud_cbr_vault" "test" {
  name            = "%[1]s"
  type            = "disk"
  protection_type = "backup"
  size            = 50

  resources {
    includes = [huaweicloud_evs_volume.test.id]
  }
}
`, name)
}

func testResourceVaultSetResource_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault_set_resource" "test" {
  vault_id     = huaweicloud_cbr_vault.test.id
  resource_ids = [huaweicloud_evs_volume.test.id]
  action       = "suspend"
}
`, testResourceVaultSetResource_base(name))
}
