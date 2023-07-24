package mrs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDataConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getDataConnectionHttpUrl = "v2/{project_id}/data-connectors"
		getDataConnectionProduct = "mrs"
	)
	getDataConnectionClient, err := cfg.NewServiceClient(getDataConnectionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating MRS Client: %s", err)
	}

	getDataConnectionPath := getDataConnectionClient.Endpoint + getDataConnectionHttpUrl
	getDataConnectionPath = strings.ReplaceAll(getDataConnectionPath, "{project_id}", getDataConnectionClient.ProjectID)

	getDataConnectionqueryParams := fmt.Sprintf("?id=%v", state.Primary.ID)
	getDataConnectionPath += getDataConnectionqueryParams

	getDataConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getDataConnectionResp, err := getDataConnectionClient.Request("GET", getDataConnectionPath, &getDataConnectionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data connection: %s", err)
	}

	getDataConnectionRespBody, err := utils.FlattenResponse(getDataConnectionResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data connection: %s", err)
	}

	jsonPath := fmt.Sprintf("data_connectors[?connector_id =='%s']|[0]", state.Primary.ID)
	getDataConnectionRespBody = utils.PathSearch(jsonPath, getDataConnectionRespBody, nil)
	if getDataConnectionRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getDataConnectionRespBody, nil
}

func TestAccDataConnection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_mapreduce_data_connection.test"
	password := acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataConnection_basic(name, "root", password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source_type", "RDS_MYSQL"),
					resource.TestCheckResourceAttr(rName, "source_info.0.db_name", name),
					resource.TestCheckResourceAttr(rName, "source_info.0.user_name", "root"),
					resource.TestCheckResourceAttrPair(rName, "source_info.0.db_instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
				),
			},
			{
				Config: testDataConnection_basic(name, "root_2", password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source_type", "RDS_MYSQL"),
					resource.TestCheckResourceAttr(rName, "source_info.0.db_name", name),
					resource.TestCheckResourceAttr(rName, "source_info.0.user_name", "root_2"),
					resource.TestCheckResourceAttrPair(rName, "source_info.0.db_instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_info.0.password"},
			},
		},
	})
}

func testDataConnection_basic(name, userName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_data_connection" "test" {
  name        = "%s"
  source_type = "RDS_MYSQL"
  source_info {
    db_instance_id = huaweicloud_rds_instance.test.id
    db_name        = huaweicloud_rds_mysql_database.test.name
    user_name      = "%s"
    password       = "%s"
  }
}
`, testRdsAccount_base(name), name, userName, password)
}

func testRdsAccount_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type           = "MySQL"
  group_type        = "general"
  db_version        = "5.7"
  instance_mode     = "single"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors.0.name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.58"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%s"
  character_set = "utf8"
}
`, common.TestBaseNetwork(name), name, name)
}
