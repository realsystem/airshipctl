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

package get

import (
	"fmt"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/k8s/client"
)

// NewGetCommand creates a new command for getting secret information
func NewGetCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	getRootCmd := &cobra.Command{
		Use:   "get",
		Short: "Get secrets",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient(rootSettings)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("client ready")
			res, err1 := c.ClientSet().CoreV1().Secrets("default").List(metav1.ListOptions{})
			if err1 != nil {
				fmt.Println(err1)
			}
			fmt.Println(res)
		},
	}

	return getRootCmd
}
