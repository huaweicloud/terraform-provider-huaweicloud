package dds

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

func getDdsParameterTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getParameterTemplate: Query DDS parameter template
	var (
		getParameterTemplateHttpUrl = "v3/{project_id}/configurations/{config_id}"
		getParameterTemplateProduct = "dds"
	)
	getParameterTemplateClient, err := cfg.NewServiceClient(getParameterTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS Client: %s", err)
	}

	getParameterTemplatePath := getParameterTemplateClient.Endpoint + getParameterTemplateHttpUrl
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{project_id}",
		getParameterTemplateClient.ProjectID)
	getParameterTemplatePath = strings.ReplaceAll(getParameterTemplatePath, "{config_id}", state.Primary.ID)

	getParameterTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getParameterTemplateResp, err := getParameterTemplateClient.Request("GET",
		getParameterTemplatePath, &getParameterTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDS parameter template: %s", err)
	}
	return utils.FlattenResponse(getParameterTemplateResp)
}

func TestAccDdsParameterTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dds_parameter_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsParameterTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsParameterTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name",
						"connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "800"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name",
						"connPoolMaxShardedConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "800"),
				),
			},
			{
				Config: testDdsParameterTemplate_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "node_version", "4.0"),
					resource.TestCheckResourceAttr(rName, "parameters.0.name",
						"connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "500"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name",
						"connPoolMaxShardedConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "500"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_type", "parameter_values"},
			},
		},
	})
}

func testDdsParameterTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_parameter_template" "test" {
  name         = "%s"
  description  = "test description"
  node_type    = "mongos"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost        = 800
    connPoolMaxShardedConnsPerHost = 800
  }
}
`, name)
}

func testDdsParameterTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_parameter_template" "test" {
  name         = "%s"
  description  = ""
  node_type    = "mongos"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost        = 500
    connPoolMaxShardedConnsPerHost = 500
  }
}
`, name)
}

func TestAccDdsParameterTemplate_withEntityId(t *testing.T) {
	var (
		obj        interface{}
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		rName      = "huaweicloud_dds_parameter_template.test"

		rc = acceptance.InitResourceCheck(
			rName,
			&obj,
			getDdsParameterTemplateResourceFunc,
		)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDdsParameterTemplate_withEntityId_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "entity_id", "huaweicloud_dds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
				),
			},
			{
				Config: testAccDdsParameterTemplate_withEntityId_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "parameters.0.value", "800"),
					resource.TestCheckResourceAttr(rName, "parameters.1.name", "cursorTimeoutMillis"),
					resource.TestCheckResourceAttr(rName, "parameters.1.value", "800000"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"entity_id", "parameter_values"},
			},
		},
	})
}

func testAccDdsParameterTemplate_withEntityId_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "test" {
  name                  = "%[2]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  password              = "Terraform@123"
  mode                  = "ReplicaSet"
  enterprise_project_id = "%[3]s"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 3
    size      = 10
    spec_code = "dds.mongodb.s6.large.2.repset"
  }
}`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDdsParameterTemplate_withEntityId_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template" "test" {
  name        = "%[2]s"
  entity_id   = huaweicloud_dds_instance.test.id
  description = "terraform test"
}
`, testAccDdsParameterTemplate_withEntityId_base(name), name)
}

func testAccDdsParameterTemplate_withEntityId_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template" "test" {
  name        = "%s"
  entity_id   = huaweicloud_dds_instance.test.id
  description = ""

  parameter_values = {
    connPoolMaxConnsPerHost = 800
    cursorTimeoutMillis     = 800000
  }
}
`, testAccDdsParameterTemplate_withEntityId_base(name), name)
}
