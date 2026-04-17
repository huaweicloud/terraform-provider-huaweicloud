package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccVpcepEndpointUpgrade_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testVpcepEndpointUpgrade_basic(rName),
			},
		},
	})
}

func testVpcepEndpointUpgrade_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpcep_public_services" "test" {
  service_name = "iam"
}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id = data.huaweicloud_vpcep_public_services.test.services[0].id
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_vpcep_endpoint_upgrade" "test" {
  vpc_endpoint_id = huaweicloud_vpcep_endpoint.test.id
  action          = "start"
}
`, common.TestVpc(name), name)
}
