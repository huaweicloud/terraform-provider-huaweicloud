package identitycenter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCustomRoleAttachmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCustomRoleAttachment: query custom role of the permission set
	var (
		getCustomRoleAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role"
		getCustomRoleAttachmentProduct = "identitycenter"
	)
	getCustomRoleAttachmentClient, err := cfg.NewServiceClient(getCustomRoleAttachmentProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	getCustomRoleAttachmentPath := getCustomRoleAttachmentClient.Endpoint + getCustomRoleAttachmentHttpUrl
	getCustomRoleAttachmentPath = strings.ReplaceAll(getCustomRoleAttachmentPath, "{instance_id}",
		state.Primary.Attributes["instance_id"])
	getCustomRoleAttachmentPath = strings.ReplaceAll(getCustomRoleAttachmentPath, "{permission_set_id}",
		state.Primary.Attributes["permission_set_id"])

	getCustomRoleAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getCustomRoleAttachmentResp, err := getCustomRoleAttachmentClient.Request("GET", getCustomRoleAttachmentPath,
		&getCustomRoleAttachmentOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center custom role attachment: %s", err)
	}

	getCustomRoleAttachmentRespBody, err := utils.FlattenResponse(getCustomRoleAttachmentResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center custom role attachment: %s", err)
	}

	customRole := utils.PathSearch("custom_role", getCustomRoleAttachmentRespBody, "").(string)
	if customRole == "" {
		return nil, fmt.Errorf("error retrieving Identity Center custom role attachment: %s", err)
	}

	return getCustomRoleAttachmentRespBody, nil
}

func TestAccCustomRoleAttachment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_custom_role_attachment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomRoleAttachmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomRoleAttachment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttrPair(rName, "permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
					resource.TestCheckResourceAttrWith(rName, "custom_role",
						checkCustomRole("{\"Version\":\"1.1\",\"Statement\":[{\"Effect\":\"Allow\","+
							"\"Action\":[\"iam:users:listUsers\",\"iam:users:getUser\"]}]}")),
				),
			},
			{
				Config: testCustomRoleAttachment_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttrPair(rName, "permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
					resource.TestCheckResourceAttrWith(rName, "custom_role",
						checkCustomRole("{\"Version\":\"1.1\",\"Statement\":[{\"Effect\":\"Allow\","+
							"\"Action\":[\"iam:users:createUser\",\"iam:users:deleteUser\"]}]}")),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPolicyAttachmentImportState(rName),
			},
		},
	})
}

func checkCustomRole(targetRole string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		var targetJson, valueJson interface{}
		if err := json.Unmarshal([]byte(targetRole), &targetJson); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &valueJson); err != nil {
			return err
		}
		if reflect.DeepEqual(targetJson, valueJson) {
			return nil
		}
		return fmt.Errorf("%#v is not equal target %#v", value, targetRole)
	}
}

func testCustomRoleAttachment_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_custom_role_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  custom_role       = jsonencode(
  {
    "Version":"1.1",
    "Statement":[
      {
        "Effect":"Allow",
        "Action":[
          "iam:users:listUsers",
          "iam:users:getUser"
        ]
      }
    ]
  })
}
`, testPermissionSet_basic(name))
}

func testCustomRoleAttachment_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_custom_role_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  custom_role       = jsonencode(
  {
    "Version":"1.1",
    "Statement":[
      {
        "Effect":"Allow",
        "Action":[
          "iam:users:createUser",
          "iam:users:deleteUser"
        ]
      }
    ]
  })
}
`, testPermissionSet_basic(name))
}
