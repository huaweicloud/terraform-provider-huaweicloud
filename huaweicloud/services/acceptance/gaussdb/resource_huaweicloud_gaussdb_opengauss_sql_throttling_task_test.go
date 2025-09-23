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

func getOpenGaussSqlThrottlingTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/limit-task-list"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{project_id}", client.ProjectID)
	listBasePath = strings.ReplaceAll(listBasePath, "{instance_id}", state.Primary.Attributes["instance_id"])

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var sqlThrottlingTask interface{}
	offset := 0
	for {
		listPath := listBasePath + buildPageQueryParam(100, offset)
		requestResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		sqlThrottlingTasks := utils.PathSearch("limit_task_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(sqlThrottlingTasks) == 0 {
			break
		}
		searchExpression := fmt.Sprintf("[?task_id=='%s']|[0]", state.Primary.ID)
		sqlThrottlingTask = utils.PathSearch(searchExpression, sqlThrottlingTasks, nil)
		if sqlThrottlingTask != nil {
			break
		}

		offset += 100
		totalCount := utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if int(totalCount) <= offset {
			break
		}
	}
	if sqlThrottlingTask == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return sqlThrottlingTask, nil
}

func buildPageQueryParam(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func TestAccOpenGaussSqlThrottlingTask_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_sql_throttling_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussSqlThrottlingTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussSqlThrottlingTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "task_scope", "SQL"),
					resource.TestCheckResourceAttr(rName, "limit_type", "SQL_ID"),
					resource.TestCheckResourceAttrPair(rName, "limit_type_value",
						"data.huaweicloud_gaussdb_opengauss_sql_templates.test",
						"node_limit_sql_model_list.1.sql_id"),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "parallel_size", "4"),
					resource.TestCheckResourceAttrPair(rName, "sql_model",
						"data.huaweicloud_gaussdb_opengauss_sql_templates.test",
						"node_limit_sql_model_list.1.sql_model"),
					resource.TestCheckResourceAttrPair(rName, "node_infos.0.sql_id",
						"data.huaweicloud_gaussdb_opengauss_sql_templates.test",
						"node_limit_sql_model_list.1.sql_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
					resource.TestCheckResourceAttrSet(rName, "rule_name"),
				),
			},
			{
				Config: testOpenGaussSqlThrottlingTask_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", updateName),
					resource.TestCheckResourceAttr(rName, "parallel_size", "10"),
					resource.TestCheckResourceAttrSet(rName, "modifier"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_time", "end_time"},
				ImportStateIdFunc:       testOpenGaussSqlThrottlingTaskImportState(rName),
			},
		},
	})
}

func TestAccOpenGaussSqlThrottlingTask_sql_type(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_sql_throttling_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussSqlThrottlingTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussSqlThrottlingTask_sql_type(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "task_scope", "SQL"),
					resource.TestCheckResourceAttr(rName, "limit_type", "SQL_TYPE"),
					resource.TestCheckResourceAttr(rName, "limit_type_value", "update"),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "parallel_size", "4"),
					resource.TestCheckResourceAttr(rName, "key_words", "aaa,bbb,ccc"),
					resource.TestCheckResourceAttrPair(rName, "databases",
						"data.huaweicloud_gaussdb_opengauss_databases.test", "databases.0.name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
					resource.TestCheckResourceAttrSet(rName, "rule_name"),
				),
			},
			{
				Config: testOpenGaussSqlThrottlingTask_sql_type_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", updateName),
					resource.TestCheckResourceAttr(rName, "parallel_size", "10"),
					resource.TestCheckResourceAttr(rName, "key_words", "aaa,fff,ggg,kkk"),
					resource.TestCheckResourceAttrPair(rName, "databases",
						"data.huaweicloud_gaussdb_opengauss_databases.test", "databases.1.name"),
					resource.TestCheckResourceAttrSet(rName, "modifier"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_time", "end_time"},
				ImportStateIdFunc:       testOpenGaussSqlThrottlingTaskImportState(rName),
			},
		},
	})
}

func TestAccOpenGaussSqlThrottlingTask_session(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_sql_throttling_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussSqlThrottlingTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussSqlThrottlingTask_session(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "task_scope", "SESSION"),
					resource.TestCheckResourceAttr(rName, "limit_type", "SESSION_ACTIVE_MAX_COUNT"),
					resource.TestCheckResourceAttr(rName, "limit_type_value", "CPU_OR_MEMORY"),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "parallel_size", "4"),
					resource.TestCheckResourceAttr(rName, "cpu_utilization", "20"),
					resource.TestCheckResourceAttr(rName, "memory_utilization", "40"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
					resource.TestCheckResourceAttrSet(rName, "rule_name"),
				),
			},
			{
				Config: testOpenGaussSqlThrottlingTask_session_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", updateName),
					resource.TestCheckResourceAttr(rName, "parallel_size", "10"),
					resource.TestCheckResourceAttr(rName, "cpu_utilization", "50"),
					resource.TestCheckResourceAttr(rName, "memory_utilization", "80"),
					resource.TestCheckResourceAttrSet(rName, "modifier"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_time", "end_time"},
				ImportStateIdFunc:       testOpenGaussSqlThrottlingTaskImportState(rName),
			},
		},
	})
}

func testOpenGaussSqlThrottlingTask_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_instances" "test" {}

locals {
  instance       = [for v in  data.huaweicloud_gaussdb_opengauss_instances.test.instances : v if v.id == "%[1]s"][0]
  master_node_id = [for v in local.instance.nodes : v if v.role == "master"][0].id
}

data "huaweicloud_gaussdb_opengauss_sql_templates" "test" {
  instance_id = "%[1]s"
  node_id     = local.master_node_id
}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id      = "%[1]s"
  task_scope       = "SQL"
  limit_type       = "SQL_ID"
  limit_type_value = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
  task_name        = "%[2]s"
  parallel_size    = 4
  start_time       = "%[3]s"
  end_time         = "%[4]s"
  sql_model        = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_model

  node_infos {
    node_id = local.master_node_id
    sql_id  = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
  }
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, name, acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME,
		acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}

func testOpenGaussSqlThrottlingTask_update(updateName string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_instances" "test" {}

locals {
  instance       = [for v in  data.huaweicloud_gaussdb_opengauss_instances.test.instances : v if v.id == "%[1]s"][0]
  master_node_id = [for v in local.instance.nodes : v if v.role == "master"][0].id
}

data "huaweicloud_gaussdb_opengauss_sql_templates" "test" {
  instance_id = "%[1]s"
  node_id     = local.master_node_id
}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id      = "%[1]s"
  task_scope       = "SQL"
  limit_type       = "SQL_ID"
  limit_type_value = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
  task_name        = "%[2]s"
  parallel_size    = 10
  start_time       = "%[3]s"
  end_time         = "%[4]s"
  sql_model        = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_model

  node_infos {
    node_id = local.master_node_id
    sql_id  = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
  }
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, updateName, acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME,
		acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}

func testOpenGaussSqlThrottlingTask_sql_type(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_databases" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id      = "%[1]s"
  task_scope       = "SQL"
  limit_type       = "SQL_TYPE"
  limit_type_value = "update"
  task_name        = "%[2]s"
  parallel_size    = 4
  start_time       = "%[3]s"
  end_time         = "%[4]s"
  key_words        = "aaa,bbb,ccc"
  databases        = data.huaweicloud_gaussdb_opengauss_databases.test.databases[0].name
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, name, acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME,
		acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}

func testOpenGaussSqlThrottlingTask_sql_type_update(updateName string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_databases" "test" {
  instance_id = "%[1]s"
}
resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id      = "%[1]s"
  task_scope       = "SQL"
  limit_type       = "SQL_TYPE"
  limit_type_value = "update"
  task_name        = "%[2]s"
  parallel_size    = 10
  start_time       = "%[3]s"
  end_time         = "%[4]s"
  key_words        = "aaa,fff,ggg,kkk"
  databases        = data.huaweicloud_gaussdb_opengauss_databases.test.databases[1].name
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, updateName, acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME,
		acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}

func testOpenGaussSqlThrottlingTask_session(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id        = "%[1]s"
  task_scope         = "SESSION"
  limit_type         = "SESSION_ACTIVE_MAX_COUNT"
  limit_type_value   = "CPU_OR_MEMORY"
  task_name          = "%[2]s"
  parallel_size      = 4
  cpu_utilization    = 20
  memory_utilization = 40
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, name)
}

func testOpenGaussSqlThrottlingTask_session_update(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id        = "%[1]s"
  task_scope         = "SESSION"
  limit_type         = "SESSION_ACTIVE_MAX_COUNT"
  limit_type_value   = "CPU_OR_MEMORY"
  task_name          = "%[2]s"
  parallel_size      = 10
  cpu_utilization    = 50
  memory_utilization = 80
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, updateName)
}

func testOpenGaussSqlThrottlingTaskImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
