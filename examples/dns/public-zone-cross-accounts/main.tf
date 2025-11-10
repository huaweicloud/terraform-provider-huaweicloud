data "huaweicloud_dns_zones" "test" {
  provider = huaweicloud.domain_master

  # ST.002 Disable
  name        = var.main_domain_name
  # ST.002 Enable
  zone_type   = "public"
  search_mode = "equal"
}

resource "huaweicloud_dns_zone_authorization" "test" {
  depends_on = [data.huaweicloud_dns_zones.test]

  zone_name = format("%s.%s", var.sub_domain_prefix, try(data.huaweicloud_dns_zones.test.zones[0].name, "master_domain_not_found"))
}

resource "huaweicloud_dns_recordset" "test" {
  provider = huaweicloud.domain_master

  zone_id = try(data.huaweicloud_dns_zones.test.zones[0].id, null)
  name    = format("%s.%s", try(huaweicloud_dns_zone_authorization.test.record[0].host, "host_not_found"), try(data.huaweicloud_dns_zones.test.zones[0].name, "master_domain_not_found"))
  type    = var.recordset_type
  ttl     = var.recordset_ttl
  records = ["\"${huaweicloud_dns_zone_authorization.test.record[0].value}\""]

  provisioner "local-exec" {
    command = "sleep 10"
  }
}

resource "huaweicloud_dns_zone_authorization_verify" "test" {
  depends_on = [huaweicloud_dns_recordset.test]

  authorization_id = huaweicloud_dns_zone_authorization.test.id
}

resource "huaweicloud_dns_zone" "test" {
  depends_on = [huaweicloud_dns_zone_authorization_verify.test]

  name      = format("%s.%s", var.sub_domain_prefix, try(data.huaweicloud_dns_zones.test.zones[0].name, "master_domain_not_found"))
  zone_type = "public"
}
