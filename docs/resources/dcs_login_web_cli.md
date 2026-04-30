---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_login_web_cli"
description: |-
  Manages a DCS login web cli resource within HuaweiCloud.
---

# huaweicloud_dcs_login_web_cli

Manages a DCS login web cli resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_login_web_cli" "test" {
  instance_id = var.instance_id
  password    = "Huawei_test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to log in to WebCli.
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `password` - (Optional, String, NonUpdatable) The password of the DCS instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
