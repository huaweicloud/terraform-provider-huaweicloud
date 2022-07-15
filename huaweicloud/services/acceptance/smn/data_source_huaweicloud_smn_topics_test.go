package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTopics_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_smn_topics.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTopicsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.id",
						"huaweicloud_smn_topic.topic_1", "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.topic_urn",
						"huaweicloud_smn_topic.topic_1", "topic_urn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.display_name",
						"huaweicloud_smn_topic.topic_1", "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.enterprise_project_id",
						"huaweicloud_smn_topic.topic_1", "enterprise_project_id"),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "topics.0.tags.foo", "huaweicloud_smn_topic.topic_1", "tags.foo"),
					resource.TestCheckResourceAttrPair(
						dataSourceName, "topics.0.tags.key", "huaweicloud_smn_topic.topic_1", "tags.key"),
				),
			},
		},
	})
}

func testAccDataTopicsConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_smn_topics" "test" {
  name = "%s"

  depends_on = [
    huaweicloud_smn_topic.topic_1
  ]
}
`, testAccSMNV2TopicConfig_basic(rName), rName)
}
