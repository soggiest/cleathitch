CleatHitch



# Deploying

Modify `k8s/ch-configmap.yaml` appropriatly.

`kubectl create namespace cleathitch`
`kubectl -n cleathitch apply -f k8s/ch-configmap.yaml`

Create a secret containing a TLS certificate and it's key. During testing I used the API Server cert.

`kubectl create secret generic cleathitchsecret --from-file=./tls.cert --from-file=./tls.key`

`kubectl -n cleathitch apply -f k8s/ch-deploy.yaml`