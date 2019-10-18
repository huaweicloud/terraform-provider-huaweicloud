package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDNSV2RecordSet_importBasic(t *testing.T) {
	zoneName := randomZoneName()
	resourceName := "huaweicloud_dns_recordset_v2.recordset_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSet_basic(zoneName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
