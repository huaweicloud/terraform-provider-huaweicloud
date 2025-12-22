organization_name           = "tf_test_swr_organization_name"
repository_name             = "tf_test_swr_repository_name"
policy_type                 = "date_rule"
policy_number               = 30
tag_selectors_configuration = [
  {
    kind    = "label"
    pattern = "1.1"
  }
]
