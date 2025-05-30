package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV3VolumeTypes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_evsv3_volume_types.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceV3VolumeTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "volume_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_types.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_types.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_types.0.extra_specs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "volume_types.0.extra_specs.0.availability_zones"),
				),
			},
		},
	})
}

const testDataSourceV3VolumeTypes_basic = `
data "huaweicloud_evsv3_volume_types" "test" {}
`
