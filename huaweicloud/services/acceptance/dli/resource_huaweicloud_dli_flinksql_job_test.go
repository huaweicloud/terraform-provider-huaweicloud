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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJobResource_basic(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "job_running"),
					resource.TestCheckResourceAttr(resourceName, "type", "flink_sql_job"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccFlinkJobResource_update(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
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

func testAccFlinkJobResource_base(name, region string) string {
	return fmt.Sprintf(`
variable "sql" {
  type    = string
  default = <<EOF
CREATE SOURCE STREAM car_infos (
  car_id STRING,
  car_owner STRING,
  car_brand STRING,
  car_price INT
)
WITH (
  type = "dis",
  region = "%[1]s",
  channel = "%[2]s_input",
  partition_count = "1",
  encode = "csv",
  field_delimiter = ","
);

CREATE SINK STREAM audi_cheaper_than_30w (
  car_id STRING,
  car_owner STRING,
  car_brand STRING,
  car_price INT
)
WITH (
  type = "dis",
  region = "%[1]s",
  channel = "%[2]s_output",
  partition_key = "car_owner",
  encode = "csv",
  field_delimiter = ","
);

INSERT INTO audi_cheaper_than_30w
SELECT *
FROM car_infos
WHERE car_brand = "audia4" and car_price < 30;


CREATE SINK STREAM car_info_data (
  car_id STRING,
  car_owner STRING,
  car_brand STRING,
  car_price INT
)
WITH (
  type ="dis",
  region = "%[1]s",
  channel = "%[2]s_input",
  partition_key = "car_owner",
  encode = "csv",
  field_delimiter = ","
);

INSERT INTO car_info_data
SELECT "1", "lilei", "bmw320i", 28;
INSERT INTO car_info_data
SELECT "2", "hanmeimei", "audia4", 27;
EOF

}

resource "huaweicloud_dis_stream" "stream_input" {
  stream_name     = "%[2]s_input"
  partition_count = 1
  data_type       = "CSV"
  csv_delimiter   = ","
}

resource "huaweicloud_dis_stream" "stream_output" {
  stream_name     = "%[2]s_output"
  partition_count = 1
  data_type       = "CSV"
  csv_delimiter   = ","
}

resource "huaweicloud_dli_queue" "test" {
  name          = "%[2]s"
  cu_count      = 16
  queue_type    = "general"
  resource_mode = 1
}
`, region, name)
}

func testAccFlinkJobResource_basic(name string, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinksql_job" "test" {
  depends_on = [
    huaweicloud_dis_stream.stream_input,
    huaweicloud_dis_stream.stream_output,
  ]

  name       = "%[2]s"
  type       = "flink_sql_job"
  sql        = var.sql
  run_mode   = "exclusive_cluster"
  queue_name = huaweicloud_dli_queue.test.name

  tags = {
    foo = "bar"
  }
}
`, testAccFlinkJobResource_base(name, region), name)
}

func testAccFlinkJobResource_update(name, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinksql_job" "test" {
  depends_on = [
    huaweicloud_dis_stream.stream_input,
    huaweicloud_dis_stream.stream_output,
  ]

  name       = "%[2]s"
  type       = "flink_sql_job"
  sql        = var.sql
  run_mode   = "exclusive_cluster"
  queue_name = huaweicloud_dli_queue.test.name

  tags = {
    owner = "terraform"
  }
}
`, testAccFlinkJobResource_base(name, region), name)
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
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

func testAccResourceDliFlinkJob_streamGraph_base(name string) string {
	return fmt.Sprintf(`
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

resource "huaweicloud_dli_queue" "test" {
  name          = "%[1]s"
  cu_count      = 16
  queue_type    = "general"
  resource_mode = 1
}
`, name)
}

func testAccResourceDliFlinkJob_streamGraph_step1(name, staticEstimatorConfig string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_flinksql_job" "test" {
  name                    = "%s"
  flink_version           = "1.12"
  type                    = "flink_opensource_sql_job"
  run_mode                = "exclusive_cluster"
  sql                     = local.opensourceSql
  queue_name              = huaweicloud_dli_queue.test.name
  graph_type              = "job_graph"
  static_estimator        = true
  static_estimator_config = jsonencode(%s)
}
`, testAccResourceDliFlinkJob_streamGraph_base(name), name, staticEstimatorConfig)
}

func testAccResourceDliFlinkJob_streamGraph_step2(name, operatorConfig string) string {
	return fmt.Sprintf(`
%s
resource "huaweicloud_dli_flinksql_job" "test" {
  name            = "%s"
  flink_version   = "1.12"
  type            = "flink_opensource_sql_job"
  run_mode        = "exclusive_cluster"
  sql             = local.opensourceSql
  queue_name      = huaweicloud_dli_queue.test.name
  graph_type      = "simple_graph"
  operator_config = jsonencode(%s)
}
`, testAccResourceDliFlinkJob_streamGraph_base(name), name, operatorConfig)
}
