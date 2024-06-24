package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDMInstanceRestart_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ddm_instance_restart.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDdmInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDMInstanceRestart_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDDMInstanceRestart_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_instance_restart" "test" {
  depends_on = [huaweicloud_ddm_instance.test]

  instance_id = huaweicloud_ddm_instance.test.id
  type        = "soft"
}`, testDdmInstance_basic(rName))
}
