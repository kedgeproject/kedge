#!/bin/bash

# Copyright 2017 The Kedge Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

run_command="go test"

while getopts ":p:t:" opt; do
    case $opt in
	p ) PARALLEL=$OPTARG;;
	t ) TIMEOUT=$OPTARG;;
    esac
done

if [ -n "$PARALLEL" ]; then
    run_command+=" -parallel=$PARALLEL"
fi

if [ -n "$TIMEOUT" ]; then
    run_command+=" -timeout=$TIMEOUT"
fi

if [ -n "$VERBOSE" ]; then
    run_command+=" -v"
fi

if [ -n "$TEST" ]; then
    run_command+=" -run=$TEST"
else
    run_command+=" -run=k8s"
fi

run_command+=" github.com/kedgeproject/kedge/tests/e2e"

# Run e2e tests
eval $run_command &
TEST_PID=$!

echo "======================================"
echo "| Running end-to-end cluster tests.  |"
echo "| Tests will be ran against a k8s    |"
echo "| cluster in separate namespaces     |"
echo "|                                    |"
echo "| Use command:                       |"
if [ "$TEST" == "os" ]; then
    echo -e "| \e[1;34mVERBOSE=yes make test-e2e-os\e[0m       |"
else
    echo -e "| \e[1;34mVERBOSE=yes make test-e2e\e[0m          |"
fi
echo "| for verbosity.                     |"
echo "======================================"
echo ""

# Watch the pods being generated
# Kill processes once done
if [ -n "$VERBOSE" ]; then
  kubectl get po --all-namespaces -w &
  KUBE_PID=$!
  (while ps -p $TEST_PID > /dev/null; do sleep 5; done) && kill $KUBE_PID
else
  (while ps -p $TEST_PID > /dev/null; do sleep 5; done)
fi

# Get the exit status of the test run
wait $TEST_PID
TEST_STATUS=$?

exit $TEST_STATUS
