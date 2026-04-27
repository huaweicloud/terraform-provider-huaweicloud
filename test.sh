#!/bin/bash
set -o errexit
set -o nounset

package=${1}
pattern=${2}
output=$(mktemp -d)
parallelNum=${parallelNum:-4}

echo "TF_ACC=1 go test ${package} -run ${pattern} -timeout 360m -parallel ${parallelNum}"

TF_LOG=DEBUG TF_ACC=1 go test -covermode=atomic -v -coverprofile=${output}/coverage.out -coverpkg=./huaweicloud/... ${package} -run ${pattern} -timeout 360m -parallel ${parallelNum}

go tool cover -html=${output}/coverage.out -o ./coverageFile.html
rm -rf ${output}
exit 0