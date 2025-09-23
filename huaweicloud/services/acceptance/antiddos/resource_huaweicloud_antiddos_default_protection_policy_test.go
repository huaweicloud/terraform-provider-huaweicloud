package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/antiddos"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDefaultProtectionPolicyFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("anti-ddos", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Anti-DDoS v1 client: %s", err)
	}

	// Default protection strategies always have value.
	policyResp, err := antiddos.ReadDefaultProtectionPolicy(client)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Anti-DDoS default protection policy: %s", err)
	}

	trafficPosID := utils.PathSearch("traffic_pos_id", policyResp, nil)
	if trafficPosID == nil {
		return nil, fmt.Errorf("error retrieving Anti-DDoS default protection policy: traffic_pos_id is not found in" +
			" read API response")
	}

	// The actual operation of deleting the resource is to set `traffic_pos_id` to `99`.
	if trafficPosID.(float64) == 99 {
		return nil, golangsdk.ErrDefault404{}
	}
	return policyResp, nil
}

func TestAccDefaultProtectionPolicy_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_antiddos_default_protection_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDefaultProtectionPolicyFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDefaultProtectionPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "traffic_threshold", "100"),
				),
			},
			{
				Config: testDefaultProtectionPolicy_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "traffic_threshold", "200"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testDefaultProtectionPolicy_basic = `
resource "huaweicloud_antiddos_default_protection_policy" "test" {
  traffic_threshold = 100
}
`

const testDefaultProtectionPolicy_update = `
resource "huaweicloud_antiddos_default_protection_policy" "test" {
  traffic_threshold = 200
}
`
