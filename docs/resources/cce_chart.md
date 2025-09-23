---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_chart"
description: ""
---

# huaweicloud_cce_chart

Manages a CCE chart resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
resource "huaweicloud_cce_chart" "test" {
  content    = "./kube-prometheus-stack-55.4.1.tgz"
  parameters = "{\"override\":true,\"skip_lint\":true,\"source\":\"package\"}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE chart resource.
  If omitted, the provider-level region will be used. Changing this creates a new CCE chart resource.

* `content` - (Required, String) Specifies the path of the chart package to be uploaded.

* `parameters` - (Optional, String) Specifies the parameters of the CCE chart.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The chart name.

* `value` - The value of the chart.

* `translate` - The traslate source of the chart.

* `instruction` - The instruction of the chart.

* `version` - The chart version.

* `description` - The description of the chart.

* `source` - The source of the chart.

* `public` - Whether the chart is public.

* `chart_url` - The chart url.

* `created_at` - The create time.

* `updated_at` - The update time.

## Import

CCE chart can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_cce_chart.test 19413aa0-9fe4-11ee-83b0-0255ac10026b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`content` and `parameters`. It is generally recommended running `terraform plan` after importing an CCE chart.
You can then decide if changes should be applied to the chart, or the resource definition should be updated to align
with the chart. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cce_chart" "test" {
    ...

  lifecycle {
    ignore_changes = [
      content, parameters,
    ]
  }
}
```
