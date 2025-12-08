package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_smn_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourceSmnTags_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_smn_tags" "test" {
  depends_on = [huaweicloud_smn_topic.test]

  resource_type = "smn_topic"
}
`, name)
}
