package taurusdb

import (
	"errors"
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

func getTaurusDBAccountPrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getTaurusDBAccountPrivilege: Query the TaurusDB account privilege
	var (
		getTaurusDBAccountPrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		getTaurusDBAccountPrivilegeProduct = "gaussdb"
	)
	getTaurusDBAccountPrivilegeClient, err := cfg.NewServiceClient(getTaurusDBAccountPrivilegeProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 3)
	if len(parts) != 3 {
		return nil, errors.New("invalid id format, must be <instance_id>/<name>/<host>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	host := parts[2]

	getTaurusDBAccountPrivilegeBasePath := getTaurusDBAccountPrivilegeClient.Endpoint + getTaurusDBAccountPrivilegeHttpUrl
	getTaurusDBAccountPrivilegeBasePath = strings.ReplaceAll(getTaurusDBAccountPrivilegeBasePath, "{project_id}",
		getTaurusDBAccountPrivilegeClient.ProjectID)
	getTaurusDBAccountPrivilegeBasePath = strings.ReplaceAll(getTaurusDBAccountPrivilegeBasePath, "{instance_id}", instanceID)

	getTaurusDBAccountPrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var currentTotal int
	var account interface{}
	getTaurusDBAccountPrivilegePath := getTaurusDBAccountPrivilegeBasePath + buildTaurusDBQueryParams(currentTotal)

	for {
		getTaurusDBAccountResp, err := getTaurusDBAccountPrivilegeClient.Request("GET", getTaurusDBAccountPrivilegePath,
			&getTaurusDBAccountPrivilegeOpt)

		if err != nil {
			return nil, err
		}

		getTaurusDBAccountRespBody, err := utils.FlattenResponse(getTaurusDBAccountResp)
		if err != nil {
			return nil, err
		}
		res, pageNum := flattenGetTaurusDBAccountPrivilegeResponseBody(getTaurusDBAccountRespBody, accountName, host)
		if res != nil {
			account = res
			break
		}
		total := utils.PathSearch("total_count", getTaurusDBAccountRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getTaurusDBAccountPrivilegePath = getTaurusDBAccountPrivilegeBasePath + buildTaurusDBQueryParams(currentTotal)
	}
	if account == nil {
		return nil, errors.New("error retrieving TaurusDB account privilege")
	}

	databases := utils.PathSearch("databases", account, make([]interface{}, 0)).([]interface{})
	if len(databases) == 0 {
		return nil, errors.New("error retrieving TaurusDB account privilege")
	}
	return account, nil
}

func flattenGetTaurusDBAccountPrivilegeResponseBody(resp interface{}, accountName, address string) (interface{}, int) {
	if resp == nil {
		return nil, 0
	}
	curJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "").(string)
		host := utils.PathSearch("host", v, "").(string)
		if accountName == name && address == host {
			return v, len(curArray)
		}
	}
	return nil, len(curArray)
}

func TestAccTaurusDBAccountPrivilege_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_taurusdb_account_privilege.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTaurusDBAccountPrivilegeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTaurusDBAccountPrivilege_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "account_name",
						"huaweicloud_taurusdb_account.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "host",
						"huaweicloud_taurusdb_account.test", "host"),
					resource.TestCheckResourceAttrPair(rName, "databases.0.name",
						"huaweicloud_taurusdb_database.test", "name"),
					resource.TestCheckResourceAttr(rName, "databases.0.readonly", "false"),
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

func testTaurusDBAccountPrivilege_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_account" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
  host        = "10.10.10.10"
}

resource "huaweicloud_taurusdb_database" "test" {
  instance_id   = huaweicloud_taurusdb_instance.test.id
  name          = "%[2]s"
  character_set = "gbk"
}

resource "huaweicloud_taurusdb_account_privilege" "test" {
  instance_id  = huaweicloud_taurusdb_instance.test.id
  account_name = huaweicloud_taurusdb_account.test.name
  host         = huaweicloud_taurusdb_account.test.host

  databases {
    name     = huaweicloud_taurusdb_database.test.name
    readonly = false
  }
}
`, testAccTaurusDBInstanceConfig_basic(name), name)
}
