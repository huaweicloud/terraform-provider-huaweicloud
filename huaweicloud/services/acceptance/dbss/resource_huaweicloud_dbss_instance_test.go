package dbss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dbss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query the DBSS instance detail
	var (
		getInstanceHttpUrl = "v1/{project_id}/dbss/audit/instances"
		getInstanceProduct = "dbss"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := getInstanceClient.Request("GET", getInstancePath, &getInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Instance: %s", err)
	}

	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return nil, err
	}

	instances, err := jmespath.Search("servers", getInstanceRespBody)
	if err != nil {
		return nil, fmt.Errorf("error parsing servers from response= %#v", getInstanceRespBody)
	}

	return dbss.FilterInstances(instances.([]interface{}), state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dbss_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttrPair(rName, "flavor",
						"data.huaweicloud_dbss_flavors.test", "flavors.0.id"),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", "dbss.bypassaudit.low"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charging_mode", "enterprise_project_id", "flavor", "period", "period_unit",
				},
			},
		},
	})
}

func testInstance_base() string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name   = "subnet-default"
  vpc_id = data.huaweicloud_vpc.test.id
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_availability_zones" "test" {
  region = "%s"
}

data "huaweicloud_dbss_flavors" "test" {} 
`, acceptance.HW_REGION_NAME)
}

func testInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dbss_instance" "test" {
  name               = "%s"
  description        = "terraform test"
  flavor             = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code = "dbss.bypassaudit.low"
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = data.huaweicloud_networking_secgroup.test.id
  charging_mode      = "prePaid"
  period_unit        = "month"
  period             = 1
}
`, testInstance_base(), name)
}

func TestAccInstance_updateWithEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dbss_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_withEpsId(name, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testInstance_withEpsId(name, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testInstance_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dbss_instance" "test" {
  name                   = "%s"
  description            = "terraform test"
  flavor                 = data.huaweicloud_dbss_flavors.test.flavors[0].id
  resource_spec_code     = "dbss.bypassaudit.low"
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                 = data.huaweicloud_vpc.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  charging_mode          = "prePaid"
  period_unit            = "month"
  period                 = 1
  enterprise_project_id  = "%s"
}
`, testInstance_base(), name, epsId)
}
