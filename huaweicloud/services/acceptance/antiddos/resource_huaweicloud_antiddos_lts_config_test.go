package antiddos_test

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

func getLtsConfigFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("anti-ddos", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Anti-DDoS client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/antiddos/lts-config"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?enterprise_project_id=%s", state.Primary.Attributes["enterprise_project_id"])
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Anti-DDoS LTS config: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error flattening Anti-DDoS LTS config response: %s", err)
	}

	enabled := utils.PathSearch("enabled", respBody, false).(bool)
	if !enabled {
		return nil, golangsdk.ErrDefault404{}
	}
	return respBody, nil
}

func TestAccAntiddosLtsConfig_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_antiddos_lts_config.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLtsConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLtsConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "lts_group_id", "huaweicloud_lts_group.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "lts_attack_stream_id", "huaweicloud_lts_stream.test1", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testLtsConfig_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "lts_group_id", "huaweicloud_lts_group.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "lts_attack_stream_id", "huaweicloud_lts_stream.test2", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
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

func testLtsConfig_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test1" {
  group_name  = "%[1]s-1"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test1" {
  group_id    = huaweicloud_lts_group.test1.id
  stream_name = "%[1]s-1"
  is_favorite = true
}

resource "huaweicloud_lts_group" "test2" {
  group_name  = "%[1]s-2"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test2" {
  group_id    = huaweicloud_lts_group.test2.id
  stream_name = "%[1]s-2"
  is_favorite = true
}
`, name)
}

func testLtsConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_antiddos_lts_config" "test" {
  lts_group_id          = huaweicloud_lts_group.test1.id
  lts_attack_stream_id  = huaweicloud_lts_stream.test1.id
  enterprise_project_id = "0"
}
`, testLtsConfig_base(name))
}

func testLtsConfig_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_antiddos_lts_config" "test" {
  lts_group_id          = huaweicloud_lts_group.test2.id
  lts_attack_stream_id  = huaweicloud_lts_stream.test2.id
  enterprise_project_id = "0"
}
`, testLtsConfig_base(name))
}
