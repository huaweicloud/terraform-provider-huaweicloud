package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSV3InstanceRestart_basic(t *testing.T) {
	var instance instances.InstanceResponse
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance_restart.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDdsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_restart(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
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
