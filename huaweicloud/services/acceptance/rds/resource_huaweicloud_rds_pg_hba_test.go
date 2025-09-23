package rds

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

func getPgHbaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPgHba: query RDS PostgreSQL hba
	var (
		getPgHbaHttpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		getPgHbaProduct = "rds"
	)
	getPgHbaClient, err := cfg.NewServiceClient(getPgHbaProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPgHbaPath := getPgHbaClient.Endpoint + getPgHbaHttpUrl
	getPgHbaPath = strings.ReplaceAll(getPgHbaPath, "{project_id}", getPgHbaClient.ProjectID)
	getPgHbaPath = strings.ReplaceAll(getPgHbaPath, "{instance_id}", state.Primary.ID)

	getPgHbaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPgHbaResp, err := getPgHbaClient.Request("GET", getPgHbaPath, &getPgHbaOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL hba: %s", err)
	}

	getPgHbaRespBody, err := utils.FlattenResponse(getPgHbaResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL hba: %s", err)
	}

	curArray := getPgHbaRespBody.([]interface{})
	for _, v := range curArray {
		if !checkDefaultRule(v) {
			return getPgHbaRespBody, nil
		}
	}

	// if all rules are default rule, then it indicates the resource has been destroyed
	return nil, fmt.Errorf("error retrieving RDS PostgreSQL hba: %s", err)
}

func checkDefaultRule(rule interface{}) bool {
	ruleType := utils.PathSearch("type", rule, nil)
	ruleDatabase := utils.PathSearch("database", rule, nil)
	ruleUser := utils.PathSearch("user", rule, nil)
	ruleAddress := utils.PathSearch("address", rule, nil)
	ruleMethod := utils.PathSearch("method", rule, nil)
	if ruleType != "host" && ruleType != "hostssl" {
		return false
	}
	if !((ruleDatabase == "all" && ruleUser == "all") || (ruleDatabase == "replication" && ruleUser == "root")) {
		return false
	}
	if ruleAddress != "0.0.0.0/0" {
		return false
	}
	if ruleMethod != "md5" {
		return false
	}
	return true
}

func TestAccPgHba_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_hba.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgHbaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgHba_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.type", "host"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.database", "all"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.user", "all"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.method", "md5"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.type", "host"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.database", "postgres"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.user", "root"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.method", "scram-sha-256"),
				),
			},
			{
				Config: testPgHba_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.type", "host"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.database", "postgres"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.user", "root"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.address", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.0.method", "md5"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.type", "hostnossl"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.database", "all"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.user", "all"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.address", "0.0.0.0/1"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.1.method", "scram-sha-256"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.2.type", "hostssl"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.2.database", "all"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.2.user", "all"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.2.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.2.method", "reject"),
					resource.TestCheckResourceAttr(rName, "host_based_authentications.2.mask", "0.0.0.0"),
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

func testPgHba_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  db {
    type    = "PostgreSQL"
    version = "16"
  }
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testPgHba_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_hba" "test" {
  instance_id = huaweicloud_rds_instance.test.id

  host_based_authentications {
    type     = "host"
    database = "all"
    user     = "all"
    address  = "0.0.0.0/0"
    method   = "md5"
  }

  host_based_authentications {
    type     = "host"
    database = "postgres"
    user     = "root"
    address  = "0.0.0.0/0"
    method   = "scram-sha-256"
  }
}
`, testPgHba_base(name))
}

func testPgHba_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_hba" "test" {
  instance_id = huaweicloud_rds_instance.test.id

  host_based_authentications {
    type     = "host"
    database = "postgres"
    user     = "root"
    address  = "0.0.0.0/0"
    method   = "md5"
  }

  host_based_authentications {
    type     = "hostnossl"
    database = "all"
    user     = "all"
    address  = "0.0.0.0/1"
    method   = "scram-sha-256"
  }

  host_based_authentications {
    type     = "hostssl"
    database = "all"
    user     = "all"
    address  = "0.0.0.0"
    mask     = "0.0.0.0"
    method   = "reject"
  }
}
`, testPgHba_base(name))
}
