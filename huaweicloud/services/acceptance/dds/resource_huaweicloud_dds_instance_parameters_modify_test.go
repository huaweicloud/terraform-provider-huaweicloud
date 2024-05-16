package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSV3InstanceModifyParams_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3ModifyParams_basic(rName),
			},
			{
				Config: testAccDDSInstanceV3ModifyParams_update(rName),
			},
		},
	})
}

func testAccDDSInstanceV3ModifyParams_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_parameters_modify" "test" {
  depends_on = [huaweicloud_dds_instance.instance]
  
  instance_id = huaweicloud_dds_instance.instance.id

  parameters {
    name  = "oplogSizePercent"
    value = "0.2"
  }
}`, testAccDDSInstanceReplicaSetBasic(rName))
}

func testAccDDSInstanceV3ModifyParams_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_parameters_modify" "test" {
  depends_on = [huaweicloud_dds_instance.instance]
  
  instance_id = huaweicloud_dds_instance.instance.id

  parameters {
    name  = "connPoolMaxConnsPerHost"
    value = "800"
  }
}`, testAccDDSInstanceReplicaSetBasic(rName))
}
