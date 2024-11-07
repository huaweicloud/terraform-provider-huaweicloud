package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSParameterTemplateApply_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplateApply_basic(name),
			},
		},
	})
}

func testParameterTemplateApply_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template" "test" {
  name         = "%[2]s"
  node_type    = "replica"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost = 500
  }
}

resource "huaweicloud_dds_parameter_template_apply" "test" {
  configuration_id = huaweicloud_dds_parameter_template.test.id
  entity_ids       = [huaweicloud_dds_instance.instance.id]
}
`, testAccDDSInstanceReplicaSetBasic(name), name)
}
