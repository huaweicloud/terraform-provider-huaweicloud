# Create a CBH HA instance

This example creates a CBH HA instance. The CBH dedicated instance requires a VPC,
Subnet and security group. In this example, we all create them in the simplest
way. You can replace them with resources already created in Huawei Cloud.

Compared with the CBH basic instance, the most notable difference of the CBH HA instance
lies in the necessity of clearly designating the primary and backup availability zones.

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The CBH HA instance configuration

| Attributes    | Value           |
|---------------|-----------------|
| name          | Cbh_HA_demo     |
| flavor_id     | cbh.basic.50    |
| password      | Cbh@Huawei123   |
| charging_mode | prePaid         |
| period_unit   | mouth           |
| period        | 1               |

### FAQ

- **How can we configure the availability zones of the CBH HA instance?**

    If your business requires high performance (e.g., low network latency), deploy primary and standby
    availability zones in the same region. This ensures fast and stable network connections.  
    The expected `master_availability_zone`、`slave_availability_zone` can be obtained in the following way.

    ```hcl
    data "huaweicloud_cbh_instance" "cbh_HA_demo" {
        master_availability_zone = data.huaweicloud_availability_zones.default.names[0]
        slave_availability_zone  = data.huaweicloud_availability_zones.default.names[0]
    }
    ```

    For disaster recovery priorities, place primary and standby zones in different regions.  
    This improves system resilience against failures as follow:

    ```hcl
    data "huaweicloud_cbh_instance" "cbh_HA_demo" {
        master_availability_zone = data.huaweicloud_availability_zones.default.names[0]
        slave_availability_zone  = data.huaweicloud_availability_zones.default.names[1]
    }
    ```

    Note: Cross-region deployment may slow down system performance.  
    Carefully evaluate performance trade-offs and disaster recovery needs during planning.

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

It takes about 30 minutes to create a CBH HA instance with primary-backup mode.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.12.0 |
| huaweicloud | >= 1.40.0 |
