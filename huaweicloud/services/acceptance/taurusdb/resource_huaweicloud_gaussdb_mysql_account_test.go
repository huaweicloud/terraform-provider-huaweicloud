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

func getGaussDBAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGaussDBAccount: Query the GaussDB MySQL account
	var (
		getGaussDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		getGaussDBAccountProduct = "gaussdb"
	)
	getGaussDBAccountClient, err := cfg.NewServiceClient(getGaussDBAccountProduct, region)
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

	getGaussDBAccountBasePath := getGaussDBAccountClient.Endpoint + getGaussDBAccountHttpUrl
	getGaussDBAccountBasePath = strings.ReplaceAll(getGaussDBAccountBasePath, "{project_id}",
		getGaussDBAccountClient.ProjectID)
	getGaussDBAccountBasePath = strings.ReplaceAll(getGaussDBAccountBasePath, "{instance_id}", instanceID)

	getGaussDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var currentTotal int
	var account interface{}
	getGaussDBAccountPath := getGaussDBAccountBasePath + buildGaussDBMysqlQueryParams(currentTotal)

	for {
		getGaussDBAccountResp, err := getGaussDBAccountClient.Request("GET", getGaussDBAccountPath,
			&getGaussDBAccountOpt)

		if err != nil {
			return nil, err
		}

		getGaussDBAccountRespBody, err := utils.FlattenResponse(getGaussDBAccountResp)
		if err != nil {
			return nil, err
		}
		res, pageNum := flattenGetGaussDBAccountResponseBody(getGaussDBAccountRespBody, accountName, host)
		if res != nil {
			account = res
			break
		}
		total := utils.PathSearch("total_count", getGaussDBAccountRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getGaussDBAccountPath = getGaussDBAccountBasePath + buildGaussDBMysqlQueryParams(currentTotal)
	}
	if account == nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL account")
	}
	return account, nil
}

func flattenGetGaussDBAccountResponseBody(resp interface{}, accountName, address string) (interface{}, int) {
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

func TestAccGaussDBAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_mysql_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDBAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "host", "%"),
					resource.TestCheckResourceAttr(rName, "description",
						"test for gaussdb mysql description"),
				),
			},
			{
				Config: testGaussDBAccount_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "host", "%"),
					resource.TestCheckResourceAttr(rName, "description",
						"test for gaussdb mysql description update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testGaussDBAccount_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_mysql_account" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  name        = "%s"
  password    = "Test@12345678"
  description = "test for gaussdb mysql description"
}
`, testAccGaussDBInstanceConfig_basic(name), name)
}

func testGaussDBAccount_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_mysql_account" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  name        = "%s"
  password    = "Test@123456789"
  description = "test for gaussdb mysql description update"
}
`, testAccGaussDBInstanceConfig_basic(name), name)
}
