/*
Copyright © 2021 Ci4Rail GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eventhub

import (
	"fmt"
	"os"
)

const (
	// EnvEventHubConnectionsString environment variable to be passed which contains the eventhub connection string to be used
	EnvEventHubConnectionsString = "EVENTHUB_CONNECTIONSTRING"
)

//ReadConnectionStringFromEnv tries to look up the environment variable for the IoT Hub connection string
func ReadConnectionStringFromEnv() (string, error) {
	val, ok := os.LookupEnv(EnvEventHubConnectionsString)

	if !ok {
		return "", fmt.Errorf("%s not set", EnvEventHubConnectionsString)
	}
	return val, nil
}
