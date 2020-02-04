package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNatDnat_importBasic(t *testing.T) {
	resourceName := "huaweicloud_nat_dnat_rule_v2.dnat"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatDnatDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatDnat_basic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
