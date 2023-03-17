package ddm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDdmAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAccount: Query DDM account
	var (
		getAccountHttpUrl = "v1/{project_id}/instances/{instance_id}/users"
		getAccountProduct = "ddm"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<account_name>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", getAccountClient.ProjectID)
	getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", fmt.Sprintf("%v", instanceID))

	getAccountResp, err := pagination.ListAllItems(
		getAccountClient,
		"offset",
		getAccountPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDM account: %s", err)
	}

	getAccountRespJson, err := json.Marshal(getAccountResp)
	if err != nil {
		return nil, err
	}
	var getAccountRespBody interface{}
	err = json.Unmarshal(getAccountRespJson, &getAccountRespBody)
	if err != nil {
		return nil, err
	}

	accounts := utils.PathSearch("users", getAccountRespBody, nil)
	if accounts == nil {
		return nil, fmt.Errorf("the instance %s has no account", instanceID)
	}
	for _, account := range accounts.([]interface{}) {
		name := utils.PathSearch("name", account, nil)
		if accountName != name {
			continue
		}
		return account, nil
	}

	return nil, fmt.Errorf("the instance %s has no account %s", instanceID, accountName)
}

func TestAccDdmAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	instanceName := strings.ReplaceAll(name, "_", "-")
	rName := "huaweicloud_ddm_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmAccount_basic(instanceName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "permissions.0", "SELECT"),
				),
			},
			{
				Config: testDdmAccount_basic_update(instanceName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "permissions.0", "CREATE"),
					resource.TestCheckResourceAttr(rName, "description", "this is a test account"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "password"},
			},
		},
	})
}

func testDdmAccount_basic(instanceName, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_account" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
  name        = "%s"
  password    = "test_1234"

  permissions = [
    "SELECT"
 ]
}
`, testDdmInstance_basic(instanceName), name)
}

func testDdmAccount_basic_update(instanceName, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_account" "test" {
  instance_id = huaweicloud_ddm_instance.test.id
  name        = "%s"
  password    = "test_12345"
  description = "this is a test account"

  permissions = [
    "CREATE"
 ]
}
`, testDdmInstance_basic(instanceName), name)
}
