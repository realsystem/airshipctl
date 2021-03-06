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

package k8sutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/rest/fake"
	"k8s.io/client-go/tools/clientcmd"
	kubeconfig "k8s.io/client-go/tools/clientcmd/api"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/openapi"
	"k8s.io/kubectl/pkg/validation"
)

// MockKubectlFactory implements Factory interface for testing purposes.
type MockKubectlFactory struct {
	MockToDiscoveryClient     func() (discovery.CachedDiscoveryInterface, error)
	MockDynamicClient         func() (dynamic.Interface, error)
	MockOpenAPISchema         func() (openapi.Resources, error)
	MockValidator             func() (validation.Schema, error)
	MockToRESTMapper          func() (meta.RESTMapper, error)
	MockToRESTConfig          func() (*rest.Config, error)
	MockNewBuilder            func() *resource.Builder
	MockToRawKubeConfigLoader func() clientcmd.ClientConfig
	MockClientForMapping      func() (resource.RESTClient, error)
	KubeConfig                kubeconfig.Config
	genericclioptions.ConfigFlags
	cmdutil.Factory
}

// ToDiscoveryClient implements Factory interface
func (f *MockKubectlFactory) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return f.MockToDiscoveryClient()
}

// DynamicClient implements Factory interface
// Returns a mock dynamic client ready for use
func (f *MockKubectlFactory) DynamicClient() (dynamic.Interface, error) { return f.MockDynamicClient() }

// OpenAPISchema implements Factory interface
// Returns a mock openapi schema definition. Schema definition includes metadata and structural information about
// Kubernetes object definitions
func (f *MockKubectlFactory) OpenAPISchema() (openapi.Resources, error) { return f.MockOpenAPISchema() }

// Validator implements Factory interface
// Returns a mock schema that can validate objects stored on disk
func (f *MockKubectlFactory) Validator(bool) (validation.Schema, error) {
	return f.MockValidator()
}

// ToRESTMapper implements Factory interface
// Returns a mock RESTMapper
// RESTMapper allows clients to map resources to kind, and map kind and version to interfaces for manipulating
// those objects. It is primarily intended for consumers of Kubernetes compatible REST APIs
func (f *MockKubectlFactory) ToRESTMapper() (meta.RESTMapper, error) { return f.MockToRESTMapper() }

// ToRESTConfig implements Factory interface
// Returns a mock Config
// Config holds the common attributes that can be passed to a Kubernetes client on initialization
func (f *MockKubectlFactory) ToRESTConfig() (*rest.Config, error) { return f.MockToRESTConfig() }

// NewBuilder implements Factory interface
// Returns a mock object that assists in loading objects from both disk and the server
func (f *MockKubectlFactory) NewBuilder() *resource.Builder { return f.MockNewBuilder() }

// ToRawKubeConfigLoader implements Factory interface
func (f *MockKubectlFactory) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return f.MockToRawKubeConfigLoader()
}

// ClientForMapping implements Factory interface
// Returns a mock RESTClient for working with the specified RESTMapping or an error
func (f *MockKubectlFactory) ClientForMapping(*meta.RESTMapping) (resource.RESTClient, error) {
	return f.MockClientForMapping()
}

// WithToDiscoveryClientByError returns mock discovery client with its respective error
func (f *MockKubectlFactory) WithToDiscoveryClientByError(d discovery.CachedDiscoveryInterface,
	err error) *MockKubectlFactory {
	f.MockToDiscoveryClient = func() (discovery.CachedDiscoveryInterface, error) { return d, err }
	return f
}

// WithOpenAPISchemaByError returns mock openAPISchema with its respective error
func (f *MockKubectlFactory) WithOpenAPISchemaByError(r openapi.Resources, err error) *MockKubectlFactory {
	f.MockOpenAPISchema = func() (openapi.Resources, error) { return r, err }
	return f
}

// WithDynamicClientByError returns mock dynamic client with its respective error
func (f *MockKubectlFactory) WithDynamicClientByError(d dynamic.Interface, err error) *MockKubectlFactory {
	f.MockDynamicClient = func() (dynamic.Interface, error) { return d, err }
	return f
}

// WithValidatorByError returns mock validator with its respective error
func (f *MockKubectlFactory) WithValidatorByError(v validation.Schema, err error) *MockKubectlFactory {
	f.MockValidator = func() (validation.Schema, error) { return v, err }
	return f
}

// WithToRESTMapperByError returns mock RESTMapper with its respective error
func (f *MockKubectlFactory) WithToRESTMapperByError(r meta.RESTMapper, err error) *MockKubectlFactory {
	f.MockToRESTMapper = func() (meta.RESTMapper, error) { return r, err }
	return f
}

// WithToRESTConfigByError returns mock RESTConfig with its respective error
func (f *MockKubectlFactory) WithToRESTConfigByError(r *rest.Config, err error) *MockKubectlFactory {
	f.MockToRESTConfig = func() (*rest.Config, error) { return r, err }
	return f
}

// WithNewBuilderByError returns mock resource builder with its respective error
func (f *MockKubectlFactory) WithNewBuilderByError(r *resource.Builder) *MockKubectlFactory {
	f.MockNewBuilder = func() *resource.Builder { return r }
	return f
}

// WithToRawKubeConfigLoaderByError returns mock raw kubeconfig loader with its respective error
func (f *MockKubectlFactory) WithToRawKubeConfigLoaderByError(c clientcmd.ClientConfig) *MockKubectlFactory {
	f.MockToRawKubeConfigLoader = func() clientcmd.ClientConfig { return c }
	return f
}

// WithClientForMappingByError returns mock client mapping with its respective error
func (f *MockKubectlFactory) WithClientForMappingByError(r resource.RESTClient, err error) *MockKubectlFactory {
	f.MockClientForMapping = func() (resource.RESTClient, error) { return r, err }
	return f
}

// NewMockKubectlFactory defines the functions of MockKubectlFactory with nil values for testing purpose
func NewMockKubectlFactory() *MockKubectlFactory {
	return &MockKubectlFactory{MockDynamicClient: func() (dynamic.Interface, error) { return nil, nil },
		MockToDiscoveryClient:     func() (discovery.CachedDiscoveryInterface, error) { return nil, nil },
		MockOpenAPISchema:         func() (openapi.Resources, error) { return nil, nil },
		MockValidator:             func() (validation.Schema, error) { return nil, nil },
		MockToRESTMapper:          func() (meta.RESTMapper, error) { return nil, nil },
		MockToRESTConfig:          func() (*rest.Config, error) { return nil, nil },
		MockNewBuilder:            func() *resource.Builder { return nil },
		MockToRawKubeConfigLoader: func() clientcmd.ClientConfig { return nil },
		MockClientForMapping:      func() (resource.RESTClient, error) { return nil, nil },
	}
}

// MockClientConfig implements DirectClientConfig interface
// Returns mock client config for testing
type MockClientConfig struct {
	clientcmd.DirectClientConfig
	MockNamespace func() (string, bool, error)
}

// Namespace returns mock namespace for testing
func (c MockClientConfig) Namespace() (string, bool, error) { return c.MockNamespace() }

// WithNamespace returns mock namespace with its respective error
func (c *MockClientConfig) WithNamespace(s string, b bool, err error) *MockClientConfig {
	c.MockNamespace = func() (string, bool, error) { return s, b, err }
	return c
}

// NewMockClientConfig returns mock client config for testing
func NewMockClientConfig() *MockClientConfig {
	return &MockClientConfig{
		MockNamespace: func() (string, bool, error) { return "test", false, nil },
	}
}

// NewFakeFactoryForRC returns a fake Factory object for testing
// It is used to mock network interactions via a rest.Request
func NewFakeFactoryForRC(t *testing.T, filenameRC string) *cmdtesting.TestFactory {
	c := scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)

	f := cmdtesting.NewTestFactory().WithNamespace("test")

	f.ClientConfigVal = cmdtesting.DefaultClientConfig()

	pathRC := "/namespaces/test/replicationcontrollers/test-rc"
	get := "GET"
	_, rcBytes := readReplicationController(t, filenameRC, c)

	f.UnstructuredClient = &fake.RESTClient{
		GroupVersion:         schema.GroupVersion{Version: "v1"},
		NegotiatedSerializer: resource.UnstructuredPlusDefaultContentConfig().NegotiatedSerializer,
		Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			switch p, m := req.URL.Path, req.Method; {
			case p == pathRC && m == get:
				bodyRC := ioutil.NopCloser(bytes.NewReader(rcBytes))
				return &http.Response{StatusCode: http.StatusOK,
					Header: cmdtesting.DefaultHeader(),
					Body:   bodyRC}, nil
			case p == "/namespaces/test/replicationcontrollers" && m == get:
				bodyRC := ioutil.NopCloser(bytes.NewReader(rcBytes))
				return &http.Response{StatusCode: http.StatusOK,
					Header: cmdtesting.DefaultHeader(),
					Body:   bodyRC}, nil
			case p == "/namespaces/test/replicationcontrollers/no-match" && m == get:
				return &http.Response{StatusCode: http.StatusNotFound,
					Header: cmdtesting.DefaultHeader(),
					Body:   cmdtesting.ObjBody(c, &corev1.Pod{})}, nil
			case p == "/api/v1/namespaces/test" && m == get:
				return &http.Response{StatusCode: http.StatusOK,
					Header: cmdtesting.DefaultHeader(),
					Body:   cmdtesting.ObjBody(c, &corev1.Namespace{})}, nil
			default:
				t.Fatalf("unexpected request: %#v\n%#v", req.URL, req)
				return nil, nil
			}
		}),
	}
	return f
}

// Below functions are taken from Kubectl library.
// https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/apply/apply_test.go
func readReplicationController(t *testing.T, filenameRC string, c runtime.Codec) (string, []byte) {
	t.Helper()
	rcObj := readReplicationControllerFromFile(t, filenameRC, c)
	metaAccessor, err := meta.Accessor(rcObj)
	require.NoError(t, err, "Could not read replcation controller")
	rcBytes, err := runtime.Encode(c, rcObj)
	require.NoError(t, err, "Could not read replcation controller")
	return metaAccessor.GetName(), rcBytes
}

func readReplicationControllerFromFile(t *testing.T,
	filename string, c runtime.Decoder) *corev1.ReplicationController {
	data := readBytesFromFile(t, filename)
	rc := corev1.ReplicationController{}
	require.NoError(t, runtime.DecodeInto(c, data, &rc), "Could not read replcation controller")

	return &rc
}

func readBytesFromFile(t *testing.T, filename string) []byte {
	file, err := os.Open(filename)
	require.NoError(t, err, "Could not read file")
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	require.NoError(t, err, "Could not read file")

	return data
}
