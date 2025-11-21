package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityV5VirtualMFADeviceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getMFAClient, err := conf.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM Client: %s", err)
	}
	getMFAHttpUrl := "v5/mfa-devices?user_id=" + state.Primary.Attributes["user_id"]
	getMFAPath := getMFAClient.Endpoint + getMFAHttpUrl
	getMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getMFAResp, err := getMFAClient.Request("GET", getMFAPath, &getMFAOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM virtual MFA device: %s", err)
	}
	respBody, err := utils.FlattenResponse(getMFAResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten IAM virtual MFA device: %s", err)
	}
	check := utils.PathSearch("mfa_devices[0]", respBody, nil)
	if check == nil {
		return nil, fmt.Errorf("error retrieving IAM virtual MFA device: %s", err)
	}
	return check, nil
}

func TestAccIdentityV5VirtualMFADevice_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_identityv5_virtual_mfa_device.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getIdentityV5VirtualMFADeviceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV5VirtualMFADevice_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "base32_string_seed"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityV5VirtualMFADevice(resourceName),
				ImportStateVerifyIgnore: []string{
					"name", "base32_string_seed",
				},
			},
		},
	})
}

func testAccIdentityV5VirtualMFADevice_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_virtual_mfa_device" "test" {
  name    = "%s"
  user_id = "%s"
}
`, name, acceptance.HW_USER_ID)
}

func testIdentityV5VirtualMFADevice(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		userID := rs.Primary.Attributes["user_id"]
		if userID == "" {
			return "", fmt.Errorf("attribute (user_id) of Resource (%s) not found", name)
		}

		return userID, nil
	}
}
