package er

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInstanceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getInstance: Query the Enterprise router instance detail
	var (
		getInstanceHttpUrl = "v3/{project_id}/enterprise-router/instances/{id}"
		getInstanceProduct = "er"
	)
	getInstanceClient, err := cfg.NewServiceClient(getInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Instance Client: %s", err)
	}

	getInstancePath := getInstanceClient.Endpoint + getInstanceHttpUrl
	getInstancePath = strings.ReplaceAll(getInstancePath, "{project_id}", getInstanceClient.ProjectID)
	getInstancePath = strings.ReplaceAll(getInstancePath, "{id}", state.Primary.ID)

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
	return utils.FlattenResponse(getInstanceResp)
}

func TestAccInstance_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_er_instance.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic_step1(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_er_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "enable_default_propagation", "true"),
					resource.TestCheckResourceAttr(rName, "enable_default_association", "true"),
					resource.TestCheckResourceAttr(rName, "auto_accept_shared_attachments", "true"),
					resource.TestCheckResourceAttrSet(rName, "default_propagation_route_table_id"),
					resource.TestCheckResourceAttrSet(rName, "default_association_route_table_id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testInstance_basic_step2(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_er_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(rName, "asn", fmt.Sprintf("%v", bgpAsNum)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "enable_default_propagation", "false"),
					resource.TestCheckResourceAttr(rName, "enable_default_association", "false"),
					resource.TestCheckResourceAttr(rName, "auto_accept_shared_attachments", "false"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(rName, "tags.newkey", "value"),
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

func testInstance_basic_step1(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name        = "%[1]s"
  asn         = %[2]d
  description = "Created by script"

  enable_default_propagation     = true
  enable_default_association     = true
  auto_accept_shared_attachments = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, bgpAsNum)
}

func testInstance_basic_step2(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 1, 2)

  name = "%[1]s"
  asn  = %[2]d

  enable_default_propagation     = false
  enable_default_association     = false
  auto_accept_shared_attachments = false

  tags = {
    foo    = "baar"
    newkey = "value"
  }
}
`, name, bgpAsNum)
}

func TestAccInstance_updateWithEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_er_instance.test"
	bgpAsNum := acctest.RandIntRange(64512, 65534)

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
				Config: testInstance_withEpsId(name, bgpAsNum, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testInstance_withEpsId(name, bgpAsNum, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testInstance_withEpsId(name string, bgpAsNum int, epsId string) string {
	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = slice(data.huaweicloud_er_availability_zones.test.names, 0, 1)

  name                  = "%[1]s"
  asn                   = %[2]d
  enterprise_project_id = "%[3]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, bgpAsNum, epsId)
}
