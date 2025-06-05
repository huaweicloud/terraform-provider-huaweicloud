package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsVolumeTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_volume_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsVolumeTags_base(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "tags.0.key", "foo"),
					resource.TestCheckResourceAttr(dataSource, "tags.0.value", "bar"),
					resource.TestCheckResourceAttr(dataSource, "tags.1.key", "key"),
					resource.TestCheckResourceAttr(dataSource, "tags.1.value", "value"),
				),
			},
		},
	})
}

func testDataSourceEvsVolumeTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%s"
  size              = 100
  volume_type       = "GPSSD"
  device_type       = "VBD"
  multiattach       = false

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_evs_volume_tags" "test" {
  volume_id = huaweicloud_evs_volume.test.id
}
`, name)
}
