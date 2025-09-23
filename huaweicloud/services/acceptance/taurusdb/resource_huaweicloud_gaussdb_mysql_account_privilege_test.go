package taurusdb

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

func getGaussDBAccountPrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGaussDBAccountPrivilege: Query the GaussDB MySQL account privilege
	var (
		getGaussDBAccountPrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		getGaussDBAccountPrivilegeProduct = "gaussdb"
	)
	getGaussDBAccountPrivilegeClient, err := cfg.NewServiceClient(getGaussDBAccountPrivilegeProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<name>/<host>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	host := parts[2]

	getGaussDBAccountPrivilegeBasePath := getGaussDBAccountPrivilegeClient.Endpoint + getGaussDBAccountPrivilegeHttpUrl
	getGaussDBAccountPrivilegeBasePath = strings.ReplaceAll(getGaussDBAccountPrivilegeBasePath, "{project_id}",
		getGaussDBAccountPrivilegeClient.ProjectID)
	getGaussDBAccountPrivilegeBasePath = strings.ReplaceAll(getGaussDBAccountPrivilegeBasePath, "{instance_id}", instanceID)

	getGaussDBAccountPrivilegeOpt := golangsdk.RequestOpts{
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
	getGaussDBAccountPrivilegePath := getGaussDBAccountPrivilegeBasePath + buildGaussDBMysqlQueryParams(currentTotal)

	for {
		getGaussDBAccountResp, err := getGaussDBAccountPrivilegeClient.Request("GET", getGaussDBAccountPrivilegePath,
			&getGaussDBAccountPrivilegeOpt)

		if err != nil {
			return nil, err
		}

		getGaussDBAccountRespBody, err := utils.FlattenResponse(getGaussDBAccountResp)
		if err != nil {
			return nil, err
		}
		res, pageNum := flattenGetGaussDBAccountPrivilegeResponseBody(getGaussDBAccountRespBody, accountName, host)
		if res != nil {
			account = res
			break
		}
		total := utils.PathSearch("total_count", getGaussDBAccountRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getGaussDBAccountPrivilegePath = getGaussDBAccountPrivilegeBasePath + buildGaussDBMysqlQueryParams(currentTotal)
	}
	if account == nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL account privilege")
	}

	databases := utils.PathSearch("databases", account, make([]interface{}, 0)).([]interface{})
	if len(databases) == 0 {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL account privilege")
	}
	return account, nil
}

func flattenGetGaussDBAccountPrivilegeResponseBody(resp interface{}, accountName, address string) (interface{}, int) {
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

func TestAccGaussDBAccountPrivilege_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_mysql_account_privilege.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBAccountPrivilegeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDBAccountPrivilege_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "account_name",
						"huaweicloud_gaussdb_mysql_account.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "host",
						"huaweicloud_gaussdb_mysql_account.test", "host"),
					resource.TestCheckResourceAttrPair(rName, "databases.0.name",
						"huaweicloud_gaussdb_mysql_database.test", "name"),
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

func testGaussDBAccountPrivilege_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_mysql_account" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
  host        = "10.10.10.10"
}

resource "huaweicloud_gaussdb_mysql_database" "test" {
  instance_id   = huaweicloud_gaussdb_mysql_instance.test.id
  name          = "%[2]s"
  character_set = "gbk"
}

resource "huaweicloud_gaussdb_mysql_account_privilege" "test" {
  instance_id  = huaweicloud_gaussdb_mysql_instance.test.id
  account_name = huaweicloud_gaussdb_mysql_account.test.name
  host         = huaweicloud_gaussdb_mysql_account.test.host

  databases {
    name     = huaweicloud_gaussdb_mysql_database.test.name
    readonly = false
  }
}
`, testAccGaussDBInstanceConfig_basic(name), name)
}
