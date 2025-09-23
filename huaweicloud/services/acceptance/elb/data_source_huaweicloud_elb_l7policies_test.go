package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceL7polices_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_l7policies.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceL7Policies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "l7policies.#"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.id"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.description"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.priority"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.action"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.listener_id"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_pool_id"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.rewrite_url_enabled"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.rewrite_url_config.0.host"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.rewrite_url_config.0.path"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.rewrite_url_config.0.query"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.insert_headers_config.0.configs.0.key"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value_type"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.insert_headers_config.0.configs.0.value"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.remove_headers_config.0.configs.0.key"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.traffic_limit_config.0.qps"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.traffic_limit_config.0.per_source_ip_qps"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_extend_config.0.traffic_limit_config.0.burst"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_sticky_session_config.0.enable"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_pools_sticky_session_config.0.timeout"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("l7policy_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),
					resource.TestCheckOutput("action_filter_is_useful", "true"),
					resource.TestCheckOutput("priority_filter_is_useful", "true"),
					resource.TestCheckOutput("provisioning_status_filter_is_useful", "true"),
					resource.TestCheckOutput("redirect_listener_id_filter_is_useful", "true"),
					resource.TestCheckOutput("redirect_pool_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceL7Policies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_l7policies" "test" {
  depends_on = [huaweicloud_elb_l7policy.test]

  l7policy_id = huaweicloud_elb_l7policy.test.id
}

data "huaweicloud_elb_l7policies" "name_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  name = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.name_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.name_filter.l7policies[*].name :v == "%[2]s"]
  )  
}

locals {
  l7policy_id = huaweicloud_elb_l7policy.test.id
}
data "huaweicloud_elb_l7policies" "l7policy_id_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  l7policy_id = huaweicloud_elb_l7policy.test.id
}
output "l7policy_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.l7policy_id_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.l7policy_id_filter.l7policies[*].id : v == local.l7policy_id]
  )  
}

locals {
  description = huaweicloud_elb_l7policy.test.description
}
data "huaweicloud_elb_l7policies" "description_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  description = huaweicloud_elb_l7policy.test.description
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.description_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.description_filter.l7policies[*].description : v == local.description]
  )  
}

locals {
  listener_id = huaweicloud_elb_l7policy.test.listener_id
}
data "huaweicloud_elb_l7policies" "listener_id_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  listener_id = huaweicloud_elb_l7policy.test.listener_id
}
output "listener_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.listener_id_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.listener_id_filter.l7policies[*].listener_id : v == local.listener_id]
  )  
}

locals {
  redirect_listener_id = huaweicloud_elb_l7policy.test.redirect_listener_id
}
data "huaweicloud_elb_l7policies" "redirect_listener_id_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  redirect_listener_id = huaweicloud_elb_l7policy.test.redirect_listener_id
}
output "redirect_listener_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.redirect_listener_id_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.redirect_listener_id_filter.l7policies[*].redirect_listener_id :
  v == local.redirect_listener_id]
  )  
}

locals {
  action = huaweicloud_elb_l7policy.test.action
}
data "huaweicloud_elb_l7policies" "action_id_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  action = huaweicloud_elb_l7policy.test.action
}
output "action_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.action_id_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.action_id_filter.l7policies[*].action : v == local.action]
  )
}

locals {
  priority = huaweicloud_elb_l7policy.test.priority
}
data "huaweicloud_elb_l7policies" "priority_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  priority = huaweicloud_elb_l7policy.test.priority
}
output "priority_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.priority_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.priority_filter.l7policies[*].priority : v == local.priority]
  )  
}

locals {
  provisioning_status = huaweicloud_elb_l7policy.test.provisioning_status
}
data "huaweicloud_elb_l7policies" "provisioning_status_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  provisioning_status = huaweicloud_elb_l7policy.test.provisioning_status
}
output "provisioning_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.provisioning_status_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.provisioning_status_filter.l7policies[*].provisioning_status :
  v == local.provisioning_status]
  )  
}

locals {
  redirect_pool_id = huaweicloud_elb_l7policy.test.redirect_pool_id
}
data "huaweicloud_elb_l7policies" "redirect_pool_id_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  redirect_pool_id = huaweicloud_elb_l7policy.test.redirect_pool_id
}
output "redirect_pool_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.redirect_pool_id_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.redirect_pool_id_filter.l7policies[*].redirect_pool_id :
  v == local.redirect_pool_id]
  )  
}

`, testAccCheckElbV3L7PolicyConfig_basic(name), name)
}

func TestAccDatasourceL7polices_url(t *testing.T) {
	rName := "data.huaweicloud_elb_l7policies.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceL7Policies_url(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "l7policies.#"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.id"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.description"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.priority"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.action"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.listener_id"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_url_config.0.status_code"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_url_config.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_url_config.0.host"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_url_config.0.port"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_url_config.0.path"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.redirect_url_config.0.query"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_url_config.0.insert_headers_config.0.configs.0.key"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_url_config.0.insert_headers_config.0.configs.0.value"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.redirect_url_config.0.insert_headers_config.0.configs.0.value_type"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.updated_at"),
					resource.TestCheckOutput("url_action_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceL7Policies_url(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_l7policies" "test" {
  depends_on = [huaweicloud_elb_l7policy.test]

  l7policy_id = huaweicloud_elb_l7policy.test.id
}

locals {
  action = huaweicloud_elb_l7policy.test.action
}
data "huaweicloud_elb_l7policies" "url_action_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  action = huaweicloud_elb_l7policy.test.action
}
output "url_action_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.url_action_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.url_action_filter.l7policies[*].action : v == local.action]
  )
}

`, testAccCheckElbV3L7PolicyConfig_url(name), name)
}

func TestAccDatasourceL7polices_fixed_response(t *testing.T) {
	rName := "data.huaweicloud_elb_l7policies.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceL7Policies_fixed_response(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "l7policies.#"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.id"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.description"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.priority"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.listener_id"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.action"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.fixed_response_config.0.status_code"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.fixed_response_config.0.content_type"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.fixed_response_config.0.message_body"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.insert_headers_config.0.configs.0.key"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.insert_headers_config.0.configs.0.value"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.insert_headers_config.0.configs.0.value_type"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.remove_headers_config.0.configs.0.key"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.traffic_limit_config.0.qps"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.traffic_limit_config.0.per_source_ip_qps"),
					resource.TestCheckResourceAttrSet(rName,
						"l7policies.0.fixed_response_config.0.traffic_limit_config.0.burst"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "l7policies.0.updated_at"),
					resource.TestCheckOutput("fixed_response_action_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceL7Policies_fixed_response(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_l7policies" "test" {
  depends_on = [huaweicloud_elb_l7policy.test]

  l7policy_id = huaweicloud_elb_l7policy.test.id
}

locals {
  action = huaweicloud_elb_l7policy.test.action
}
data "huaweicloud_elb_l7policies" "fixed_response_action_filter" {
  depends_on = [huaweicloud_elb_l7policy.test]

  action = huaweicloud_elb_l7policy.test.action
}
output "fixed_response_action_filter_is_useful" {
  value = length(data.huaweicloud_elb_l7policies.fixed_response_action_filter.l7policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_l7policies.fixed_response_action_filter.l7policies[*].action : v == local.action]
  )
}
`, testAccCheckElbV3L7PolicyConfig_fixed_response(name), name)
}
