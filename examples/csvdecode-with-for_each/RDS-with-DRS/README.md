# RDS-with-DRS

## Description

This example is use csvdecode Function with the for_each meta-argument.

RDS instances of different specifications can be created in batches.

After each RDS DB instance is created, DRS task is automatically created in batches.

## CSV file Argument

* `name` -RDS instance name.

* `flavor` -RDS instance flavor.

* `ha` -synchronization mode. This parameter is left empty for single-node instances.

* `size` -volume size.

* `az` -availability_zone.

* `name1` -DRS task name.

* `sip` -Source DB IP.

* `suser` -Source DB user.

* `spass` -Source DB password.
