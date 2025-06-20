package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPipelineAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPipelineAction_run(name),
			},
		},
	})
}

func testPipelineAction_run_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_pipeline" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[3]s"
  description = "test"
  is_publish  = false
  definition  = jsonencode({
    "stages": [
      {
      	"name": "stage_1",
      	"identifier": "17501613644926b1155c7-cb8a-4af8-87b7-80c73ddb266a",
      	"run_condition": null,
      	"type": null,
      	"sequence": 0,
      	"parallel": null,
      	"pre": [
      	  {
            "runtime_attribution": null,
            "multi_step_editable": 0,
            "official_task_version": null,
            "icon_url": null,
            "name": null,
            "task": "official_devcloud_autoTrigger",
            "business_type": null,
            "inputs": null,
            "env": null,
            "sequence": 0,
            "identifier": null,
            "endpoint_ids": null
          }
      	],
      	"post": [],
      	"jobs": [
      	  {
            "id": "",
            "identifier_old": null,
            "stage_index": null,
            "type": null,
            "name": "new_job",
            "async": null,
            "identifier": "JOB_ijsHS",
            "sequence": 0,
            "condition": "$${{ default() }}",
            "strategy": null,
            "timeout": "",
            "resource": null,
            "steps": [],
            "stage_id": "1750161364492",
            "pipeline_id": null,
            "unfinished_steps": null,
            "condition_tag": null,
            "exec_type": "AGENTLESS_JOB",
            "depends_on": [],
            "reusable_job_id": null
      	  }
      	],
      	"pipeline_id": null,
      	"depends_on": [],
      	"run_always": false
      }
    ]
  })

  sources {
    type = "code"

    params {
      codehub_id     = huaweicloud_codearts_repository.test.id
      git_type       = "codehub"
      git_url        = huaweicloud_codearts_repository.test.https_url
      ssh_git_url    = huaweicloud_codearts_repository.test.ssh_url
      repo_name      = huaweicloud_codearts_repository.test.name
      default_branch = "master"
    }
  }

  concurrency_control {
    concurrency_number = 1
    exceed_action      = "QUEUE"
    enable             = true
  }

  lifecycle {
    ignore_changes = [
      definition,
    ]
  }
}
`, testProject_basic(name), testRepository_basic(name), name)
}

func testPipelineAction_run(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_action" "run" {
  action      = "run"
  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id

  sources {
    type = "code"

    params {
      codehub_id     = huaweicloud_codearts_repository.test.id
      git_type       = "codehub"
      git_url        = huaweicloud_codearts_repository.test.https_url
      default_branch = "master"

      build_params {
        build_type    = "branch"
        event_type    = "Manual"
        target_branch = "master"
      }
    }
  }

  choose_jobs   = ["JOB_ijsHS"]
  choose_stages = ["17501613644926b1155c7-cb8a-4af8-87b7-80c73ddb266a"]
  description   = "demo"
}
`, testPipelineAction_run_basic(name))
}
