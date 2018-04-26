 CleatHitch

A Kubernetes Webhook Token Authenticator service that provides group information to authentication requests. 

## Deploying

Before deploying CleatHitch, modify k8s/ch-configmap.yaml file to point to the appropriate Identity Provider and lookup information to gather group information for user requests.

* NOTE: At the moment CleatHitch only supports:
  * LDAP 

* `kubectl create namespace cleathitch`
* `kubectl -n cleathitch apply -f k8s/ch-configmap.yaml`

Deploy the CleatHitch NodePort Service:
* `kubectl -n cleathitch apply -f k8s/ch-svc.yaml`

The Kuberentes API requires all communication with is secured. To facilitate this a TLS key and cert will need to be created that the API server cert trusts. If a trusted TLS cert already exists that points to every node the cluster via an IP SAN field, feel free to skip the following steps. To create a self-signed certificate that the API server trusts:

* Copy a valid SSL configuration file so it can be modified:
  * `cp /etc/ssl/openssl.cnf /tmp/openssl.cnf`

* Modify /tmp/openssl.cnf to include:
  * In the `[ req ]` section include:
    * `req_extensions = v3_req`
  * In the `[ v3_req ]` section include both:
    * `basicConstraints = CA:TRUE`
    * `subjectAltName=@alternate_names`
  * Add the following as a section:
    ```
[ alternate_names ]
IP.1 = <WORKER NODE 1 IP ADDRESS>
IP.2 = <WORKER NODE 2 IP ADDRESS>
    ```

    * NOTE: Add each worker node's IP address as an IP alternate name
  * Using the SSL config defined above create an SSL signing cert using a CA trusted by the Kubernetes API server:
    * If the Kubernetes cluster was created using `kubeadm` a tursted CA cert can most likely be found at `/etc/kubernetes/pki` on one of the Control Plane nodes.
    * Otherwise a trusted CA could most likely be found at `/etc/ssl/` on one of the Control Plane nodes.
  * Create a SSL Key:
    * `openssl genrsa -out server.key 2048`
  * Create a SSL Signing Request:
    * `openssl req -new -key tls.key -out tls.csr -config /tmp/openssl.cnf`
  * Create the Self Signed Certificate:
    * `openssl x509 -req -days 3650 -in tls.csr -signkey tls.key -out tls.crt -CA /var/lib/kubernetes/ca.pem -CAkey /var/lib/kubernetes/ca-key.pem -CAcreateserial -extensions v3_req -extfile /tmp/openssl.cnf`

Create a directory to store the webhook config file and SSL certificates:
  * `mkdir /etc/kubernetes/webhook`

Copy the SSL certs created above into this new directory:
  * `cp tls* /etc/kubernetes/webhook/`

Copy the Webhook config into the webhook directory:
  * `cp k8s/ch-webhook.yaml /etc/kubernetes/webhook`

Modify the CleatHitch Webhook config file:
  * Modify the `server` line to point to the CleatHitch NodePort service.
  * Modify `certificate-authority` to point to the certificate authority used to create the CleatHitch SSL cert.


Create a secret containing a TLS certificate and it's key.

`kubectl create secret generic cleathitchsecret --from-file=<PATH TO TLS CERT>/tls.cert --from-file=<PATH TO TLS KEY>/tls.key`

 * NOTE: Cleathitch currently expects the TLS cert and key to be named `tls.crt` and `tls.key`

Deploy the Cleathitch App:

`kubectl -n cleathitch apply -f k8s/ch-deploy.yaml`


## Verifying Access Controls are Working


