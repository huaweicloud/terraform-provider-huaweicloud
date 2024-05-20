---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_certificates"
description: |-
  Use this data source to get the list of the CSS logstash certificates.
---

# huaweicloud_css_logstash_certificates

Use this data source to get the list of the CSS logstash certificates.

## Example Usage

```hcl
variable "cluster_id" {}
variable "certs_type" {}

data "huaweicloud_css_logstash_certificates" "test" {
  cluster_id = var.cluster_id
  certs_type = var.certs_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies ID of the CSS logstash cluster.

* `certs_type` - (Optional, String) Specifies the certificate type.
  The **defaultCerts** is the default certificate type. If you do not specify the query certificate type,
  it will search the custom certificate list by default.

* `file_name` - (Optional, String) Specifies the file name of the certificate.

* `status` - (Optional, String) Specifies the status of the certificate.
  The values can be **available** and **unavailable**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - The list of the certificates.

  The [certificates](#certificates_struct) structure is documented below.

<a name="certificates_struct"></a>
The `certificates` block supports:

* `id` - The ID of the certificate.

* `file_name` - The name of the certificate.

* `file_location` - The file location of the certificate.

* `status` - The status of the certificate.

* `updated_at` - The upload time of the certificate.
