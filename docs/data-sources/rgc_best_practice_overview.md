---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_best_practice_overview"
description: |-
  Use this data source to get the best practice overview information in Resource Governance Center.
---

# huaweicloud_rgc_best_practice_overview

Use this data source to get the best practice overview information in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_best_practice_overview" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `detect_time` - The detect time.

* `compliance_audit` - Information about the result of compliance audit.

  The [compliance_audit](#compliance_audit_struct) structure is documented below.

* `data_boundary` - Information about the result of data boundary.

  The [data_boundary](#data_boundary_struct) structure is documented below.

* `financial_governance` - Information about the result of financial governance.

  The [financial_governance](#financial_governance_struct) structure is documented below.

* `identity_permission` - Information about the result of identity permission.

  The [identity_permission](#identity_permission_struct) structure is documented below.

* `network_planning` - Information about the result of network planning.

  The [network_planning](#network_planning_struct) structure is documented below.

* `om_monitor` - Information about the result of O&M.

  The [om_monitor](#om_monitor_struct) structure is documented below.

* `organization_account` - Information about the result of organization.

  The [organization_account](#organization_account_struct) structure is documented below.

* `security_management` - Information about the result of security management.

  The [security_management](#security_management_struct) structure is documented below.

* `total_score` - Total score.

<a name="compliance_audit_struct"></a>
The `compliance_audit` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="data_boundary_struct"></a>
The `data_boundary` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="financial_governance_struct"></a>
The `financial_governance` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="identity_permission_struct"></a>
The `identity_permission` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="network_planning_struct"></a>
The `network_planning` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="om_monitor_struct"></a>
The `om_monitor` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="organization_account_struct"></a>
The `organization_account` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.

<a name="security_management_struct"></a>
The `security_management` block supports:

* `detection_count` - Number of detection items.

* `high_risk_count` - Number of high-risk items.

* `low_risk_count` - Number of low-risk items.

* `medium_risk_count` - Number of medium-risk items.

* `risk_item_description` - Risk description.

* `score` - Detection score.
