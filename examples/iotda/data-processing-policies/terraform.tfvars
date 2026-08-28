flow_control_policy_name        = "tf_test_iotda_flow_control_policy"
flow_control_policy_description = "Limit-the-data-forwarding-tps-of-the-tenant"
flow_control_policy_limit       = 500
backlog_policy_name             = "tf_test_iotda_backlog_policy"
backlog_policy_description      = "Control-the-size-and-time-of-forwarded-data-backlog"
backlog_policy_size             = "524288000"
backlog_policy_time             = "3600"
iotda_access_address            = "https://779f0f0dd5.st1.iotda-app.cn-north-4.myhuaweicloud.com"
