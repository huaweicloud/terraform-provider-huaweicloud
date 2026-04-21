package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getSecurityDynamicMaskingPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	return dataarts.GetSecurityDynamicMaskingPolicyById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

// HW_DATAARTS_CONNECTION_ID must be the connection ID corresponding to a DWS type data connection.
func TestAccSecurityDynamicMaskingPolicy_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dataarts_security_dynamic_masking_policy.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getSecurityDynamicMaskingPolicyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityDynamicMaskingPolicy_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "datasource_type", "DWS"),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(rName, "cluster_name"),
					resource.TestCheckResourceAttr(rName, "database_name", "gaussdb"),
					resource.TestCheckResourceAttr(rName, "table_name", name),
					resource.TestCheckResourceAttr(rName, "conn_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttrSet(rName, "conn_name"),
					resource.TestCheckResourceAttr(rName, "policy_list.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "policy_list.*", map[string]string{
						"column_name":          "total",
						"column_type":          "int8",
						"algorithm_type":       "DWS_SELF_CONFIG",
						"algorithm_detail_dto": `{"end":2,"int_target":0,"start":1}`,
					}),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "policy_list.*", map[string]string{
						"column_name":    "enabled",
						"column_type":    "bool",
						"algorithm_type": "MASK",
					}),
					resource.TestCheckResourceAttrSet(rName, "users"),
					resource.TestCheckResourceAttr(rName, "user_groups", fmt.Sprintf("%s0,%s1", name, name)),
					resource.TestCheckResourceAttr(rName, "schema_name", "public"),
					resource.TestCheckResourceAttrSet(rName, "sync_status"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "create_user"),
				),
			},
			{
				Config: testAccSecurityDynamicMaskingPolicy_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "policy_list.#", "1"),
					resource.TestCheckResourceAttr(rName, "policy_list.0.column_name", "description"),
					resource.TestCheckResourceAttr(rName, "policy_list.0.column_type", "text"),
					resource.TestCheckResourceAttr(rName, "policy_list.0.algorithm_type", "DWS_SELF_CONFIG"),
					resource.TestCheckResourceAttrSet(rName, "policy_list.0.algorithm_detail_dto"),
					resource.TestCheckResourceAttrPair(rName, "users", "data.huaweicloud_identity_users.test", "users.0.name"),
					resource.TestCheckResourceAttrPair(rName, "user_groups", "huaweicloud_identity_group.test.0", "name"),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "update_user"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSecurityDynamicMaskingPolicyImportStateFunc(rName),
			},
		},
	})
}

func testAccSecurityDynamicMaskingPolicyImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		policyId := rs.Primary.ID
		if workspaceId == "" || policyId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, policyId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, policyId), nil
	}
}

func testAccSecurityDynamicMaskingPolicy_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[1]s"
  connection_id = "%[2]s"
}

data "huaweicloud_dws_clusters" "test" {}

resource "huaweicloud_identity_group" "test" {
  count = 2
  name  = "%[3]s${count.index}"
}

data "huaweicloud_identity_users" "test" {}

locals {
  connection_name  = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, null)
  dws_cluster_name = try([for v in data.huaweicloud_dws_clusters.test.clusters : v.name if v.id == "%[4]s"][0], null)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CONNECTION_ID, name, acceptance.HW_DWS_CLUSTER_ID)
}

func testAccSecurityDynamicMaskingPolicy_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_dynamic_masking_policy" "test" {
  workspace_id    = "%[2]s"
  name            = "%[3]s"
  datasource_type = "DWS"
  cluster_id      = "%[4]s"
  cluster_name    = local.dws_cluster_name
  database_name   = "gaussdb"
  table_name      = "%[3]s"
  users           = join(",", slice(data.huaweicloud_identity_users.test.users[*].name, 0, 2))
  user_groups     = join(",", huaweicloud_identity_group.test[*].name)
  conn_id         = "%[5]s"
  conn_name       = local.connection_name
  table_id        = "NativeTable-%[4]s-gaussdb-public-%[3]s"
  schema_name     = "public"

  policy_list {
    column_name          = "total"
    column_type          = "int8"
    algorithm_type       = "DWS_SELF_CONFIG"
    algorithm_detail_dto = jsonencode({
      start      = 1
      end        = 2
      int_target = 0
    })
  }
  policy_list {
    column_name    = "enabled"
    column_type    = "bool"
    algorithm_type = "MASK"
  }
}
`, testAccSecurityDynamicMaskingPolicy_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name,
		acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DATAARTS_CONNECTION_ID)
}

func testAccSecurityDynamicMaskingPolicy_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_dynamic_masking_policy" "test" {
  workspace_id    = "%[2]s"
  name            = "%[3]s"
  datasource_type = "DWS"
  cluster_id      = "%[4]s"
  cluster_name    = local.dws_cluster_name
  database_name   = "gaussdb"
  table_name      = "%[5]s"
  users           = try(data.huaweicloud_identity_users.test.users[0].name, null)
  user_groups     = try(huaweicloud_identity_group.test[0].name, null)
  conn_id         = "%[6]s"
  conn_name       = local.connection_name
  table_id        = "NativeTable-%[4]s-gaussdb-public-%[3]s"
  schema_name     = "public"

  policy_list {
    column_name          = "description"
    column_type          = "text"
    algorithm_type       = "DWS_SELF_CONFIG"
    algorithm_detail_dto = jsonencode({
      start         = 1
      end           = 2
	  string_target = "*"
    })
  }
}
`, testAccSecurityDynamicMaskingPolicy_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, updateName,
		acceptance.HW_DWS_CLUSTER_ID, name, acceptance.HW_DATAARTS_CONNECTION_ID,
	)
}
