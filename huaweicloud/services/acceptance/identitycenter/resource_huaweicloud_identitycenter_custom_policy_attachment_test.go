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

func getCustomPolicyAttachmentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCustomPolicyAttachment: query custom policy of the permission set
	var (
		getCustomPolicyAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy"
		getCustomPolicyAttachmentProduct = "identitycenter"
	)
	getCustomPolicyAttachmentClient, err := cfg.NewServiceClient(getCustomPolicyAttachmentProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	getCustomPolicyAttachmentPath := getCustomPolicyAttachmentClient.Endpoint + getCustomPolicyAttachmentHttpUrl
	getCustomPolicyAttachmentPath = strings.ReplaceAll(getCustomPolicyAttachmentPath, "{instance_id}",
		state.Primary.Attributes["instance_id"])
	getCustomPolicyAttachmentPath = strings.ReplaceAll(getCustomPolicyAttachmentPath, "{permission_set_id}",
		state.Primary.Attributes["permission_set_id"])

	getCustomPolicyAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getCustomPolicyAttachmentResp, err := getCustomPolicyAttachmentClient.Request("GET",
		getCustomPolicyAttachmentPath, &getCustomPolicyAttachmentOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center custom policy attachment: %s", err)
	}

	getCustomPolicyAttachmentRespBody, err := utils.FlattenResponse(getCustomPolicyAttachmentResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center custom policy attachment: %s", err)
	}

	customPolicy := utils.PathSearch("custom_policy", getCustomPolicyAttachmentRespBody, "").(string)
	if customPolicy == "" {
		return nil, fmt.Errorf("error retrieving Identity Center custom policy attachment: %s", err)
	}

	return getCustomPolicyAttachmentRespBody, nil
}

func TestAccCustomPolicyAttachment_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_custom_policy_attachment.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomPolicyAttachmentResourceFunc,
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
				Config: testCustomPolicyAttachment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttrPair(rName, "permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
					resource.TestCheckResourceAttrWith(rName, "custom_policy",
						checkCustomPolicy("{\"Version\":\"5.0\",\"Statement\":[{\"Effect\":\"Allow\","+
							"\"Action\":[\"billing:subscription:renew\"]}]}")),
				),
			},
			{
				Config: testCustomPolicyAttachment_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttrPair(rName, "permission_set_id",
						"huaweicloud_identitycenter_permission_set.test", "id"),
					resource.TestCheckResourceAttrWith(rName, "custom_policy",
						checkCustomPolicy("{\"Version\":\"5.0\",\"Statement\":[{\"Effect\":\"Allow\","+
							"\"Action\":[\"eps:resources:add\"]}]}")),
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

func checkCustomPolicy(targetPolicy string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		var targetJson, valueJson interface{}
		if err := json.Unmarshal([]byte(targetPolicy), &targetJson); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &valueJson); err != nil {
			return err
		}
		if reflect.DeepEqual(targetJson, valueJson) {
			return nil
		}
		return fmt.Errorf("%#v is not equal target %#v", value, targetPolicy)
	}
}

func testCustomPolicyAttachment_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_custom_policy_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  custom_policy     = jsonencode(
{
	"Version":"5.0",
	"Statement":[
		{
			"Effect":"Allow",
			"Action":[
				"billing:subscription:renew"
			]
		}
	]
}
)
}
`, testPermissionSet_basic(name))
}

func testCustomPolicyAttachment_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_custom_policy_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
  custom_policy     = jsonencode(
{
	"Version":"5.0",
	"Statement":[
		{
			"Effect":"Allow",
			"Action":[
				"eps:resources:add"
			]
		}
	]
}
)
}
`, testPermissionSet_basic(name))
}
