package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPipelineGroupSwap_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPipelineGroupSwap_base(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_codearts_pipeline_group.test2", "ordinal", "1"),
				),
			},
			{
				Config: testPipelineGroupSwap_basic(name),
			},
			{
				Config: testPipelineGroupSwap_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_codearts_pipeline_group.test2", "ordinal", "0"),
				),
			},
		},
	})
}

func testPipelineGroupSwap_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_group" "test1" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-1"
}

resource "huaweicloud_codearts_pipeline_group" "test2" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-2"
}
`, testProject_basic(name), name)
}

func testPipelineGroupSwap_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_group_swap" "test" {
  project_id = huaweicloud_codearts_project.test.id
  group_id1  = huaweicloud_codearts_pipeline_group.test1.id
  group_id2  = huaweicloud_codearts_pipeline_group.test2.id
}
`, testPipelineGroupSwap_base(name))
}
