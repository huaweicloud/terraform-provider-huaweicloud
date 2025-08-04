package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSfsTurboShareTypes_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_sfs_turbo_share_types.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSfsTurboShareTypes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.share_type"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.capacity.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.capacity.0.max"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.capacity.0.min"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.capacity.0.step"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.bandwidth.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.bandwidth.0.max"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.bandwidth.0.min"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.bandwidth.0.step"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.bandwidth.0.density"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.bandwidth.0.base"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.iops.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.iops.0.max"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.iops.0.min"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.single_channel_4k_latency.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.single_channel_4k_latency.0.max"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.attribution.0.single_channel_4k_latency.0.min"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.support_period"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.available_zones.#"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.available_zones.0.available_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.available_zones.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.spec_code"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.storage_media"),
					resource.TestCheckResourceAttrSet(resourceName, "share_types.0.features.#"),
				),
			},
		},
	})
}

func testAccDataSourceSfsTurboShareTypes_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {}

data "huaweicloud_vpc_subnet" "test" {
  vpc_id = data.huaweicloud_vpc.test.id
}

data "huaweicloud_networking_secgroup" "test" {}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%s"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  
  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_sfs_turbo_share_types" "test" {}
`, rName)
}
