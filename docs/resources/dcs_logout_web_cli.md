---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_logout_web_cli"
description: |-
  Manages a DCS logout web cli resource within HuaweiCloud.
---

# huaweicloud_dcs_logout_web_cli

Manages a DCS logout web cli resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "client_id" {}

resource "huaweicloud_dcs_logout_web_cli" "test" {
  instance_id = var.instance_id
  client_id   = var.client_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to log in to WebCli.
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `client_id` - (Optional, String, NonUpdatable) Specifies the client ID returned by the WebCli authentication.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
