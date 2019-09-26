#!/usr/bin/env bash

# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

project_id=$(terraform output project_id)

sa_json=$(terraform output sa_key | base64 --decode)
cred_file="/workspace/test/credentials.json"
echo $sa_json > $cred_file
{
  echo "#!/usr/bin/env bash"
  echo "export TEST_PROJECT='$project_id'"
  echo "export TEST_CREDENTIALS='$cred_file'"
} > ../source.sh

