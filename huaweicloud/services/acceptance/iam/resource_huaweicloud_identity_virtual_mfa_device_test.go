package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityVirtualMFADeviceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getMFAClient, err := conf.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM Client: %s", err)
	}
	getMFAHttpUrl := "v3.0/OS-MFA/users/{user_id}/virtual-mfa-device"
	getMFAPath := getMFAClient.Endpoint + getMFAHttpUrl
	getMFAPath = strings.ReplaceAll(getMFAPath, "{user_id}", state.Primary.Attributes["user_id"])
	getMFAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getMFAResp, err := getMFAClient.Request("GET", getMFAPath, &getMFAOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM virtual MFA device: %s", err)
	}
	return utils.FlattenResponse(getMFAResp)
}

func TestAccIdentityVirtualMFADevice_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_virtual_mfa_device.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getIdentityVirtualMFADeviceResourceFunc)

		name = acceptance.RandomAccResourceName()
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
				Config: testAccIdentityVirtualMFADevice_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_USER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "base32_string_seed"),
					resource.TestCheckResourceAttrSet(resourceName, "qr_code_png"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityVirtualMFADevice(resourceName),
				ImportStateVerifyIgnore: []string{
					"base32_string_seed", "qr_code_png",
				},
			},
		},
	})
}

func testAccIdentityVirtualMFADevice_basic(name string) string {
	return fmt.Sprintf(`
// Only a user can create a virtual MFA device for himself.
// Even a primary user cannot create an MFA for a sub-user.

resource "huaweicloud_identity_virtual_mfa_device" "test" {
  name    = "%[1]s"
  user_id = "%[2]s"
}
`, name, acceptance.HW_USER_ID)
}

func testIdentityVirtualMFADevice(name string) resource.ImportStateIdFunc {
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
