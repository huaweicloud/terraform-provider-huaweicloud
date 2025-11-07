---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_ingress_associated_domains"
description: |-
  Use this data source to get the list of domain information bound to the specified custom ingress port of the APIG
  instance within HuaweiCloud.
---

# huaweicloud_apig_instance_ingress_associated_domains

Use this data source to get the list of domain information bound to the specified custom ingress port of the APIG
instance within HuaweiCloud.

## Example Usage

### Query all domains bound to the specified custom ingress port

```hcl
variable "instance_id" {}
variable "ingress_port_id" {}

data "huaweicloud_apig_instance_ingress_associated_domains" "test" {
  instance_id     = var.instance_id
  ingress_port_id = var.ingress_port_id
}
```

### Query domains bound to the specified custom ingress port by domain name

```hcl
variable "instance_id" {}
variable "ingress_port_id" {}
variable "domain_name" {}

data "huaweicloud_apig_instance_ingress_associated_domains" "test" {
  instance_id     = var.instance_id
  ingress_port_id = var.ingress_port_id
  domain_name     = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the domains are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the ingress port belongs.

* `ingress_port_id` - (Required, String) Specifies the ID of the custom ingress port.

* `domain_name` - (Optional, String) Specifies the domain name that uses the ingress port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `domains` - The list of domain information bound to the ingress port that matched the filter parameters.  
  The [domains](#apig_data_instance_ingress_associated_domains) structure is documented below.

<a name="apig_data_instance_ingress_associated_domains"></a>
The `domains` block supports:

* `name` - The domain name bound to the ingress port.

* `group_id` - The ID of the API group bound to the ingress port.

* `group_name` - The name of the API group bound to the ingress port.
