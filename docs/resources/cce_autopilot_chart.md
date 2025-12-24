---
subcategory: "Cloud Container Engine Autopilot (CCE Autopilot)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_autopilot_chart"
description: |-
  Manages a CCE Autopilot chart resource within HuaweiCloud.
---

# huaweicloud_cce_autopilot_chart

Manages a CCE Autopilot chart resource within HuaweiCloud.

## Example Usage

### Basic Example

```hcl
resource "huaweicloud_cce_autopilot_chart" "test" {
  content    = "./kube-prometheus-stack-55.4.1.tgz"
  parameters = "{\"override\":true,\"skip_lint\":true,\"source\":\"package\"}"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE Autopilot chart resource.
  If omitted, the provider-level region will be used. Changing this creates a new CCE Autopilot chart resource.

* `content` - (Required, String) Specifies the path of the Autopilot chart package to be uploaded.

* `parameters` - (Optional, String) Specifies the parameters of the CCE Autopilot chart.

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

CCE Autopilot chart can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_cce_autopilot_chart.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`content` and `parameters`. It is generally recommended running `terraform plan` after importing an CCE Autopilot chart.
You can then decide if changes should be applied to the chart, or the resource definition should be updated to align
with the chart. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cce_autopilot_chart" "test" {
    ...

  lifecycle {
    ignore_changes = [
      content, parameters,
    ]
  }
}
```
