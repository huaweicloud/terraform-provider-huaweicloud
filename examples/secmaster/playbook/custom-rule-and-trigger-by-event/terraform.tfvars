workspace_name   = "tf_test_workspace"
playbook_name    = "tf_test_playbook"
rule_conditions  = [
  {
    name   = "condition1",
    detail = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    data   = [
      "environment.domain_id",
      "==",
      "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
    ]
  },
  {
    name   = "condition2",
    detail = "cn-xxx-x",
    data   = [
      "environment.region_id",
      "==",
      "cn-xxx-x",
    ]
  }
]
approval_content = "Approved for production use"
