package gaussdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDbAspCollectResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/asp"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}
	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody any
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}

	asp := utils.PathSearch(fmt.Sprintf("asp[?job_id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if asp == nil {
		return nil, errors.New("error retrieving GaussDB ASP collect")
	}

	return asp, nil
}

func TestAccResourceGaussDbAspCollect_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_asp_collect.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDbAspCollectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGaussdbAspCollect_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "start_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "file_size"),
					resource.TestCheckResourceAttrSet(rName, "download_url"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "file_path"),
					resource.TestCheckResourceAttrSet(rName, "file_name"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.#"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.name"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.type"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.url"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.port"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.domain_id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccGaussDbAspImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{"download_url"},
			},
		},
	})
}

func testAccResourceGaussdbAspCollect_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "test_123456"
  port              = "9000"
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "enterprise"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, common.TestVpc(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccResourceGaussdbAspCollect_basic(name string) string {
	startTime := time.Now().UTC().Format("2006-01-02T15:04:05+0000")
	endTime := time.Now().UTC().Add(1 * time.Hour).Format("2006-01-02T15:04:05+0000")
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_asp_collect" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, testAccResourceGaussdbAspCollect_base(name), startTime, endTime)
}

func testAccGaussDbAspImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
