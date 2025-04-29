# Create a postpaid EVS volume

This example creates a postpaid EVS volume based on the example
[examples/evs/postpaid-volume](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/evs/postpaid-volume).

To run, configure your HuaweiCloud provider as described in the
[document](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs).

## The postpaid EVS volume configuration

| Attributes   | Value  |
|--------------|--------|
| Size         | 100    |
| Volume_type  | GPSSD2 |
| IOPS         | 3000   |
| Throughput   | 125    |
| Device type  | SCSI   |
| Multi attach | false  |

## Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

It takes about several minutes to create a postpaid EVS volume.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.12.0 |
| huaweicloud | >= 1.49.0 |
