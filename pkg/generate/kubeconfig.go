package generate

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

const (
	kubeconfigClusterName = "default-cluster"
	kubeconfigUserName    = "default-user"
	kubeconfigContextName = "default-context"
)

func newCmdKubeconfig() *cobra.Command {
	o := newKubeconfigOptions()

	cmd := &cobra.Command{
		Use: "kubeconfig",
		// Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return fmt.Errorf("Validate: %w", err)
			}

			if err := o.Complete(); err != nil {
				return err
			}

			return o.Run(cmd.Context())
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&o.serviceAccount, "sa", "", "serviceaccounts to generate kubeconfig")
	flags.Var(&o.format, "format", o.format.Usage())
	o.configFlags.AddFlags(flags)
	return cmd
}

type kubeconfigOptions struct {
	configFlags *genericclioptions.ConfigFlags

	// serviceAccount 待生成 kubeconfig 的服务账号名
	serviceAccount string
	// eg. https://127.0.0.1:45847
	host string
	// defalut is os.Stdout
	out    io.Writer
	format format

	kubeClient kubernetes.Interface
}

func newKubeconfigOptions() *kubeconfigOptions {
	return &kubeconfigOptions{
		configFlags: genericclioptions.NewConfigFlags(false),
	}
}

// Complete sets all information required for updating the current context
func (o *kubeconfigOptions) Complete() error {
	o.out = os.Stdout

	var err error
	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("Complete: %w", err)
	}

	o.host = restConfig.Host
	o.kubeClient, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("Complete: %w", err)
	}

	return nil
}

func (o *kubeconfigOptions) Validate() error {
	if o.serviceAccount == "" {
		return fmt.Errorf("serviceaccount is required")
	}

	return nil
}

func (o *kubeconfigOptions) Run(ctx context.Context) error {
	sa, err := o.kubeClient.CoreV1().ServiceAccounts(*o.configFlags.Namespace).Get(ctx, o.serviceAccount, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	secret, err := o.kubeClient.CoreV1().Secrets(sa.GetNamespace()).Get(ctx, sa.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	if _, ok := secret.Data["ca.crt"]; !ok {
		return fmt.Errorf("secret [%s] should have data, which key is ca.crt", secret.Name)
	}
	if _, ok := secret.Data["token"]; !ok {
		return fmt.Errorf("secret [%s] should have data, which key is token", secret.Name)
	}

	kubeconfig := api.Config{
		Clusters: map[string]*api.Cluster{
			kubeconfigClusterName: {
				CertificateAuthorityData: secret.Data["ca.crt"],
				Server:                   o.host,
			},
		},
		Contexts: map[string]*api.Context{
			kubeconfigContextName: {
				Cluster:  kubeconfigClusterName,
				AuthInfo: kubeconfigUserName,
			},
		},
		CurrentContext: kubeconfigContextName,
		AuthInfos: map[string]*api.AuthInfo{
			kubeconfigUserName: {
				Token: string(secret.Data["token"]),
			},
		},
	}

	if err := clientcmd.Validate(kubeconfig); err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	// 输出
	out, err := clientcmd.Write(kubeconfig)
	if err != nil {
		return fmt.Errorf("Run: %w", err)
	}

	if o.format.Is(formatJSON) {
		out, err = yaml.YAMLToJSON(out)
		if err != nil {
			return fmt.Errorf("Run: %w", err)
		}
	}

	fmt.Fprintf(o.out, "%s", out)
	return nil
}
