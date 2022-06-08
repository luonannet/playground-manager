package util

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

var (
	clusterMap map[string]*rest.Config
)

func InitK8s() {
	k8sconfigBytes := `{
		"kind": "Config",
		"apiVersion": "v1",
		"preferences": {},
		"clusters": [
			{
				"name": "osinfra-test-cluster",
				"cluster": {
					"server": "https://159.138.30.232:5443",
					"insecure-skip-tls-verify": true
				}
			}
		],
		"users": [
			{
				"name": "osinfra-test-cluster-user",
				"user": {
					"client-certificate-data": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURQekNDQWllZ0F3SUJBZ0lJU3Y1QnFUbzFIaXN3RFFZSktvWklodmNOQVFFTEJRQXdLVEVaTUJjR0ExVUUKQ2hNUVEwTkZJRlJsWTJodWIyeHZaMmxsY3pFTU1Bb0dBMVVFQXhNRFEwTkZNQjRYRFRJeU1ETXdNakF6TURreQpOVm9YRFRJM01ETXdNakF6TURreU5Wb3diVEZBTUJVR0ExVUVDaE1PYzNsemRHVnRPbTFoYzNSbGNuTXdKd1lEClZRUUtFeUE1TTJObFpqYzRaVEJpWTJVMFpUY3dZV0prT1RJek9UTm1NMlJqWTJRellqRXBNQ2NHQTFVRUF4TWcKT0RZM1ltWXpOVE0xTVRrM05ESmhZV0kxTURNM1pERmlZVEJoWmpGa05HVXdnZ0VpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQ3YrbjhoazltbjBCL0ZDaXM0QjljTjFieElIVitLTHhqdkJGTWJ6SEJvCk0xNUt2aDFtbVllZlFKbldyTUczRGJRMWlXVFhQdDh0UStzR3Q3RmNUUnlsTzZMLzJjNHJrVTV2U25wOWtrUkoKeEVwT2xRbllnOFJYZ3cwU0QzakZtUjQxczdobzgyQzBjUERSMXZmMHVFZ3M3MWhBVlpMVDcreCt3QzBaNmZKSgpGZjhEcmtHK3lUZTdLa3B3dnV2bGJJamNVeVdSbFpnTkgxdzYxRUh1ZDNRbzk4TUdzcFFYbjg0dCtGNzFRbHYyCmdIY0tTdkt0SEhrd2VOM0d3ZHN6RmFqZlVORmpwbWNZZFJ6SHhldVFSRGoxS2hBMDhDYnF4S3QrWWlFUmZWT3MKYnRmM3RBaXJTeHp6VEdWdWJsN2J1ck5xOHJEbmo4d2pkL3lKWFhtUGloVmhBZ01CQUFHakp6QWxNQTRHQTFVZApEd0VCL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUZCUWNEQWpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DCkFRRUFNV1VFYmpzTHFnOGdTQ0NNUHNIQjQ3TGpCZzRBeHN0OTV2aG9WMk9JVy9yTWpUQkVPS3NCMU9yY290aVkKNkhDK0NnRWJPM3B5c2h5RnlMUTRhdTJ1MmhXQ2Q5YVMzaUp6bTdQMVFtUnR5R0ppK0k1OThVSTMwTjhnUTNEYgo3azVqMXZDTFNybTA2ejhvUGRzQys0KzUyTEhRUy81RVQxWXc0ZzZwUzNxNGhkbGJSdnhpVzd3ZzFKL3lmOE8wCkdGU2pHZlArYjArVXRzT2FGYXZRMy9mcjZkcFJSSS91OUlKYnZyMXZYTUNZYnhWNXVOWUg4ODhjbHYybEF0RDIKUWxReEM4bzZmRU5IN2FuZXVHcVVEYmdJemZnZWJOYU9yaVJBOStWSllERHR2c3VaUTNOMFBaK1FRNzc3K05MbAovS05GY0xoUUV3R01IK1JpOENrb1FCcTA3Zz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
					"client-key-data": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBci9wL0laUFpwOUFmeFFvck9BZlhEZFc4U0IxZmlpOFk3d1JURzh4d2FETmVTcjRkClpwbUhuMENaMXF6QnR3MjBOWWxrMXo3ZkxVUHJCcmV4WEUwY3BUdWkvOW5PSzVGT2IwcDZmWkpFU2NSS1RwVUoKMklQRVY0TU5FZzk0eFprZU5iTzRhUE5ndEhEdzBkYjM5TGhJTE85WVFGV1MwKy9zZnNBdEdlbnlTUlgvQTY1Qgp2c2szdXlwS2NMN3I1V3lJM0ZNbGtaV1lEUjljT3RSQjduZDBLUGZEQnJLVUY1L09MZmhlOVVKYjlvQjNDa3J5CnJSeDVNSGpkeHNIYk14V28zMURSWTZabkdIVWN4OFhya0VRNDlTb1FOUEFtNnNTcmZtSWhFWDFUckc3WDk3UUkKcTBzYzgweGxibTVlMjdxemF2S3c1NC9NSTNmOGlWMTVqNG9WWVFJREFRQUJBb0lCQVFDaUM0eWRvc08rSDR2aApxU2Q3Qk4wbXhEWUlRZVFFSHJKYkJpUnhhS1BwajhPNEY3Q0RGY1VwQkJlazhwZSt0RVBKT0tjKy9Zb002SW0zCk9kZDhJeFhKb1V6TFJBanhYWEJZVXNEUWVLdWhNYnNxK1ZxRG9YSmZFeklwKzAwK3gyeG1Ed1Ewd0VmQVlHK2gKQ3M1dzduQ24wQTMzejlUUWpCRzk5Y0NTS0RjT3pFWjlTdE9tWlJxeXpmY1ZoQlJoRkhjNWNnYkhCVTAwYmFyMApIVGM0YjFXb0xMWkpIOFlVU3NLdFhqYjBHK3BRWFdyTzM2NVNYOEh4VGduclFyN3llelNNczhVRHgybWlsVzBJCkNWN2tHdzRFUjJaWXhKTVMyZ3ZTRGdBOXg5V09EcnMrVVVpcU54OGpKbHRQS1NjM0RvTmY4MS8zdTZpTzQ1K0wKN2lINFNDZnhBb0dCQU9SL3pFRXY0ek1FZXlPbjRhTEc1NG5GL2hEWlpKMUwzZkFKVm1FWGtmYXcrRGEwYjRyaAo5Rm9JTG9vbzlSNG1CalVrZzZmS1lobmxHdDR6Tlk5L3Qrc2loVG95RlBuLzdWUVRBWnEyVDMzTHBwN0tKanQzCm5kdzF6MVRhREFpSUVMWFAvYndHRXlYOU1BQ3V6TWxlNUxxSCs0cHJiMmpBZ2JFQ0NNcXhDcTh2QW9HQkFNVW8KZ05JdHhOL3huTEx4ei9VZkY4VkFWeVp1RG9RYkJic2ZRQWdsMnhpL1RrVTNjWHRBdzI4b2pPbnlsajdkZmt5agpJalFOL1B6K0tXWnFXS2JtV0wxWS9KWnVsWWJHL0NkbVZGUjFtUFV5UjYvS3JsQWtUL2FkajBhS3A0K1JKbU1mClFucDZqN09uMk5VZmN1MmFWaGJUT1dVTTg3NVhzY1JUMTQwWGtlQnZBb0dBRXZCa3hhSnFlbmdNbk0vMWE4Q0EKd0hLZ3RGVFFlaGtudEJXU0Mwb2cxd21rQ1NTN2VnOXdhV1gwTlk5ZmdiZzFHNEtEUncwSFJJNHd3K29Lcm5JdgpsRld3SlRxeGNtYzhESlBtdGlRV1FwNzJtSUs0MklUNW1KNzlLRk5LWEFxckYrYTlhVEEzdGhaZVBEQkEyUS84CmRWbUFCK0VQd3VieDJQSUtPSUtrR0JFQ2dZRUFyU3R0TXE2bXZMaEFGV0NWY3N4em1YcHcwYjFiSEtlTGVoYngKcS9ac0lhbHVleGMrb3E1NHVmclpBbWRqbGlYLzJCcDFZVmxqKzJiV1FERnFXblg2UFoyYWhocnRWT3NUOFJ0TgpPTDN0c25nT1JSWjEwc2lDUDYrV2Q3UXpoc05MLzBZbW91Tmhzci9ia1I1RldQTDFhS2p3UVpTRnhvcktCaEpLClFwQVFQaHNDZ1lCUVpQQmNmVyt3cnBDS1Z3cC9GcnVXYjBQYlJNczRncGZ5VDNVUHZGTWl5UThKQy9aVzJPZnYKNmpQSGs2eVFRSnZNU3dHc2toT1VLNWdhUDF4UTFLbHplMTdKRFlSVDdGQnRyQjZFdWR1R3ZwcUJjeGpGU0dnVQpmYjljaE9lSDE4QjQ1Q3pZVUo1cHFyMDVVNkcxcWs2elQ5WjJmQjYvRWg2Ni8yWUdQV0JaRWc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
				}
			}
		],
		"contexts": [
			{
				"name": "osinfra-test-cluster-external",
				"context": {
					"cluster": "osinfra-test-cluster",
					"user": "osinfra-test-cluster-user"
				}
			}
		],
		"current-context": "osinfra-test-cluster-external"
	}`
	// use the current context in kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(k8sconfigBytes)) // clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientcmd.RESTConfigFromKubeConfig([]byte(k8sconfigBytes))

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message

		namespace := "default"
		pod := "example-xxxxx"
		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}