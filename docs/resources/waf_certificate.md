---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_certificate"
description: |-
  Manages a WAF certificate resource within HuaweiCloud.
---

# huaweicloud_waf_certificate

Manages a WAF certificate resource within HuaweiCloud.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The certificate resource can be used in Cloud Mode and Dedicated Mode.

## Example Usage

```hcl
variable enterprise_project_id {}

resource "huaweicloud_waf_certificate" "test" {
  name                  = "test-name"
  enterprise_project_id = var.enterprise_project_id
  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIFmQl5dh2QUAeo39TIKtadgAgh4zHx09kSgayS9Wph9LEqq7MA+2042L3J9aOa
DAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUR+SosWwALt6PkP0J9iOIxA6RW8gVsLwq
...
+HhDvD/VeOHytX3RAs2GeTOtxyAV5XpKY5r+PkyUqPJj04t3d0Fopi0gNtLpMF=
-----END CERTIFICATE-----
EOT
  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIJwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAM
ATAwMC4GCCsGAQUFBwIBFiJodHRwOi8vY3BzLnJvb3QteDEubGV0c2VuY3J5cHQu
...
he8Y4IWS6wY7bCkjCWDcRQJMEhg76fsO3txE+FiYruq9RUWhiF1myv4Q6W+CyBFC
1qoJFlcDyqSMo5iHq3HLjs
-----END PRIVATE KEY-----
EOT
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the WAF certificate. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the certificate name. The maximum length is `256` characters. Only digits,
  letters, underscores(_), and hyphens(-) are allowed.

* `certificate` - (Required, String) Specifies the certificate content.

* `private_key` - (Required, String) Specifies the private key. This field does not support individual editing.
  Changes to this field will only take effect when `certificate` changes.

-> Only `PEM` format supported for `certificate` and `private_key`, and the newline characters in the file must be
replaced with `\n`.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF certificate.
  For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The certificate ID in UUID format.

* `created_at` - Indicates the time when the certificate uploaded, in RFC3339 format.

* `expired_at` - Indicates the time when the certificate expires, in RFC3339 format.

## Import

There are two ways to import WAF certificate state.

* Using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_certificate.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_certificate.test <id>/<enterprise_project_id>
```

Note that the imported state is not identical to your resource definition, due to security reason. The missing
attributes include `certificate`, and `private_key`. You can ignore changes as below.

```hcl
resource "huaweicloud_waf_certificate" "test" {
    ...
  lifecycle {
    ignore_changes = [
      certificate, private_key
    ]
  }
}
```
