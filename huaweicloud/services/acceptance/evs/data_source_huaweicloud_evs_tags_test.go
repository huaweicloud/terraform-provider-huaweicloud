package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("tags_datasource_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEvsTags_basic(name string) string {
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

data "huaweicloud_evs_tags" "test" {
  depends_on = [huaweicloud_evs_volume.test]
}

output "tags_datasource_useful" {
  value = strcontains(data.huaweicloud_evs_tags.test.tags["foo"], "bar") && strcontains(data.huaweicloud_evs_tags.test.
  tags["key"], "value")
}
`, name)
}
