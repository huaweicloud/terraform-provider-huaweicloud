---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_es_enabled_elb_certificates"
description: |-
  Use this data source to get the list of CSS es cluster enabled certificates.
---

# huaweicloud_css_es_enabled_elb_certificates

Use this data source to get the list of CSS es cluster enabled certificates.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_es_enabled_elb_certificates" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the Elasticsearth cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - The certificates list.

  The [certificates](#certificates_struct) structure is documented below.

<a name="certificates_struct"></a>
The `certificates` block supports:

* `id` - The certificate ID.

* `name` - The certificate name.

* `type` - The type of SL certificate. Divided into server certificate and CA certificate.
