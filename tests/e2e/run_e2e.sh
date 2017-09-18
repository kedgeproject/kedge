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


# Run e2e tests
make test-e2e &
TEST_PID=$!

# Watch the pods being generated
kubectl get po --all-namespaces -w &
KUBE_PID=$!

# Kill processes once done
(while kill -0 $TEST_PID; do sleep 5; done) && kill $KUBE_PID
