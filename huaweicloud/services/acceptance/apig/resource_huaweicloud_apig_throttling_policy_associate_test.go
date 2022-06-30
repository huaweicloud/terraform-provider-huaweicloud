package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opt := throttles.ListBindOpts{
		InstanceId: state.Primary.Attributes["instance_id"],
		ThrottleId: state.Primary.Attributes["policy_id"],
	}
	return throttles.ListBind(c, opt)
}

func TestAccThrottlingPolicyAssociate_basic(t *testing.T) {
	var apiDetails []throttles.ApiForThrottle

	// The dedicated instance name only allow letters, digits and underscores (_).
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_apig_throttling_policy_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&apiDetails,
		getAssociateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccThrottlingPolicyAssociate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "policy_id",
						"huaweicloud_apig_throttling_policy.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "publish_ids.#", "1"),
				),
			}, {
				Config: testAccThrottlingPolicyAssociate_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "policy_id",
						"huaweicloud_apig_throttling_policy.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "publish_ids.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccThrottlingPolicyAssociate_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_environment" "test" {
  count = 2

  name        = "%s_${count.index}"
  instance_id = huaweicloud_apig_instance.test.id
}

resource "huaweicloud_apig_api_publishment" "test" {
  count = 2

  instance_id = huaweicloud_apig_instance.test.id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[count.index].id
}

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"
}
`, testAccApigAPI_basic(rName), rName, rName)
}

func testAccThrottlingPolicyAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_throttling_policy_associate" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  policy_id   = huaweicloud_apig_throttling_policy.test.id

  publish_ids = [
    huaweicloud_apig_api_publishment.test[0].publish_id
  ]
}
`, testAccThrottlingPolicyAssociate_base(rName))
}

func testAccThrottlingPolicyAssociate_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_throttling_policy_associate" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  policy_id   = huaweicloud_apig_throttling_policy.test.id

  publish_ids = [
    huaweicloud_apig_api_publishment.test[1].publish_id
  ]
}
`, testAccThrottlingPolicyAssociate_base(rName))
}
