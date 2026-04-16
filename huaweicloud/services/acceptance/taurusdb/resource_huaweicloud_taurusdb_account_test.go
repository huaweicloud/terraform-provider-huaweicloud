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

func getTaurusDBAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getTaurusDBAccount: Query the TaurusDB account
	var (
		getTaurusDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		getTaurusDBAccountProduct = "gaussdb"
	)
	getTaurusDBAccountClient, err := cfg.NewServiceClient(getTaurusDBAccountProduct, region)
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

	getTaurusDBAccountBasePath := getTaurusDBAccountClient.Endpoint + getTaurusDBAccountHttpUrl
	getTaurusDBAccountBasePath = strings.ReplaceAll(getTaurusDBAccountBasePath, "{project_id}",
		getTaurusDBAccountClient.ProjectID)
	getTaurusDBAccountBasePath = strings.ReplaceAll(getTaurusDBAccountBasePath, "{instance_id}", instanceID)

	getTaurusDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var currentTotal int
	var account interface{}
	getTaurusDBAccountPath := getTaurusDBAccountBasePath + buildTaurusDBQueryParams(currentTotal)

	for {
		getTaurusDBAccountResp, err := getTaurusDBAccountClient.Request("GET", getTaurusDBAccountPath,
			&getTaurusDBAccountOpt)

		if err != nil {
			return nil, err
		}

		getTaurusDBAccountRespBody, err := utils.FlattenResponse(getTaurusDBAccountResp)
		if err != nil {
			return nil, err
		}
		res, pageNum := flattenGetTaurusDBAccountResponseBody(getTaurusDBAccountRespBody, accountName, host)
		if res != nil {
			account = res
			break
		}
		total := utils.PathSearch("total_count", getTaurusDBAccountRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getTaurusDBAccountPath = getTaurusDBAccountBasePath + buildTaurusDBQueryParams(currentTotal)
	}
	if account == nil {
		return nil, errors.New("error retrieving TaurusDB account")
	}
	return account, nil
}

func flattenGetTaurusDBAccountResponseBody(resp interface{}, accountName, address string) (interface{}, int) {
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

func TestAccTaurusDBAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_taurusdb_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTaurusDBAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTaurusDBAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "host", "%"),
					resource.TestCheckResourceAttr(rName, "description",
						"test for gaussdb mysql description"),
				),
			},
			{
				Config: testTaurusDBAccount_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
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

func testTaurusDBAccount_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_taurusdb_account" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
  name        = "%s"
  password    = "Test@12345678"
  description = "test for gaussdb mysql description"
}
`, testAccTaurusDBInstanceConfig_basic(name), name)
}

func testTaurusDBAccount_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_taurusdb_account" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
  name        = "%s"
  password    = "Test@123456789"
  description = "test for gaussdb mysql description update"
}
`, testAccTaurusDBInstanceConfig_basic(name), name)
}
