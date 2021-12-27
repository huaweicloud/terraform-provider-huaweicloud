package deprecated

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpnIPSecPolicyV2_basic(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"huaweicloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "tenant_id", &policy.TenantID),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "transform_protocol", &policy.TransformProtocol),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "encapsulation_mode", &policy.EncapsulationMode),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "encryption_algorithm", &policy.EncryptionAlgorithm),
				),
			},
			{
				Config: testAccIPSecPolicyV2_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"huaweicloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVpnIPSecPolicyV2_withLifetime(t *testing.T) {
	var policy ipsecpolicies.Policy
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"huaweicloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("huaweicloud_vpnaas_ipsec_policy_v2.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIPSecPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSecPolicyV2Exists(
						"huaweicloud_vpnaas_ipsec_policy_v2.policy_1", &policy),
				),
			},
		},
	})
}

func testAccCheckIPSecPolicyV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpnaas_ipsec_policy_v2" {
			continue
		}
		_, err = ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IPSec policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIPSecPolicyV2Exists(n string, policy *ipsecpolicies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := ipsecpolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIPSecPolicyV2_basic = `
resource "huaweicloud_vpnaas_ipsec_policy_v2" "policy_1" {
}
`

const testAccIPSecPolicyV2_Update = `
resource "huaweicloud_vpnaas_ipsec_policy_v2" "policy_1" {
	name = "updatedname"
}
`

const testAccIPSecPolicyV2_withLifetime = `
resource "huaweicloud_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "md5"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIPSecPolicyV2_withLifetimeUpdate = `
resource "huaweicloud_vpnaas_ipsec_policy_v2" "policy_1" {
	auth_algorithm = "md5"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
