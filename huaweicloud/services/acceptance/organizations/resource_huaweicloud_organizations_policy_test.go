package organizations

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

func getPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPolicy: Query Organizations policy
	var (
		getPolicyHttpUrl = "v1/organizations/policies/{policy_id}"
		getPolicyProduct = "organizations"
	)
	getPolicyClient, err := cfg.NewServiceClient(getPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	getPolicyPath := getPolicyClient.Endpoint + getPolicyHttpUrl
	getPolicyPath = strings.ReplaceAll(getPolicyPath, "{policy_id}", state.Primary.ID)

	getPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPolicyResp, err := getPolicyClient.Request("GET", getPolicyPath, &getPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy: %s", err)
	}

	getPolicyRespBody, err := utils.FlattenResponse(getPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations policy: %s", err)
	}

	return getPolicyRespBody, nil
}

func TestAccPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_organizations_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrWith(rName, "content",
						checkPolicyContent("{\"Version\":\"5.0\",\"Statement\":[{\"Effect\":\"Deny\","+
							"\"Action\":[]}]}")),
					resource.TestCheckResourceAttr(rName, "type", "service_control_policy"),
					resource.TestCheckResourceAttr(rName, "description",
						"test service control policy description"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
				),
			},
			{
				Config: testPolicy_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttrWith(rName, "content",
						checkPolicyContent("{\"Version\":\"5.0\",\"Statement\":[{\"Sid\":\"Statement1\","+
							"\"Effect\":\"Allow\",\"Action\":[\"*\"]}]}")),
					resource.TestCheckResourceAttr(rName, "type", "service_control_policy"),
					resource.TestCheckResourceAttr(rName, "description",
						"test service control policy description update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
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

func checkPolicyContent(targetContent string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		var targetJson, valueJson interface{}
		if err := json.Unmarshal([]byte(targetContent), &targetJson); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value), &valueJson); err != nil {
			return err
		}
		if reflect.DeepEqual(targetJson, valueJson) {
			return nil
		}
		return fmt.Errorf("%#v is not equal target %#v", value, targetContent)
	}
}

func testPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name        = "%s"
  type        = "service_control_policy"
  description = "test service control policy description"
  content     = jsonencode(
{
	"Version":"5.0",
	"Statement":[
		{
			"Effect":"Deny",
			"Action":[]
		}
	]
}
)

  tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}
`, name)
}

func testPolicy_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name        = "%s"
  type        = "service_control_policy"
  description = "test service control policy description update"
  content     = jsonencode(
{
	"Version":"5.0",
	"Statement":[
		{
			"Sid":"Statement1",
			"Effect":"Allow",
			"Action":["*"]
		}
	]
}
)

  tags = {
    "key3" = "value3"
    "key4" = "value4"
  }
}
`, name)
}
