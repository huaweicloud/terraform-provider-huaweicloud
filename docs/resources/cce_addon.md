---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_addon"
description: ""
---

# huaweicloud_cce_addon

Provides a CCE add-on resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
variable "cluster_id" {}

resource "huaweicloud_cce_addon" "addon_test" {
  cluster_id    = var.cluster_id
  template_name = "metrics-server"
  version       = "1.0.0"
}
```

### Using the Add-on Template as a Starting Point

The parameters that can be passed in `basic_json`, `custom_json` and `flavor_json` depend on the add-on type and
version. The [huaweicloud_cce_addon_template](../data-sources/cce_addon_template.md) data source returns the default
parameters of an add-on version in its `spec` attribute, a JSON string whose `basic` and `parameters` keys map to
these arguments:

* `basic` holds the value for `basic_json`.
* `parameters.custom` holds the value for `custom_json`.
* `parameters.flavor1`, `parameters.flavor2` and so on hold the predefined specifications, one of which can be used as
  the value for `flavor_json`.

The keys inside each object, and the flavors that exist, depend on the add-on and its version, so decode `spec` with
`jsondecode()` (or look at it with `terraform console`) to see them. Use `merge()` to change single values and keep
the defaults for the rest.

```hcl
variable "cluster_id" {}
variable "addon_name" {}
variable "addon_version" {}

data "huaweicloud_cce_addon_template" "test" {
  cluster_id = var.cluster_id
  name       = var.addon_name
  version    = var.addon_version
}

locals {
  addon_spec = jsondecode(data.huaweicloud_cce_addon_template.test.spec)
}

resource "huaweicloud_cce_addon" "test" {
  cluster_id    = var.cluster_id
  template_name = var.addon_name
  version       = var.addon_version

  values {
    basic_json  = jsonencode(local.addon_spec.basic)
    custom_json = jsonencode(merge(local.addon_spec.parameters.custom, {
      # Override single values of the template here.
    }))
    flavor_json = jsonencode(local.addon_spec.parameters.flavor1)
  }
}
```

### More Examples

For a complete example, see [Deploy AutoScaler Addon][autoscaler-best-practice] and the add-on examples in the
[examples/cce][cce-examples] directory of this repository.

[autoscaler-best-practice]: https://hcbp.gitbook.io/huaweicloud-provider/best-practices/cce/addon_autoscaler
[cce-examples]: https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/cce

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE add-on resource.
  If omitted, the provider-level region will be used. Changing this creates a new CCE add-on resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the cluster ID.
  Changing this parameter will create a new resource.

* `template_name` - (Required, String, ForceNew) Specifies the name of the add-on template.
  Changing this parameter will create a new resource.

* `version` - (Optional, String) Specifies the version of the add-on.

* `values` - (Optional, List) Specifies the add-on template installation parameters.
  These parameters vary depending on the add-on. The [structure](#cce_addon_values) is documented below.

<a name="cce_addon_values"></a>
The `values` block supports:

* `basic_json` - (Optional, String) Specifies the json string vary depending on the add-on.

* `custom_json` - (Optional, String) Specifies the json string vary depending on the add-on.

* `flavor_json` - (Optional, String) Specifies the json string vary depending on the add-on.

* `basic` - (Optional, Map) Specifies the key/value pairs vary depending on the add-on.
  Only supports non-nested structure and only supports string type elements.
  This is an alternative to `basic_json`, but it is not recommended.

* `custom` - (Optional, Map) Specifies the key/value pairs vary depending on the add-on.
  Only supports non-nested structure and only supports string type elements.
  This is an alternative to `custom_json`, but it is not recommended.

* `flavor` - (Optional, Map) Specifies the key/value pairs vary depending on the add-on.
  Only supports non-nested structure and only supports string type elements.
  This is an alternative to `flavor_json`, but it is not recommended.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the add-on instance.
* `status` - Add-on status information.
* `description` - Description of add-on instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

CCE add-on can be imported using the cluster ID and add-on ID separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_cce_addon.my_addon <cluster_id>/<id>
```
