package geminidb

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

func getGeminiDbAccount(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/db-users?name={name}"
		product = "geminidb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	accountName := fmt.Sprintf("%v", state.Primary.Attributes["name"])
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{name}", accountName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	user := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getRespBody, nil)
	if user == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return user, nil
}

func TestAccGeminiDbAccount_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_geminidb_account.test"
	name := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getGeminiDbAccount,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDbAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "password", "Test@1234567"),
					resource.TestCheckResourceAttr(resourceName, "privilege", "ReadWrite"),
					resource.TestCheckResourceAttr(resourceName, "databases.#", "2"),
				),
			},
			{
				Config: testAccGeminiDbAccount_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "password", "Test@123456789"),
					resource.TestCheckResourceAttr(resourceName, "privilege", "ReadOnly"),
					resource.TestCheckResourceAttr(resourceName, "databases.#", "3"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccGeminiDbAccount_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_account" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  name        = "%[2]s"
  password    = "Test@1234567"
  privilege   = "ReadWrite"
  databases   = [ "1", "2" ]
}
`, testAccGeminiDbInstance_basic(name), name)
}

func testAccGeminiDbAccount_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_account" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  name        = "%[2]s"
  password    = "Test@123456789"
  privilege   = "ReadOnly"
  databases   = [ "1", "2", "3" ]
}
`, testAccGeminiDbInstance_basic(name), name)
}
