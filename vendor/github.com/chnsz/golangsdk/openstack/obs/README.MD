# README

## Version

3.23.12

## How to change local obs sdk
- Try to submit pull request to the official repository：
[OBS-Repository](https://github.com/huaweicloud/huaweicloud-sdk-go-obs), and waiting for the PR closing.
- If you have to change local obs sdk, you need to record this change in the table below

| commit                                      | pull request                                          | description                 | date       |
|---------------------------------------------|-------------------------------------------------------|-----------------------------|------------|
| dab1291b283ee292638537274b9506de6eec5241    | [GH-11](https://github.com/chnsz/golangsdk/pull/11)   | get proxy URL from env      | 2021.09.10 |
| 893f81801b653bada3a912b293a92bade2adb399    | [GH-109](https://github.com/chnsz/golangsdk/pull/109) | support bucket cross-region | 2022.01.29 |


## How to update local obs sdk from official repository
- Cover your local obs sdk package with the official repository and commit changes.

- Cherry-pick local change commit records, and handling conflicts.

- Update the Version
