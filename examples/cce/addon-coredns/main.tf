data "huaweicloud_cce_clusters" "test" {
  count = var.cluster_id == "" ? 1 : 0

  name = var.cluster_name
}

data "huaweicloud_cce_addon_template" "test" {
  cluster_id = var.cluster_id != "" ? var.cluster_id : try(data.huaweicloud_cce_clusters.test[0].clusters[0].id, null)
  name       = var.addon_template_name
  # ST.002 Disable
  version    = var.addon_version
  # ST.002 Enable
}

resource "huaweicloud_cce_addon" "test" {
  cluster_id    = var.cluster_id
  template_name = var.addon_template_name
  version       = var.addon_version

  values {
    basic_json  = jsonencode(jsondecode(data.huaweicloud_cce_addon_template.test.spec).basic)
    custom_json = jsonencode(jsondecode(data.huaweicloud_cce_addon_template.test.spec).parameters.custom)
    flavor_json = jsonencode(jsondecode(data.huaweicloud_cce_addon_template.test.spec).parameters.flavor1)
  }
}
