---
subcategory: "SSL Certificate Manager (SCM)"
---

# huaweicloud\_scm\_certificate\_push

Used to push an SSL certificate to another HUAWEI CLOUD service, such as Elastic Load Balance (ELB),
Web Application Firewall (WAF), and Content Delivery Network (CDN).

## Example Usage

```hcl
# Upload the certificate to HUAWEI CLOUD service
resource "huaweicloud_scm_certificate" "certificate_2" {
  name              ="certificate_2"
  certificate       = "-----BEGIN CERTIFICATE-----***-----END CERTIFICATE-----\n"
  certificate_chain = "-----BEGIN CERTIFICATE-----***-----END CERTIFICATE-----\n"
  private_key       = "-----BEGIN PRIVATE KEY-----***-----END PRIVATE KEY-----\n"
}

# Push the certificate to the WAF service in the "la-south-2" region.
resource "huaweicloud_scm_certificate_push" "push_1" {
  certificate_id = huaweicloud_scm_certificate.certificate_2.id
  target_service = "WAF"
  target_project = "la-south-2"
}

# Push the certificate to the CDN service.
resource "huaweicloud_scm_certificate_push" "push_2" {
  certificate_id = huaweicloud_scm_certificate.certificate_2.id
  target_service = "CND"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB certificate resource.
    If omitted, the provider-level region will be used.
    Changing this will push the certificate again.
* `certificate_id` - (Required, String, ForceNew) The id of certificate which uploaded to HUAWEI CLOUD.
* `target_service` - (Required, String, ForceNew) Service to which the certificate is pushed.
    The value can be:
     * `CDN` - The Content Delivery Network.
     * `WAF` - The Web Application Firewall.
     * `Enhance_ELB` - The Elastic Load Balance.
* `target_project` - (Optional, String, ForceNew) The region where the service you want to push a certificate to.
    The same certificate can be pushed repeatedly to the same WAF or ELB service in the same `target_project`,
    but the CDN service can only be pushed once.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 5 minute.

## Error Codes
See Error Codes and Solution: https://support.huaweicloud.com/intl/en-us/api-scm/PushCertificate.html
