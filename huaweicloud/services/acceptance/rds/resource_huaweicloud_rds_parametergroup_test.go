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

func getConfiguration(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccRdsConfiguration_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_parametergroup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getConfiguration,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "description_1"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "datastore.0.version", "8.0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccRdsConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"values"},
			},
		},
	})
}

func testAccRdsConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_parametergroup" "test" {
  name        = "%s"
  description = "description_1"

  values = {
    auto_increment_increment     = "2"
    binlog_rows_query_log_events = "ON"
  }

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}
`, rName)
}

func testAccRdsConfig_update(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_parametergroup" "test" {
  name        = "%s"
  description = ""

  values = {
    bulk_insert_buffer_size = "10"
    connect_timeout         = "10"
  }

  datastore {
    type    = "mysql"
    version = "8.0"
  }
}
`, updateName)
}
