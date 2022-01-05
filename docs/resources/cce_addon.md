---
subcategory: "Cloud Container Engine (CCE)"
---

# huaweicloud_cce_addon

Provides a CCE add-on resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_cce_addon" "addon_test" {
  cluster_id    = var.cluster_id
  template_name = "metrics-server"
  version       = "1.0.0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the CCE add-on resource. If omitted, the
  provider-level region will be used. Changing this creates a new CCE add-on resource.

* `cluster_id` - (Required, String, ForceNew) ID of the cluster. Changing this parameter will create a new resource.

* `template_name` - (Required, String, ForceNew) Name of the add-on template. Changing this parameter will create a new
  resource.

* `version` - (Required, String, ForceNew) Version of the add-on. Changing this parameter will create a new resource.

* `values` - (Optional, List, ForceNew) Add-on template installation parameters. These parameters vary depending on the
  add-on. Structure is documented below. Changing this parameter will create a new resource.

* The `values` block supports:

* `basic_json` - (Optional, String, ForceNew) A json string vary depending on the add-on.
  Changing this parameter will create a new resource.

* `custom_json` - (Optional, String, ForceNew) A json string vary depending on the add-on.
  Changing this parameter will create a new resource.

* `flavor_json` - (Optional, String, ForceNew) A json string vary depending on the add-on.
  Changing this parameter will create a new resource.

* `basic` - (Optional, Map, ForceNew) Key/Value pairs vary depending on the add-on.
  Only supports non-nested structure and only supports string type elements.
  This is an alternative to `basic_json`, but it is not recommended.
  Changing this parameter will create a new resource.

* `custom` - (Optional, Map, ForceNew) Key/Value pairs vary depending on the add-on.
  Only supports non-nested structure and only supports string type elements.
  This is an alternative to `custom_json`, but it is not recommended.
  Changing this parameter will create a new resource.

* `flavor` - (Optional, Map, ForceNew) Key/Value pairs vary depending on the add-on.
  Only supports non-nested structure and only supports string type elements.
  This is an alternative to `flavor_json`, but it is not recommended.
  Changing this parameter will create a new resource.

Arguments which can be passed to the `basic_json`, `custom_json` and `flavor_json` add-on parameters depends on
the add-on type and version. For more detailed description of add-ons
see [add-ons description](https://github.com/huaweicloud/terraform-provider-huaweicloud/blob/master/examples/cce/basic/cce-addon-templates.md)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the add-on instance.
* `status` - Add-on status information.
* `description` - Description of add-on instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 3 minute.

## Import

CCE add-on can be imported using the cluster ID and add-on ID separated by a slash, e.g.:

```
$ terraform import huaweicloud_cce_addon.my_addon bb6923e4-b16e-11eb-b0cd-0255ac101da1/c7ecb230-b16f-11eb-b3b6-0255ac1015a3
```
