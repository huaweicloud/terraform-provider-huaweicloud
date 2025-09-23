package lts

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

func getLtsGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups"
	client, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing the log group: %s", err)
	}

	groupId := state.Primary.ID
	groupResult := utils.PathSearch(fmt.Sprintf("log_groups|[?log_group_id=='%s']|[0]", groupId), respBody, nil)
	if groupResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return groupResult, nil
}

func TestAccGroup_basic(t *testing.T) {
	var (
		group        interface{}
		resourceName = "huaweicloud_lts_group.test"
		rName        = acceptance.RandomAccResourceName()
		rc           = acceptance.InitResourceCheck(resourceName, &group, getLtsGroupResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsGroup_basic(rName, 30),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLtsGroup_basic(rName, 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
				),
			},
			{
				Config: testAccLtsGroup_step3(rName, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.terraform", ""),
				),
			},
			{
				Config: testAccLtsGroup_step4(rName, 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccLtsGroup_basic(name string, ttl int) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = %[2]d

  tags = {
    owner = "terraform"
  }
}
`, name, ttl)
}

func testAccLtsGroup_step3(name string, ttl int) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%s"
  ttl_in_days = %d

  tags = {
    foo       = "bar"
    key       = "value"
    terraform = ""
  }
}
`, name, ttl)
}

func testAccLtsGroup_step4(name string, ttl int) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%s"
  ttl_in_days = %d
}
`, name, ttl)
}

func TestAccGroup_withEpsId(t *testing.T) {
	var (
		group interface{}

		resourceName = "huaweicloud_lts_group.test"
		rName        = acceptance.RandomAccResourceName()
		rc           = acceptance.InitResourceCheck(resourceName, &group, getLtsGroupResourceFunc)
		epsId        = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		updateEpsId  = acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_withEpsId_step1(rName, epsId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", epsId),
				),
			},
			{
				Config: testAccGroup_withEpsId_step2(rName, updateEpsId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", updateEpsId),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroup_withEpsId_step1(name, epsId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name            = "%[1]s"
  ttl_in_days           = 30
  enterprise_project_id = "%[2]s"
}
`, name, epsId)
}

func testAccGroup_withEpsId_step2(name, epsId string) string {
	return testAccGroup_withEpsId_step1(name, epsId)
}
