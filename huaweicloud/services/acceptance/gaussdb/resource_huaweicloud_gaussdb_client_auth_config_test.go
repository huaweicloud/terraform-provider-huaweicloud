package gaussdb

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

func getGaussDBClientAuthConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var offset int

	typeRaw := state.Primary.Attributes["type"]
	databaseRaw := state.Primary.Attributes["database"]
	userRaw := state.Primary.Attributes["user"]
	addressRaw := state.Primary.Attributes["address"]

	var hbaConf interface{}

	for {
		getPath := getBasePath + buildGaussDBClientAuthConfigQueryParams(offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}
		hbaConf = utils.PathSearch(
			fmt.Sprintf("hba_confs|[?type=='%s' && database=='%s' && user=='%s' && address=='%s']|[0]",
				typeRaw, databaseRaw, userRaw, addressRaw),
			getRespBody,
			nil,
		)

		if hbaConf != nil {
			break
		}

		hbaConfs := utils.PathSearch("hba_confs", getRespBody, []interface{}{}).([]interface{})

		if len(hbaConfs) < 100 {
			break
		}
		offset += 100
	}

	if hbaConf == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return hbaConf, nil
}

func buildGaussDBClientAuthConfigQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%v", offset)
}

func TestAccGaussDBClientAuthConfig_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_gaussdb_client_auth_config.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBClientAuthConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDBClientAuthConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "host"),
					resource.TestCheckResourceAttr(rName, "database", "all"),
					resource.TestCheckResourceAttr(rName, "user", "root"),
					resource.TestCheckResourceAttr(rName, "address", "10.10.0.0/16"),
					resource.TestCheckResourceAttr(rName, "method", "md5"),
				),
			},

			{
				Config: testGaussDBClientAuthConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "method", "sha256"),
				),
			},

			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGaussDBClientAuthConfigImportStateFunc(rName),
			},
		},
	})
}

func testGaussDBClientAuthConfig_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_client_auth_config" "test" {
  instance_id = "%s"
  type        = "host"
  database    = "all"
  user        = "root"
  address     = "10.10.0.0/16"
  method      = "md5"
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID)
}

func testGaussDBClientAuthConfig_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_client_auth_config" "test" {
  instance_id = "%s"
  type        = "host"
  database    = "all"
  user        = "root"
  address     = "10.10.0.0/16"
  method      = "sha256" 
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID)
}

func testAccGaussDBClientAuthConfigImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["type"] == "" ||
			rs.Primary.Attributes["database"] == "" || rs.Primary.Attributes["user"] == "" ||
			rs.Primary.Attributes["address"] == "" {
			return "", fmt.Errorf("one or more required attributes are missing")
		}
		return fmt.Sprintf("%s:%s:%s:%s:%s",
			rs.Primary.Attributes["instance_id"],
			rs.Primary.Attributes["type"],
			rs.Primary.Attributes["database"],
			rs.Primary.Attributes["user"],
			rs.Primary.Attributes["address"]), nil
	}
}
