package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGroupTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_group_tags.test"
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGroupTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourceGroupTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_configuration" "test"{
  scaling_configuration_name = "%[2]s"
  instance_config {
    image              = data.huaweicloud_images_image.test.id
    flavor             = data.huaweicloud_compute_flavors.test.ids[0]
    key_name           = huaweicloud_kps_keypair.acc_key.id
    security_group_ids = [huaweicloud_networking_secgroup.test.id]

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}

resource "huaweicloud_as_group" "test"{
  scaling_group_name       = "%[2]s"
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  vpc_id                   = huaweicloud_vpc.test.id
  delete_publicip          = true
  delete_volume            = true

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_as_group_tags" "test" {
  resource_type = "scaling_group_tag"

  depends_on = [huaweicloud_as_group.test]
}`, testAccASConfiguration_base(name), name)
}
