/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package client

import (
	"fmt"
)

// ErrProviderNotDefined is returned when wrong AuthType is provided
type ErrProviderNotDefined struct {
	ProviderName string
}

func (e ErrProviderNotDefined) Error() string {
	return fmt.Sprintf("provider %s is not defined in Clusterctl document", e.ProviderName)
}

// ErrProviderRepoNotFound is returned when wrong AuthType is provided
type ErrProviderRepoNotFound struct {
	ProviderName string
	ProviderType string
}

func (e ErrProviderRepoNotFound) Error() string {
	return fmt.Sprintf("failed to find repository for provider %s of type %s", e.ProviderName, e.ProviderType)
}
