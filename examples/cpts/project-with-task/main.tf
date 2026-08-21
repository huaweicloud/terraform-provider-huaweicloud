resource "huaweicloud_cpts_project" "test" {
  name        = var.cpts_project_name
  description = var.cpts_project_description
}

resource "huaweicloud_cpts_task" "test" {
  name                  = var.cpts_task_name
  project_id            = huaweicloud_cpts_project.test.id
  benchmark_concurrency = var.cpts_task_benchmark_concurrency
}
