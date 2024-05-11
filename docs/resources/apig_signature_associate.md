---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_signature_associate"
description: ""
---

# huaweicloud_apig_signature_associate

Use this resource to bind the APIs to the signature within HuaweiCloud.

-> A signature can only create one `huaweicloud_apig_signature_associate` resource.
   And a published ID for API can only bind a signature.

## Example Usage

```hcl
variable "instance_id" {}
variable "signature_id" {}
variable "api_publish_ids" {
  type = list(string)
}

resource "huaweicloud_apig_signature_associate" "test" {
  instance_id  = var.instance_id
  signature_id = var.signature_id
  publish_ids  = var.api_publish_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the signature and the APIs are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the APIs and the
  signature belong.  
  Changing this will create a new resource.

* `signature_id` - (Required, String, ForceNew) Specifies the signature ID for APIs binding.  
  Changing this will create a new resource.

* `publish_ids` - (Required, List) Specifies the publish IDs corresponding to the APIs bound by the signature.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<instance_id>/<signature_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

Associate resources can be imported using their `signature_id` and the APIG dedicated instance ID to which the signature
belongs, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_signature_associate.test <instance_id>/<signature_id>
```
