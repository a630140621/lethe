# your server name goes here
server=https://172.24.4.140:8899
# the name of the secret containing the service account token goes here
name=artoo-controller-manager-token-tprtb
P="--context tsh1-spock-ak8s -n spock-system"

ca=$(kubectl get secret/$name -o jsonpath='{.data.ca\.crt}' $P)
token=$(kubectl get secret/$name -o jsonpath='{.data.token}' $P | base64 --decode)
namespace=$(kubectl get secret/$name -o jsonpath='{.data.namespace}' $P | base64 --decode)

echo "
apiVersion: v1
kind: Config
clusters:
- name: default-cluster
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: default-context
  context:
    cluster: default-cluster
    namespace: default
    user: default-user
current-context: default-context
users:
- name: default-user
  user:
    token: ${token}
" > sa.kubeconfig
