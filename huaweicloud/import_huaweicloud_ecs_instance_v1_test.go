package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccEcsV1Instance_importBasic(t *testing.T) {
	resourceName := "huaweicloud_ecs_instance_v1.instance_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1Instance_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}
