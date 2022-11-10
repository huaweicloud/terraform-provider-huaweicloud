#!/usr/bin/env bash

function usage() {
    echo "Usage: ./scripts/codecheck.sh {package}"
    echo "Example1: ./scripts/codecheck.sh ./huaweicloud/services/vpc"
    echo "Example2: ./scripts/codecheck.sh ./huaweicloud/services/..."
    echo ""
}

# Check parameters
package=$1
if [ "X$package" == "X" ]; then
    echo -e "error: package is missing!\n"
    usage
    exit 1
fi
# trim right "/" if necessary
package=${package%/}
packageDir=${package%...}
service=${package##*/}

# Check working directory
workDir=`pwd`
thisDir=${workDir##*/}
if [ "X$thisDir" != "Xterraform-provider-huaweicloud" ]; then
    echo -e "error: the working directory must be terraform-provider-huaweicloud!\n"
    usage
    exit 1
fi

git status >/dev/null
if [ $? -ne 0 ]; then
    echo -e "error: the working directory is not a git repository!\n"
    exit 2
fi

# Check running environment
echo -e "\n==> Checking for running environment..."
LINT=$(which golangci-lint)
SCC=$(which scc)
MISSPELL=$(which misspell)

if [ "X$LINT" == "X" ] || [ "X$SCC" == "X" ] || [ "X$MISSPELL" == "X" ]; then
    echo "    ==> Checking PATH..."
    GOBIN=$(go env GOPATH)/bin
    added=$(echo $PATH | grep -w $GOBIN)
    if [ "X$added" == "X" ]; then
        echo -e "error: the GOBIN is not in PATH, please add it manually!\n"
        exit 2
    fi
fi

if [ "X$LINT" == "X" ]; then
    echo "    ==> Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

if [ "X$SCC" == "X" ]; then
    echo "    ==> Installing boyter/scc..."
    go install github.com/boyter/scc/v3@latest
fi

if [ "X$MISSPELL" == "X" ]; then
    echo "    ==> Installing misspell..."
    go install github.com/client9/misspell/cmd/misspell@latest
fi

# Apply patch
echo -e "\n==> Applying patch..."
git apply --check ./scripts/0001-deprecate-fmtp-and-logp.patch
if [ $? -ne 0 ]; then
    echo -e "warning: cannot apply patch\n"
else
    git apply ./scripts/0001-deprecate-fmtp-and-logp.patch
    applied=TRUE
fi

# Check Code Complexity
echo -e "\n==> Checking for code complexity..."
scc --by-file -s complexity --no-cocomo -w $packageDir | grep -v "/deprecated/"
if [ $? -ne 0 ]; then
    exit 1
fi

# Check golangci-lint
echo -e "\n==> Checking for golangci-lint..."
golangci-lint run $package

# Nolint Directiving
echo -e "\n==> Nolint Directiving..."
grep -rn "nolint:" $packageDir | grep -v "/deprecated/"
grep -rn "lintignore:" $packageDir | grep -v "/deprecated/"

if [ "X$service" != "X..." ] && [[ $package == ./huaweicloud/services/* ]]; then
    grep -rn "markdownlint" ./docs | grep "/${service}_"

    echo -e "\n==> Checking for misspell in $service..."
    misspell ./docs | grep "/${service}_"
    misspell ./examples | grep -w "${service}"

    # update path to "./huaweicloud/services/acceptance/xxx"
    testpackage=${package/"services"/"services/acceptance"}
    if [ ! -d $testpackage ]; then
        echo -e "error: the acceptance directory is not exist!\n"
        exit 1
    fi
    
    echo -e "\n==> Checking for code complexity in $testpackage..."
    scc --by-file -s complexity --no-cocomo -w $testpackage

    echo -e "\n==> Checking for golangci-lint in $testpackage..."
    golangci-lint run $testpackage

    echo -e "\n==> Nolint Directiving in $testpackage..."
    grep -rn "nolint:" $testpackage
    grep -rn "lintignore:" $testpackage
fi

# cleanup
if [ "X$applied" == "XTRUE" ]; then
    echo -e "\n==> Cleanup patch..."
    git checkout -- huaweicloud/utils/fmtp/errors.go
    git checkout -- huaweicloud/utils/logp/log.go
fi

echo -e "\nCheck Completed!\n"
exit 0
