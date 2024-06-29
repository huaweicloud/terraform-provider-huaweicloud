package dcs

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

func getDcsAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAccount: get acount
	var (
		getAccountHttpUrl = "v2/{project_id}/instances/{instance_id}/accounts"
		getAccountProduct = "dcs"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", getAccountClient.ProjectID)
	getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", instanceId)

	getAccountOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getAccountResp, err := getAccountClient.Request("GET", getAccountPath, &getAccountOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS accounts: %s", err)
	}

	getAccountRespBody, err := utils.FlattenResponse(getAccountResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS accounts: %s", err)
	}
	account := utils.PathSearch(fmt.Sprintf("accounts|[?account_id =='%s']|[0]", state.Primary.ID), getAccountRespBody, nil)
	if account == nil {
		return nil, fmt.Errorf("the account (%s) is not found", state.Primary.ID)
	}

	return account, nil
}

func TestAccDcsAccount_basic(t *testing.T) {
	var obj interface{}

	var accountName = acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDcsAccount_basic(accountName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "account_name", accountName),
					resource.TestCheckResourceAttr(rName, "account_role", "read"),
					resource.TestCheckResourceAttr(rName, "account_type", "normal"),
					resource.TestCheckResourceAttr(rName, "description", "add account"),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.instance_1", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAccDcsAccount_updated(accountName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "account_name", accountName),
					resource.TestCheckResourceAttr(rName, "account_role", "write"),
					resource.TestCheckResourceAttr(rName, "account_type", "normal"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.instance_1", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
				ImportStateIdFunc:       testDcsAccountImportState(rName),
			},
		},
	})
}

func testAccDcsAccount_basic(accountName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_account" "test" {
  instance_id      = huaweicloud_dcs_instance.instance_1.id
  account_name     = "%s"
  account_role     = "read"
  account_password = "Terraform@123"
  description      = "add account"
}`, testAccDcsV1Instance_basic(accountName), accountName)
}

func testAccDcsAccount_updated(accountName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_account" "test" {
  instance_id      = huaweicloud_dcs_instance.instance_1.id
  account_name     = "%s"
  account_role     = "write"
  account_password = "Terraform@1234"
  description      = ""
}`, testAccDcsV1Instance_basic(accountName), accountName)
}

func testDcsAccountImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceId, rs.Primary.ID), nil
	}
}
