package dis

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDisPartion_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_dis_partitions.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDisPartion_basic(randName, 2),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "partitions.#", "2"),
				),
			},
		},
	})
}

func testAccDataSourceDisPartion_basic(name string, partitionCount int) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dis_partitions" "test" {
  stream_name = "${huaweicloud_dis_stream.test.stream_name}"
}
`, testAccDisStream_basic(name, partitionCount))
}
