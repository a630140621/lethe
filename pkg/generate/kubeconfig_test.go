package generate

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGenerate(t *testing.T) {
	t.Run("generate kubeconfig", func(t *testing.T) {
		fakeClient := fake.NewSimpleClientset()
		ca := []byte("test")
		token := []byte("test")
		fakeClient.CoreV1().Secrets("default").Create(context.Background(), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-secret",
			},
			Data: map[string][]byte{
				"ca.crt": ca,
				"token":  token,
			},
		}, metav1.CreateOptions{})
		fakeClient.CoreV1().ServiceAccounts("default").Create(context.Background(), &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-service-account",
			},
			Secrets: []corev1.ObjectReference{
				{Name: "test-secret"},
			},
		}, metav1.CreateOptions{})

		var out bytes.Buffer
		opt := &kubeconfigOptions{
			configFlags: &genericclioptions.ConfigFlags{
				Namespace: stringptr("default"),
			},
			serviceAccount: "test-service-account",
			host:           "https://127.0.0.1:45847",
			out:            &out,
			kubeClient:     fakeClient,
		}

		assert.NoError(t, opt.Run(context.Background()))
		c, err := clientcmd.Load(out.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, c.Clusters[kubeconfigClusterName].CertificateAuthorityData, ca)
		assert.Equal(t, c.AuthInfos[kubeconfigUserName].Token, string(token))
	})
}

func stringptr(raw string) *string {
	return &raw
}
