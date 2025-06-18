package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartspipeline"
)

func getPipelineResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipeline(client, state.Primary.Attributes["project_id"], state.Primary.ID)
}

func TestAccPipeline_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipeline_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "is_publish", "false"),
					resource.TestCheckResourceAttr(rName, "banned", "true"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "sources.#"),
					resource.TestCheckResourceAttrSet(rName, "triggers.#"),
					resource.TestCheckResourceAttrSet(rName, "schedules.#"),
					resource.TestCheckResourceAttrSet(rName, "concurrency_control.0.concurrency_number"),
					resource.TestCheckResourceAttrSet(rName, "variables.#"),
				),
			},
			{
				Config: testPipeline_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "test-update"),
					resource.TestCheckResourceAttr(rName, "is_publish", "false"),
					resource.TestCheckResourceAttr(rName, "banned", "false"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "sources.#"),
					resource.TestCheckResourceAttrSet(rName, "triggers.#"),
					resource.TestCheckResourceAttrSet(rName, "schedules.#"),
					resource.TestCheckResourceAttrSet(rName, "concurrency_control.0.concurrency_number"),
					resource.TestCheckResourceAttrSet(rName, "variables.#"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPipelineImportState(rName),
			},
		},
	})
}

func testPipelineImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["project_id"] == "" {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.ID), nil
	}
}

func testRepository_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_repository" "test" {
  project_id       = huaweicloud_codearts_project.test.id
  name             = "%[1]s"
  description      = "Created by terraform acc test"
  gitignore_id     = "Go"
  enable_readme    = 0
  visibility_level = 20
  license_id       = 2
  import_members   = 0
}
`, name)
}

//nolint:revive
func testPipeline_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_pipeline" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[3]s"
  description = "test"
  is_publish  = false
  banned      = true
  definition  = jsonencode({
    "stages": [
      {
        "name": "Stage_1",
        "identifier": "1749524616050e405ad46-29bf-4944-964b-0def6b6fb65e",
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
        "post": null,
        "jobs": [
          {
            "id": "",
            "identifier_old": null,
            "stage_index": null,
            "type": null,
            "name": "ManualReview",
            "async": null,
            "identifier": "JOB_EyJYf",
            "sequence": 0,
            "condition": "$${{ default() }}",
            "strategy": {
              "select_strategy": "selected"
            },
            "timeout": "",
            "resource": "{\"type\":\"system\",\"arch\":\"x86\"}",
            "steps": [
              {
                "runtime_attribution": "agentless",
                "multi_step_editable": 0,
                "official_task_version": "0.0.5",
                "icon_url": "/api/v1/*/common/get-plugin-icon?object_key=official_devcloud_checkpoint-9cdc4b215e47465ab664fe4eaaa4e148-人工审核.png-4b267c81b33648d982e5c2067e4c77b9&icon_type=plugin",
                "name": "ManualReview",
                "task": "official_devcloud_checkpoint",
                "business_type": "Normal",
                "inputs": [
                  {
                    "key": "audit_source",
                    "value": "members"
                  },
                  {
                    "key": "approvers",
                    "value": "%[4]s"
                  },
                  {
                    "key": "audit_role",
                    "value": ""
                  },
                  {
                    "key": "check_strategy",
                    "value": "all"
                  },
                  {
                    "key": "timeout_strategy",
                    "value": "reject"
                  },
                  {
                    "key": "timeout",
                    "value": 3600
                  },
                  {
                    "key": "comment",
                    "value": ""
                  }
                ],
                "env": [],
                "sequence": 0,
                "identifier": "174952489372242b019d1-6fad-47dd-a1b5-1cc561a1cae5",
                "endpoint_ids": []
              }
            ],
            "unfinished_steps": [],
            "condition_tag": "",
            "exec_type": "AGENTLESS_JOB",
            "depends_on": [],
            "reusable_job_id": null
          }
        ],
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

  variables {
    name        = "test_var"
    type        = "string"
    value       = "test_value"
    description = "test variable"
    is_runtime  = true
  }

  triggers {
    git_url        = huaweicloud_codearts_repository.test.https_url
    git_type       = "codehub"
    is_auto_commit = false
    repo_id        = huaweicloud_codearts_repository.test.id

    events {
      type   = "push"
      enable = true
    }

    events {
      type   = "merge_request"
      enable = true
    }

    events {
      type   = "tag_push"
      enable = true
    }
  }

  schedules {
    type          = "periodic"
    name          = "test"
    enable        = true
    days_of_week  = [2,3]
    time_zone     = "China Standard Time"
    start_time    = "18:51"
    end_time      = "19:27"
    interval_time = "3600"
    interval_unit = "s"
  }

  lifecycle {
    ignore_changes = [
      definition,
    ]
  }
}
`, testProject_basic(name), testRepository_basic(name), name, acceptance.HW_USER_ID)
}

//nolint:revive
func testPipeline_update(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_pipeline" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[3]s-update"
  description = "test-update"
  is_publish  = false
  banned      = false
  definition  = jsonencode({
    "stages": [
      {
        "name": "Stage_1",
        "identifier": "1749524616050e405ad46-29bf-4944-964b-0def6b6fb65e",
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
        "post": null,
        "jobs": [
          {
            "id": "",
            "identifier_old": null,
            "stage_index": null,
            "type": null,
            "name": "ManualReview",
            "async": null,
            "identifier": "JOB_EyJYf",
            "sequence": 0,
            "condition": "$${{ default() }}",
            "strategy": {
              "select_strategy": "selected"
            },
            "timeout": "",
            "resource": "{\"type\":\"system\",\"arch\":\"x86\"}",
            "steps": [
              {
                "runtime_attribution": "agentless",
                "multi_step_editable": 0,
                "official_task_version": "0.0.5",
                "icon_url": "/api/v1/*/common/get-plugin-icon?object_key=official_devcloud_checkpoint-9cdc4b215e47465ab664fe4eaaa4e148-人工审核.png-4b267c81b33648d982e5c2067e4c77b9&icon_type=plugin",
                "name": "ManualReview",
                "task": "official_devcloud_checkpoint",
                "business_type": "Normal",
                "inputs": [
                  {
                    "key": "audit_source",
                    "value": "members"
                  },
                  {
                    "key": "approvers",
                    "value": "%[4]s"
                  },
                  {
                    "key": "audit_role",
                    "value": ""
                  },
                  {
                    "key": "check_strategy",
                    "value": "all"
                  },
                  {
                    "key": "timeout_strategy",
                    "value": "reject"
                  },
                  {
                    "key": "timeout",
                    "value": 3600
                  },
                  {
                    "key": "comment",
                    "value": ""
                  }
                ],
                "env": [],
                "sequence": 0,
                "identifier": "174952489372242b019d1-6fad-47dd-a1b5-1cc561a1cae5",
                "endpoint_ids": []
              }
            ],
            "unfinished_steps": [],
            "condition_tag": "",
            "exec_type": "AGENTLESS_JOB",
            "depends_on": [],
            "reusable_job_id": null
          }
        ],
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

  variables {
    name        = "test_var_update"
    type        = "string"
    value       = "test_value"
    description = "test variable"
    is_runtime  = true
  }

  triggers {
    git_url        = huaweicloud_codearts_repository.test.https_url
    git_type       = "codehub"
    is_auto_commit = false
    repo_id        = huaweicloud_codearts_repository.test.id

    events {
      type   = "push"
      enable = false
    }

    events {
      type   = "merge_request"
      enable = true
    }

    events {
      type   = "tag_push"
      enable = true
    }
  }

  schedules {
    type          = "periodic"
    name          = "test_update"
    enable        = true
    days_of_week  = [2,3]
    time_zone     = "China Standard Time"
    start_time    = "18:51"
    end_time      = "19:27"
    interval_time = "3600"
    interval_unit = "s"
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
`, testProject_basic(name), testRepository_basic(name), name, acceptance.HW_USER_ID)
}
