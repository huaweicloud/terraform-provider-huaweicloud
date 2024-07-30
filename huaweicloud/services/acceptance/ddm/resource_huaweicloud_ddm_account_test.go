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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	accountName := state.Primary.Attributes["name"]
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{project_id}", getAccountClient.ProjectID)
	getAccountPath = strings.ReplaceAll(getAccountPath, "{instance_id}", instanceID)

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
	account := utils.PathSearch(fmt.Sprintf("users|[?name=='%s']|[0]", accountName), getAccountRespBody, nil)
	if account == nil {
		return nil, fmt.Errorf("the instance %s has no account %s", instanceID, accountName)
	}
	return account, nil
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
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "permissions.0", "SELECT"),
					resource.TestCheckResourceAttr(rName, "description", "this is a test account"),
					resource.TestCheckResourceAttrPair(rName, "schemas.0.name",
						"huaweicloud_ddm_schema.test.0", "name"),
				),
			},
			{
				Config: testDdmAccount_basic_update(instanceName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "permissions.0", "CREATE"),
					resource.TestCheckResourceAttrPair(rName, "schemas.0.name",
						"huaweicloud_ddm_schema.test.1", "name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "password"},
				ImportStateIdFunc:       testDdmAccountImportState(rName),
			},
		},
	})
}

func testDdmAccount_base(instanceName, name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" test {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mysql.n1.large.4"
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  db {
    password = "test_1234"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_ddm_schema" "test" {
  count = 2

  instance_id  = huaweicloud_ddm_instance.test.id
  name         = "%[3]s_${count.index}"
  shard_mode   = "single"
  shard_number = "1"

  data_nodes {
    id             = huaweicloud_rds_instance.test.id
    admin_user     = "root"
    admin_password = "test_1234"
  }

  delete_rds_data = "true"

  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}

`, common.TestVpc(name), instanceName, name)
}

func testDdmAccount_basic(instanceName, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_account" "test" {
  depends_on = [huaweicloud_ddm_schema.test[0]]

  instance_id = huaweicloud_ddm_instance.test.id
  name        = "%s"
  password    = "test_1234"
  description = "this is a test account"

  permissions = [
    "SELECT"
  ]

  schemas {
    name = huaweicloud_ddm_schema.test[0].name
  }
}
`, testDdmAccount_base(instanceName, name), name)
}

func testDdmAccount_basic_update(instanceName, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_account" "test" {
  depends_on = [huaweicloud_ddm_schema.test[1]]

  instance_id = huaweicloud_ddm_instance.test.id
  name        = "%s"
  password    = "test_12345"
  description = ""

  permissions = [
    "CREATE"
  ]

  schemas {
    name = huaweicloud_ddm_schema.test[1].name
  }
}
`, testDdmAccount_base(instanceName, name), name)
}

func testDdmAccountImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		accountName := rs.Primary.Attributes["name"]
		return fmt.Sprintf("%s/%s", instanceId, accountName), nil
	}
}
