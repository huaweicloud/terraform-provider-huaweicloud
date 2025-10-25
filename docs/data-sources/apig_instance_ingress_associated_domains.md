---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_ingress_associated_domains"
description: |-
  Use this data source to query the domain information bound to the ingress of dedicated instance within HuaweiCloud.
---

# huaweicloud_apig_instance_ingress_associated_domains

Use this data source to query the domain information bound to the ingress of dedicated instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "ingress_port_id" {}

data "huaweicloud_apig_instance_ingress_associated_domains" "test" {
  instance_id     = var.instance_id
  ingress_port_id = var.ingress_port_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the dedicated instance is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the API gateway instance.

* `ingress_port_id` - (Required, String) Specifies the ID of the custom ingress port.

* `domain_name` - (Optional, String) Specifies the domain name bound to the ingress port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domain_infos` - The list of domain information bound to the ingress port.  
  The [domain_infos](#attrblock_domain_infos) structure is documented below.

<a name="attrblock_domain_infos"></a>
The `domain_infos` block supports:

* `group_id` - The ID of the API group bound to the ingress port.

* `group_name` - The name of the API group bound to the ingress port.

* `domain_name` - The domain name bound to the ingress port.
