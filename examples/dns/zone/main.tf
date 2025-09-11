# Create a DNS public zone
resource "huaweicloud_dns_zone" "test" {
  name                  = var.dns_public_zone_name
  email                 = var.dns_public_zone_email
  zone_type             = var.dns_public_zone_type
  description           = var.dns_public_zone_description
  ttl                   = var.dns_public_zone_ttl
  enterprise_project_id = var.dns_public_zone_enterprise_project_id
  status                = var.dns_public_zone_status
  dnssec                = var.dns_public_zone_dnssec

  dynamic "router" {
    for_each = var.dns_public_zone_router
    content {
      router_id     = router.value.router_id
      router_region = router.value.router_region
    }
  }
}
