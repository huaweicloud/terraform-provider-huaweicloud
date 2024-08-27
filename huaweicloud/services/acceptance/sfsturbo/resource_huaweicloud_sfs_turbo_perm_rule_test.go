package sfsturbo

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

func sfsTurboPermRuleReadfunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getPath := sfsClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", sfsClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{share_id}", state.Primary.Attributes["share_id"])
	getPath = strings.ReplaceAll(getPath, "{rule_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := sfsClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SFS Turbo permission rule %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccSFSTruboPermRule_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo_perm_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		sfsTurboPermRuleReadfunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSFSTruboPermRuleBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ip_cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "rw_type", "rw"),
					resource.TestCheckResourceAttr(resourceName, "user_type", "no_root_squash"),
				),
			},
			{
				Config: testSFSTruboPermRuleUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "rw_type", "ro"),
					resource.TestCheckResourceAttr(resourceName, "user_type", "root_squash"),
				),
			},
		},
	})
}

func testSFSTruboPermRuleBasic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_perm_rule" "test" {
  share_id  = huaweicloud_sfs_turbo.test.id
  ip_cidr   = "192.168.0.0/16"
  rw_type   = "rw"
  user_type = "no_root_squash"
}
`, testAccSFSTurbo_basic(rName))
}

func testSFSTruboPermRuleUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sfs_turbo_perm_rule" "test" {
  share_id  = huaweicloud_sfs_turbo.test.id
  ip_cidr   = "192.168.0.0/16"
  rw_type   = "ro"
  user_type = "root_squash"
}
`, testAccSFSTurbo_basic(rName))
}
