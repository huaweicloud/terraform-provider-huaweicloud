package gaussdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDBInstanceDatabaseAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("opengauss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	accountName := state.Primary.Attributes["name"]

	httpUrl := "v3/{project_id}/instances/{instance_id}/db-users"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB instance database accounts: %s", err)
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err), nil
	}

	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err), nil
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getRespBody, nil)
	if account == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return account, nil
}

func TestAccGaussDBInstanceDatabaseAccount_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_instance_database_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBInstanceDatabaseAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBInstanceDatabaseAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "is_login_only", "false"),
					resource.TestCheckResourceAttrSet(rName, "attribute.#"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolsuper"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolinherit"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolcreaterole"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolcreatedb"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolcanlogin"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolconnlimit"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolreplication"),
					resource.TestCheckResourceAttrSet(rName, "attribute.0.rolbypassrls"),
					resource.TestCheckResourceAttrSet(rName, "memberof"),
					resource.TestCheckResourceAttrSet(rName, "lock_status"),
				),
			},
			{
				Config: testAccGaussDBInstanceDatabaseAccount_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "is_login_only", "false"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccGaussDBInstanceDatabaseAccountImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{"password", "is_login_only"},
			},
		},
	})
}

func testAccGaussDBInstanceDatabaseAccount_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance_database_account" "test" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  name          = "%[2]s"
  password      = "Test@963852"
  is_login_only = "false"
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}

func testAccGaussDBInstanceDatabaseAccount_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_instance_database_account" "test" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  name          = "%[2]s"
  password      = "Asd@741258"
  is_login_only = "false"
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}

func testAccGaussDBInstanceDatabaseAccountImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		name := rs.Primary.Attributes["name"]
		if instanceID == "" || name == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<name>', but '%s/%s'",
				instanceID, name)
		}
		return fmt.Sprintf("%s/%s", instanceID, name), nil
	}
}
