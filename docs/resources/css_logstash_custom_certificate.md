---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_custom_certificate"
description: ""
---

# huaweicloud_css_logstash_custom_certificate

Manages CSS logstash cluster custom certificate resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "bucket_name" {}
variable "cert_object" {}

resource "huaweicloud_css_logstash_custom_certificate" "test"  {
  cluster_id  = var.cluster_id
  bucket_name = var.bucket_name
  cert_object = var.cert_object
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies ID of the CSS logstash cluster.
  Changing this creates a new resource.

* `bucket_name` - (Required, String, ForceNew) Specifies the OBS bucket name where the certificate file is stored.
  Changing this creates a new resource.

* `cert_object` - (Required, String, ForceNew) Specifies the certificate file path to upload in the OBS bucket.
  The certificate name ranges from `4` to `32` digits, must start with a letter and end with
  (.cer|.crt|.rsa|.jks|.pem|.p10|.pfx|.p12|.csr|.der|.keystore), it can contain letters, numbers,
  dashes, underlines or decimal points, but cannot contain other special characters.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `name` - The custom certificate name.

* `path` - The custom certificate path after uploading.

* `status` - The custom certificate status.

* `updated_at` - The custom certificate upload time.

## Import

The CSS logstash cluster custom certificate can be imported using `cluster_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_css_logstash_custom_certificate.test <cluster_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attributes include: `bucket_name`, `cert_object`.
It is generally recommended running `terraform plan` after importing a CSS logstash cluster custom certificate.
You can then decide if changes should be applied to the CSS logstash cluster custom certificate, or the resource
definition should be updated to align with the CSS logstash cluster custom certificate. Also you can ignore changes
as below.

```hcl
resource "huaweicloud_css_logstash_custom_certificate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      bucket_name, cert_object,
    ]
  }
}
```
