---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_endpoint_whitelist"
description: ""
---

# huaweicloud_apig_endpoint_whitelist

Manage the endpoint service whitelist records of APIG service within HuaweiCloud.

-> There is not need to add a whitelist for the current account. After starting the endpoint service, the whitelist
   of the current account will automatically be added, and this whitelist is not managed by terraform.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_apig_endpoint_whitelist" "test" {
  instance_id = var.instance_id
  whitelists  = [
    "iam:domain::1cc2018e40394f7c9692f1713e76234d",
    "iam:domain::2cc2018e40394f7c9692f1713e76234d",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the endpoint service is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the endpoint service
  belongs. Changing this will create a new resource.

* `whitelists` - (Required, List) Specifies the whitelist records of the endpoint service.
  This is a list of strings. The value must comply with regular validation.
  It consists of two parts, the first part is the fixed format **iam: domain::**, and the second part is the
  enterprise project ID. e.g. **iam:domain::XXX**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also instance ID).

## Import

Whitelist records can be imported using resource ID (`id`), e.g.

```shell
$ terraform import huaweicloud_apig_endpoint_whitelist.test <id>
```
