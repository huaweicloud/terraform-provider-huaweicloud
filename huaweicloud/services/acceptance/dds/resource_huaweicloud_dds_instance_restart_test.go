package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSV3InstanceRestart_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_restart(rName),
			},
		},
	})
}

func testAccDDSInstanceV3Config_restart(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_restart" "test" {
  depends_on = [huaweicloud_dds_instance.instance]
  
  instance_id = huaweicloud_dds_instance.instance.id
}`, testAccDDSInstanceV3Config_basic(rName, 8800))
}
