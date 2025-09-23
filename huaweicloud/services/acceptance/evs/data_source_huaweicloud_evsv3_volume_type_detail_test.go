package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsv3VolumeTypeDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_evsv3_volume_type_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsv3VolumeTypeDetail_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_type.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_type.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_type.0.extra_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_type.0.extra_specs.0.reskey_availability_zones"),
				),
			},
		},
	})
}

const testDataSourceEvsv3VolumeTypeDetail_basic = `
data "huaweicloud_evs_volume_types" "test" {}

data "huaweicloud_evsv3_volume_type_detail" "test" {
  type_id = data.huaweicloud_evs_volume_types.test.types[0].id
}
`
