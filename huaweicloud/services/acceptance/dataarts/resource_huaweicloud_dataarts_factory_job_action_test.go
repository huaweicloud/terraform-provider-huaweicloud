package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFactoryJobAction_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_factory_job_action.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFactoryJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFactoryJobAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "job_name", name),
					resource.TestCheckResourceAttr(rName, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(rName, "action", "start"),
					resource.TestCheckResourceAttr(rName, "status", "NORMAL"),
				),
			},
			{
				Config: testFactoryJobAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "job_name", name),
					resource.TestCheckResourceAttr(rName, "process_type", "REAL_TIME"),
					resource.TestCheckResourceAttr(rName, "action", "stop"),
					resource.TestCheckResourceAttr(rName, "status", "STOPPED"),
				),
			},
		},
	})
}

func testFactoryJobAction_realTime_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[1]s"
  workspace_id = "%[2]s"
  process_type = "REAL_TIME"

  nodes {
    name = "SMN_%[1]s"
    type = "SMN"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "topic"
      value = huaweicloud_smn_topic.test.topic_urn
    }

    properties {
      name  = "messageType"
      value = "NORMAL"
    }

    properties {
      name  = "message"
      value = "terraform acceptance test"
    }
  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testFactoryJobAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job_action" "test" {
  depends_on = [
    "huaweicloud_dataarts_factory_job.test"
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test.name
  process_type = huaweicloud_dataarts_factory_job.test.process_type
}
`, testFactoryJobAction_realTime_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testFactoryJobAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job_action" "test" {
  depends_on = [
    "huaweicloud_dataarts_factory_job.test"
  ]

  workspace_id = "%[2]s"
  action       = "stop"
  job_name     = huaweicloud_dataarts_factory_job.test.name
  process_type = huaweicloud_dataarts_factory_job.test.process_type
}
`, testFactoryJobAction_realTime_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func TestAccFactoryJobAction_batchJob(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_factory_job_action.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getFactoryJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsCdmName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testFactoryJobAction_batchPinelineJob_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "job_name", name),
					resource.TestCheckResourceAttr(rName, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(rName, "action", "start"),
					resource.TestCheckResourceAttr(rName, "status", "SCHEDULING"),
				),
			},
			{
				Config: testFactoryJobAction_batchPinelineJob_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "job_name", name),
					resource.TestCheckResourceAttr(rName, "process_type", "BATCH"),
					resource.TestCheckResourceAttr(rName, "action", "stop"),
					resource.TestCheckResourceAttr(rName, "status", "STOPPED"),
				),
			},
		},
	})
}

func testFactoryJobAction_batchPinelineJob_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[1]s"
  workspace_id = "%[2]s"
  process_type = "BATCH"

  nodes {
    name = "Rest_client_%[1]s"
    type = "RESTAPI"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "url"
      value = "https://www.huaweicloud.com/"
    }

    properties {
      name  = "method"
      value = "GET"
    }

    properties {
      name  = "retry"
      value = "false"
    }

    properties {
      name  = "requestMode"
      value = "sync"
    }

    properties {
      name  = "securityAuthentication"
      value = "NONE"
    }

    properties {
      name  = "agentName"
      value = "%[3]s"
    }
  }

  schedule {
    type = "CRON"
    cron {
      expression = "0 0 0 * * ?"
      start_time = "2024-07-24T16:14:04+08"
    }
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CDM_NAME)
}

func testFactoryJobAction_batchPinelineJob_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job_action" "test" {
  depends_on = [
    "huaweicloud_dataarts_factory_job.test"
  ]

  workspace_id = "%[2]s"
  action       = "start"
  job_name     = huaweicloud_dataarts_factory_job.test.name
  process_type = huaweicloud_dataarts_factory_job.test.process_type
}
`, testFactoryJobAction_batchPinelineJob_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testFactoryJobAction_batchPinelineJob_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job_action" "test" {
  depends_on = [
    "huaweicloud_dataarts_factory_job.test"
  ]

  workspace_id = "%[2]s"
  action       = "stop"
  job_name     = huaweicloud_dataarts_factory_job.test.name
  process_type = huaweicloud_dataarts_factory_job.test.process_type
}
`, testFactoryJobAction_batchPinelineJob_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
