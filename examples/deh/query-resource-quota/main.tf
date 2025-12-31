data "huaweicloud_deh_quotas" "test" {
  resource = var.host_type != "" ? var.host_type : null
}

locals {
  # Query used quotas
  quotas_with_usage = [for v in data.huaweicloud_deh_quotas.test.quota_set : v if v.used > 0]

  # Query available quotas
  quotas_available = [for v in data.huaweicloud_deh_quotas.test.quota_set : v if v.hard_limit > v.used]

  # Query exhausted quotas
  quotas_exhausted = [for v in data.huaweicloud_deh_quotas.test.quota_set : v if v.hard_limit == v.used]
}
