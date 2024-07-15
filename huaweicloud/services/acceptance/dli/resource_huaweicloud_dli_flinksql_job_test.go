package dli

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDliFlinkSqlJobResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}
	jobId, _ := strconv.Atoi(state.Primary.ID)
	return flinkjob.Get(client, jobId)
}

func TestAccResourceDliFlinkJob_basic(t *testing.T) {
	var obj flinkjob.CreateSqlJobOpts
	resourceName := "huaweicloud_dli_flinksql_job.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliFlinkSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliGenaralQueueName(t)
			acceptance.TestAccPreCheckDliFlinkVersion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJobResource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "job_running"),
					resource.TestCheckResourceAttr(resourceName, "run_mode", "exclusive_cluster"),
					resource.TestCheckResourceAttr(resourceName, "type", "flink_opensource_sql_job"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", acceptance.HW_DLI_GENERAL_QUEUE_NAME),
					resource.TestCheckResourceAttr(resourceName, "runtime_config.dli_sql_sqlasync_enabled", "true"),
				),
			},
			{
				Config: testAccFlinkJobResource_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "obs_bucket", "huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "smn_topic", name),
					resource.TestCheckResourceAttr(resourceName, "runtime_config.dli_sql_sqlasync_enabled", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccFlinkJobResource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinksql_job" "test" {
  name          = "%[2]s"
  type          = "flink_opensource_sql_job"
  run_mode      = "exclusive_cluster"
  sql           = local.opensourceSql
  queue_name    = "%[3]s"
  flink_version = "%[4]s"
  
  runtime_config = {
    "dli_sql_sqlasync_enabled"= true
  }

  tags = {
    foo = "bar"
  }
}
`, testAccFlinkJobResource_base(), name, acceptance.HW_DLI_GENERAL_QUEUE_NAME, acceptance.HW_DLI_FLINK_VERSION)
}

func testAccFlinkJobResource_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = replace("%[2]s", "_", "-")
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "The display name of topic"
}

resource "huaweicloud_dli_flinksql_job" "test" {
  name          = "%[2]s"
  type          = "flink_opensource_sql_job"
  run_mode      = "exclusive_cluster"
  sql           = local.opensourceSql
  queue_name    = "%[3]s"
  flink_version = "%[4]s"
  obs_bucket    = huaweicloud_obs_bucket.test.bucket
  smn_topic     = huaweicloud_smn_topic.test.name
 
  runtime_config = {
    "dli_sql_sqlasync_enabled"= false
  }

  tags = {
    owner = "terraform"
  }
}
`, testAccFlinkJobResource_base(), name, acceptance.HW_DLI_GENERAL_QUEUE_NAME, acceptance.HW_DLI_FLINK_VERSION)
}

func TestAccResourceDliFlinkJob_streamGraph(t *testing.T) {
	var obj flinkjob.CreateSqlJobOpts
	resourceName := "huaweicloud_dli_flinksql_job.test"
	name := acceptance.RandomAccResourceName()
	// The test found that the IDs in the operatorConfig and staticEstimatorConfig parameters aren't random,
	// and the IDs corresponding to the same SQL statement are consistent.
	staticEstimatorConfig := `{"operator_list":[{"id":"0a448493b4782967b150582570326227"},{"id":"bc764cd8ddf7a0cff126f51c16239658"}]}`
	operatorConfig :=
		`{"operator_list":[{"id":"0a448493b4782967b150582570326227","parallelism":1},{"id":"bc764cd8ddf7a0cff126f51c16239658","parallelism":1}]}`

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliFlinkSqlJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliGenaralQueueName(t)
			// The current flink version does not support generating stream graph.
			acceptance.TestAccPreCheckDliFlinkStreamGraph(t)
			acceptance.TestAccPreCheckDliFlinkVersion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDliFlinkJob_streamGraph_step1(name, staticEstimatorConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "graph_type", "job_graph"),
					resource.TestCheckResourceAttrSet(resourceName, "stream_graph"),
					resource.TestCheckResourceAttr(resourceName, "operator_config", ""),
					resource.TestCheckResourceAttr(resourceName, "static_estimator_config", staticEstimatorConfig),
				),
			},
			{
				Config: testAccResourceDliFlinkJob_streamGraph_step2(name, operatorConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "graph_type", "simple_graph"),
					resource.TestCheckResourceAttrSet(resourceName, "stream_graph"),
					resource.TestCheckResourceAttr(resourceName, "static_estimator_config", ""),
					resource.TestCheckResourceAttr(resourceName, "operator_config", fmt.Sprintf("%v", operatorConfig)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"graph_type", "static_estimator",
				},
			},
		},
	})
}

func testAccFlinkJobResource_base() string {
	return `
locals {
  opensourceSql = <<EOF
create table dataGenSource(
  user_id string,
  amount int
) with (
  'connector' = 'datagen',
  'rows-per-second' = '1', --每秒生成一条数据
  'fields.user_id.kind' = 'random', --为字段user_id指定random生成器
  'fields.user_id.length' = '3' --限制user_id长度为3
);
  
create table printSink(
  user_id string,
  amount int
  ) with (
  'connector' = 'print'
  );

insert into printSink select * from dataGenSource;
EOF
}
`
}

func testAccResourceDliFlinkJob_streamGraph_step1(name, staticEstimatorConfig string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_flinksql_job" "test" {
  name                    = "%s"
  flink_version           = "%s"
  type                    = "flink_opensource_sql_job"
  run_mode                = "exclusive_cluster"
  sql                     = local.opensourceSql
  queue_name              = "%s"
  graph_type              = "job_graph"
  static_estimator        = true
  static_estimator_config = jsonencode(%s)
}
`, testAccFlinkJobResource_base(), name, acceptance.HW_DLI_GENERAL_QUEUE_NAME, acceptance.HW_DLI_FLINK_VERSION, staticEstimatorConfig)
}

func testAccResourceDliFlinkJob_streamGraph_step2(name, operatorConfig string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_dli_flinksql_job" "test" {
  name            = "%s"
  flink_version   = "%s"
  type            = "flink_opensource_sql_job"
  run_mode        = "exclusive_cluster"
  sql             = local.opensourceSql
  queue_name      = "%s"
  graph_type      = "simple_graph"
  operator_config = jsonencode(%s)
}
`, testAccFlinkJobResource_base(), name, acceptance.HW_DLI_GENERAL_QUEUE_NAME, acceptance.HW_DLI_FLINK_VERSION, operatorConfig)
}
