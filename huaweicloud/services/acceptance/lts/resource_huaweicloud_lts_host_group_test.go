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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getHostGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getHostGroup: Query the LTS HostGroup detail
	var (
		getHostGroupHttpUrl = "v3/{project_id}/lts/host-group-list"
		getHostGroupProduct = "lts"
	)
	getHostGroupClient, err := cfg.NewServiceClient(getHostGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS Client: %s", err)
	}

	getHostGroupPath := getHostGroupClient.Endpoint + getHostGroupHttpUrl
	getHostGroupPath = strings.ReplaceAll(getHostGroupPath, "{project_id}", getHostGroupClient.ProjectID)

	getHostGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getHostGroupOpt.JSONBody = utils.RemoveNil(lts.BuildGetOrDeleteHostGroupBodyParams(state.Primary.ID))
	getHostGroupResp, err := getHostGroupClient.Request("POST", getHostGroupPath, &getHostGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HostGroup: %s", err)
	}

	getHostGroupRespBody, err := utils.FlattenResponse(getHostGroupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HostGroup: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?host_group_id=='%s']|[0]", state.Primary.ID)
	getHostGroupRespBody = utils.PathSearch(jsonPath, getHostGroupRespBody, nil)
	if getHostGroupRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getHostGroupRespBody, nil
}

func TestAccHostGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_host_group.test"
	hostId := strings.Split(acceptance.HW_LTS_HOST_IDS, ",")[0]
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSHostGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testHostGroup_basic(name, hostId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckOutput("is_host_id_different", "false"),
				),
			},
			{
				Config:            testHostGroup_import(name, hostId),
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testHostGroup_basic_update_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key_update", "value"),
					resource.TestCheckOutput("is_host_id_different", "false"),
				),
			},
			{
				Config: testHostGroup_basic_update_2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "type", "linux"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key_update", "value"),
					resource.TestCheckResourceAttr(rName, "host_ids.#", "0"),
				),
			},
		},
	})
}

func testHostGroup_basic(name, hostId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_host_group" "test" {
  name     = "%[1]s"
  type     = "linux"
  host_ids = ["%[2]s"]

  tags = {
    foo = "bar"
    key = "value"
  }
}

output "is_host_id_different" {
  value = length(setsubtract(huaweicloud_lts_host_group.test.host_ids,
    tolist(["%[2]s"]))) != 0
}
`, name, hostId)
}

func testHostGroup_import(name, hostId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_host_group" "test" {
  name     = "%s"
  type     = "linux"
  host_ids = ["%s"]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, hostId)
}

func testHostGroup_basic_update_1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_host_group" "test" {
  name     = "%[1]s-update"
  type     = "linux"
  host_ids = split(",", "%[2]s")

  tags = {
    foo        = "bar_update"
    key_update = "value"
  }
}

output "is_host_id_different" {
  value = length(setsubtract(huaweicloud_lts_host_group.test.host_ids, split(",", "%[2]s"))) != 0
}
`, name, acceptance.HW_LTS_HOST_IDS)
}

func testHostGroup_basic_update_2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_host_group" "test" {
  name = "%s-update"
  type = "linux"

  tags = {
    foo        = "bar_update"
    key_update = "value"
  }
}
`, name)
}
