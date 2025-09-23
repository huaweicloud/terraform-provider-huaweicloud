package codeartsbuild

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartsbuild"
)

func getBuildTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("codearts_build", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Build client: %s", err)
	}

	return codeartsbuild.GetBuildTask(client, state.Primary.ID)
}

func TestAccBuildTask_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_build_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBuildTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBuildTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "arch", "x86-64"),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "scms.0.url", "huaweicloud_codearts_repository.test", "ssh_url"),
					resource.TestCheckResourceAttr(rName, "scms.0.scm_type", "codehub"),
					resource.TestCheckResourceAttrSet(rName, "scms.0.web_url"),
					resource.TestCheckResourceAttrPair(rName, "scms.0.repo_id", "huaweicloud_codearts_repository.test", "repository_id"),
					resource.TestCheckResourceAttr(rName, "scms.0.branch", "master"),
				),
			},
			{
				Config: testBuildTask_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "arch", "x86-64"),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "scms.0.url", "huaweicloud_codearts_repository.test", "ssh_url"),
					resource.TestCheckResourceAttr(rName, "scms.0.scm_type", "codehub"),
					resource.TestCheckResourceAttrSet(rName, "scms.0.web_url"),
					resource.TestCheckResourceAttrPair(rName, "scms.0.repo_id", "huaweicloud_codearts_repository.test", "repository_id"),
					resource.TestCheckResourceAttr(rName, "scms.0.branch", "master"),
					resource.TestCheckResourceAttrSet(rName, "parameters.#"),
					resource.TestCheckResourceAttrSet(rName, "triggers.#"),
					resource.TestCheckResourceAttrSet(rName, "steps.#"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"steps.1.properties"},
			},
		},
	})
}

func testProject_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_project" "test" {
  name = "%s"
  type = "scrum"
}
`, name)
}

func testRepository_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_repository" "test" {
  project_id       = huaweicloud_codearts_project.test.id
  name             = "%s"
  description      = "Created by terraform acc test"
  gitignore_id     = "Go"
  enable_readme    = 0
  visibility_level = 20
  license_id       = 2
  import_members   = 0
}
`, name)
}

func testBuildTask_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_build_task" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[3]s"
  arch       = "x86-64"

  scms {
    url      = huaweicloud_codearts_repository.test.ssh_url
    scm_type = "codehub"
    web_url  = huaweicloud_codearts_repository.test.web_url
    repo_id  = huaweicloud_codearts_repository.test.repository_id
    branch   = "master"
  }
}
`, testProject_basic(name), testRepository_basic(name), name)
}

func testBuildTask_update(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_build_task" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[3]s-update"
  arch       = "x86-64"

  scms {
    url      = huaweicloud_codearts_repository.test.ssh_url
    scm_type = "codehub"
    web_url  = huaweicloud_codearts_repository.test.web_url
    repo_id  = huaweicloud_codearts_repository.test.repository_id
    branch   = "master"
  }

  parameters {
    name = "hudson.model.StringParameterDefinition"

    params {
      name  = "name"
      value = "test"
    }
    params {
      name  = "type"
      value = "customizeparam"
    }
    params {
      name  = "defaultValue"
      value = "cs"

      limits {
        disable = "1"
        name    = "cs"
      }
      limits{
        name = "dfds"
      }
    }
    params {
      name  = "staticVar"
      value = "false"
    }
    params {
      name  = "sensitiveVar"
      value = "false"
    }
    params {
      name  = "deletion"
      value = "false"
    }
    params {
      name  = "defaults"
      value = "false"
    }
  }

  triggers {
    name = "hudson.triggers.TimerTrigger"

    parameters {
      name  = "spec"
      value = "53 15 * * 1,2,3,4,5"
    }
    parameters {
      name  = "timeZoneId"
      value = "Asia/Shanghai"
    }
    parameters {
      name  = "timeZone"
      value = "China_Standard_Time"
    }
    parameters {
      name  = "isDST"
      value = "false"
    }
  }

  steps {
    enable    = true
    module_id = "devcloud2018.codeci_action_20035.action"
    name      = "Docker Command"
  }

  steps {
    enable     = true
    module_id  = "devcloud2018.codeci_action_20057.action"
    name       = "Update OBS"
    properties = {
      objectKey          = jsonencode("./")
      backetName         = jsonencode("test")
      uploadDirectory    = jsonencode(true)
      artifactSourcePath = jsonencode("bin/*")
      authorizationUser  = jsonencode({
        "displayName": "current user",
        "value": "build" 
      })
      obsHeaders = jsonencode([
        {
          "headerKey": "1",
          "headerValue": "1"
        }
      ])
    }
  }
}
`, testProject_basic(name), testRepository_basic(name), name)
}
