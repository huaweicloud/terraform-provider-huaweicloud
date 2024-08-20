package aom

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

func buildHeaders(state *terraform.ResourceState) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	epsID := state.Primary.Attributes["enterprise_project_id"]
	if epsID != "" {
		moreHeaders["Enterprise-Project-Id"] = epsID
	}
	return moreHeaders
}

func getMultiAccountAggregationRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getHttpUrl := "v1/{project_id}/aom/aggr-config"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AOM multi account aggregation rule: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening AOM multi account aggregation rule: %s", err)
	}

	jsonPath := fmt.Sprintf("[?dest_prometheus_id=='%s']|[0]", state.Primary.ID)
	rule := utils.PathSearch(jsonPath, getRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func TestAccMultiAccountAggregationRule_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_multi_account_aggregation_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getMultiAccountAggregationRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccountAggregationRuleEnable(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMultiAccountAggregationRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_aom_prom_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "send_to_source_account", "true"),
				),
			},
			{
				Config: testMultiAccountAggregationRule_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_aom_prom_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "send_to_source_account", "false"),
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

// Use endpoint `https://aomperform.cn-north-7.myhuaweicloud.com/` to test if
// `https://aom.cn-north-7.myhuaweicloud.com/` continuously reporting errors.
func testMultiAccountAggregationRuleBase(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "test" {
  prom_name             = "%s"
  prom_type             = "ACROSS_ACCOUNT"
  enterprise_project_id = "0"
}
`, name)
}

func testMultiAccountAggregationRule_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_multi_account_aggregation_rule" "test" {
  instance_id            = huaweicloud_aom_prom_instance.test.id
  send_to_source_account = true

  accounts {
    id   = "%s"
    name = "%s"
  }

  services {
    service = "SYS.ELB"
    metrics = [
        "huaweicloud_sys_elb_m2_act_conn",
    ]
  }
}
`, testMultiAccountAggregationRuleBase(name), acceptance.HW_DOMAIN_ID, acceptance.HW_DOMAIN_NAME)
}

func testMultiAccountAggregationRule_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_aom_multi_account_aggregation_rule" "test" {
  instance_id            = huaweicloud_aom_prom_instance.test.id
  send_to_source_account = false

  accounts {
    id   = "%s"
    name = "%s"
  }

  services {
    service = "SYS.ELB"
    metrics = [
        "huaweicloud_sys_elb_m1_cps",
        "huaweicloud_sys_elb_m2_act_conn",
    ]
  }
}
`, testMultiAccountAggregationRuleBase(name), acceptance.HW_DOMAIN_ID, acceptance.HW_DOMAIN_NAME)
}
