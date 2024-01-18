package cnad

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPolicyAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                    = acceptance.HW_REGION_NAME
		getProtectedObjectHttpUrl = "v1/cnad/protected-ips"
		getProtectedObjectProduct = "aad"
	)
	getProtectedObjectClient, err := cfg.NewServiceClient(getProtectedObjectProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CNAD Client: %s", err)
	}

	policyID := state.Primary.Attributes["policy_id"]
	instanceID := state.Primary.Attributes["instance_id"]
	queryParam := fmt.Sprintf("?policy_id=%s&package_id=%s", policyID, instanceID)
	getProtectedObjectPath := getProtectedObjectClient.Endpoint + getProtectedObjectHttpUrl
	getProtectedObjectPath += queryParam

	getProtectedObjectsResp, err := pagination.ListAllItems(
		getProtectedObjectClient,
		"offset",
		getProtectedObjectPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		// here is no special error code
		return nil, fmt.Errorf("error retrieving policy binding protected objects, %s", err)
	}

	getProtectedObjectsRespJson, err := json.Marshal(getProtectedObjectsResp)
	if err != nil {
		return nil, err
	}
	var getProtectedObjectsRespBody interface{}
	if err := json.Unmarshal(getProtectedObjectsRespJson, &getProtectedObjectsRespBody); err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("items", getProtectedObjectsRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil, fmt.Errorf("error retrieving policy binding protected objects")
	}
	return curArray, nil
}

// Due to testing environment issues, updating operations cannot be tested
func TestAccPolicyAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cnad_advanced_policy_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPolicyAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCNADInstance(t)
			acceptance.TestAccPreCheckCNADProtectedObject(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPolicyAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_cnad_advanced_policy.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_cnad_advanced_instances.test",
						"instances.0.instance_id"),
					resource.TestCheckResourceAttr(rName, "protected_object_ids.0",
						acceptance.HW_CNAD_PROJECT_OBJECT_ID),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.id"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.type"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.instance_name"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.instance_version"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.region"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.status"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.block_threshold"),
					resource.TestCheckResourceAttrSet(rName, "protected_objects.0.clean_threshold"),
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

func testPolicyAssociate_base(name string) string {
	return fmt.Sprintf(`

data "huaweicloud_cnad_advanced_instances" "test" {}

resource "huaweicloud_cnad_advanced_policy" "test" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
  name        = "%s"
}
`, name)
}

func testPolicyAssociate_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cnad_advanced_policy_associate" "test" {
  policy_id            = huaweicloud_cnad_advanced_policy.test.id
  instance_id          = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
  protected_object_ids = ["%s"]
}
`, testPolicyAssociate_base(name), acceptance.HW_CNAD_PROJECT_OBJECT_ID)
}
