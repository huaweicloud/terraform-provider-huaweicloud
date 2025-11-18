---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_recording_rule"
description: |-
  Use this resource to manage Prometheus recording rule resource within HuaweiCloud.
---

# huaweicloud_aom_recording_rule

Use this resource to manage Prometheus recording rule resource within HuaweiCloud.

-> This is resource bind with the Prometheus instance resource. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_aom_recording_rule" "test" {
  instance_id    = var.instance_id
  recording_rule = <<EOF
groups:
  - name: node_basic_aggregation
    interval: 60s
    rules:
      - record: instance:node_memory_usage:percent
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100
EOF
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the Prometheus instance is located.
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Prometheus instance.
  Currently, only Prometheus for CCE and general instances are supported.

* `recording_rule` - (Required, String) Specifies the content of the recording rule, in YAML format.  
  The recording rule must follow the Prometheus recording rule format.  
  It supports the following sub-parameters:
  - `groups` - (Required) Rule groups. A RecordingRule.yaml can configure multiple rule groups.
  - `name` - (Required) The name of the rule group. The rule group name must be unique.
  - `interval` - (Optional) The execution period of the rule group. The default value is `60`s.
  - `rules` - (Required) Rules. A rule group can contain multiple rules.
  - `record` - (Required) The name of the rule. The aggregation rule name must comply with
    [Prometheus metric name specifications](https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels).
  - `expr` - (Required) The calculation expression. Prometheus monitoring will calculate the pre-aggregated metrics
    through this expression. The calculation expression must comply with
    [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/).
  - `labels` - (Optional) The labels of the metric. Labels must comply with
    [Prometheus metric label specifications](https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

AOM recording rules can be imported using the format `<instance_id>/<rule_id>` or `<instance_id>`, e.g.

```bash
$ terraform import huaweicloud_aom_recording_rule.test <instance_id>/<rule_id>
```

or

```bash
$ terraform import huaweicloud_aom_recording_rule.test <instance_id>
```
