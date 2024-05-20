---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_feature"
description: |-
  Manages an APIG instance feature resource within HuaweiCloud.
---

# huaweicloud_apig_instance_feature

Manages an APIG instance feature resource within HuaweiCloud.

-> For various types of feature parameter configurations, please refer to the
   [documentation](https://support.huaweicloud.com/intl/en-us/api-apig/apig-api-20200402.html).

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = var.instance_id
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 300
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specified the ID of the dedicated instance to which the feature belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specified the name of the feature.
  Changing this creates a new resource.

* `enabled` - (Optional, Bool) Specified whether to enable the feature. Default value is `false`.

* `config` - (Optional, String) Specified the detailed configuration of the feature.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the feature name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

## Import

The resource can be imported using `instance_id` and `name`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_apig_instance_feature.test <instance_id>/<name>
```
