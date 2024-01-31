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

func getProtectedObjectResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                    = acceptance.HW_REGION_NAME
		getProtectedObjectHttpUrl = "v1/cnad/protected-ips"
		getProtectedObjectProduct = "aad"
	)
	getProtectedObjectClient, err := cfg.NewServiceClient(getProtectedObjectProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CNAD Client: %s", err)
	}

	queryParam := fmt.Sprintf("?package_id=%v", state.Primary.ID)
	getProtectedObjectPath := getProtectedObjectClient.Endpoint + getProtectedObjectHttpUrl
	getProtectedObjectPath += queryParam

	getProtectedObjectsResp, err := pagination.ListAllItems(
		getProtectedObjectClient,
		"offset",
		getProtectedObjectPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		// here is no special error code
		return nil, fmt.Errorf("error retrieving protected objects, %s", err)
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
		return nil, fmt.Errorf("error retrieving advanced protected objects")
	}
	return curArray, nil
}

// Due to testing environment issues, updating operations cannot be tested
func TestAccProtectedObject_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_cnad_advanced_protected_object.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectedObjectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCNADInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProtectedObject_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"data.huaweicloud_cnad_advanced_instances.test",
						"instances.0.instance_id"),
					resource.TestCheckResourceAttrPair(rName, "protected_objects.0.id",
						"data.huaweicloud_cnad_advanced_available_objects.test",
						"protected_objects.0.id"),
					resource.TestCheckResourceAttrPair(rName, "protected_objects.0.ip_address",
						"data.huaweicloud_cnad_advanced_available_objects.test",
						"protected_objects.0.ip_address"),
					resource.TestCheckResourceAttrPair(rName, "protected_objects.0.type",
						"data.huaweicloud_cnad_advanced_available_objects.test",
						"protected_objects.0.type"),
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

const testProtectedObject_base = `
data "huaweicloud_cnad_advanced_instances" "test" {}

data "huaweicloud_cnad_advanced_available_objects" "test" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
}
`

func testProtectedObject_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cnad_advanced_protected_object" "test" {
  instance_id = data.huaweicloud_cnad_advanced_instances.test.instances.0.instance_id
  protected_objects {
    id         = data.huaweicloud_cnad_advanced_available_objects.test.protected_objects.0.id
    ip_address = data.huaweicloud_cnad_advanced_available_objects.test.protected_objects.0.ip_address
    type       = data.huaweicloud_cnad_advanced_available_objects.test.protected_objects.0.type
  }
}
`, testProtectedObject_base)
}
