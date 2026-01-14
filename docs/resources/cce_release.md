---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_release"
description: |-
  Manages a CCE release resource within HuaweiCloud.
---

# huaweicloud_cce_release

Manages a CCE release resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
variable "cluster_id" {}
variable "chart_id" {}
variable "name" {}

resource "huaweicloud_cce_release" "test" {
  cluster_id = var.cluster_id
  chart_id   = var.chart_id
  name       = var.name
  namespace  = "default"
  version    = "4.9.0"

  values_json = jsonencode({
    "key1" : ["value1"]
    "key2" : "value2"
    "key3" : {
      "key1" : "value1",
      "key2" : {
        "sub_key1" : "sub_value1",
        "sub_key2" : "sub_value2"
      }
    }
  })

  description = "created by terraform"

  parameters {
    dry_run = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE release resource.
  If omitted, the provider-level region will be used. Changing this creates a new CCE release resource.

* `cluster_id` - (Required, String) Specifies the CCE cluster ID.

* `chart_id` - (Required, String) Specifies the CCE chart ID.

* `name` - (Required, String) Specifies the release name.

* `namespace` - (Required, String) Specifies the namespace to which a chart release belongs.

* `version` - (Required, String) Specifies the release version.

* `values_json` - (Required, String) Specifies the release value. It's a JSON string.

* `description` - (Optional, String) Specifies the release description.

* `parameters` - (Optional, List) Specifies the release parameters.
  The [values](#cce_release_parameters) structure is documented below.

* `action` - (Optional, String) Specifies the release updating action, only works when updating the release.
  The value can be: **upgrade** or **rollback**.

<a name="cce_release_parameters"></a>
The `parameters` block supports:

* `dry_run` - (Optional, Bool) Specifies whether to dry run. IF set to **true**,
  only chart parameters are verified, and installation is not performed.

* `name_template` - (Optional, String) Specifies the release name template.

* `no_hooks` - (Optional, Bool) Specifies whether to disable hooks during installation.

* `replace` - (Optional, Bool) Specifies whether to replace the release with the same name.

* `recreate` - (Optional, Bool) Specifies whether to rebuild the release.

* `reset_values` - (Optional, Bool) Specifies whether to reset values during an update.

* `release_version` - (Optional, Int) Specifies the version of the rollback release.

* `include_hooks` - (Optional, Bool) Specifies whether to enable hooks during an update or deletion.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The release status.

* `status_description` - The release status description.

* `cluster_name` - The cluster name.

* `chart_name` - The chart name.

* `chart_public` - Whether the chart is public.

* `chart_version` - The chart version.

* `created_at` - The create time.

* `updated_at` - The update time.

## Import

CCE release can be imported using the `cluster_id`, `namespace` and `chart_name`, e.g.:

```bash
$ terraform import huaweicloud_cce_release.test <cluster_id>/<namespace>/<chart_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`version`, `values_json`, `chart_id`, `description`, `parameters` and `action`. It is generally recommended running
`terraform plan` after importing an CCE release. You can then decide if changes should be applied to the release,
or the resource definition should be updated to align with the release. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cce_release" "test" {
    ...

  lifecycle {
    ignore_changes = [
      version, values, chart_id, description, parameters, action,
    ]
  }
}
```
