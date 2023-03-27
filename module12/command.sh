#Note for https gateway

# create tls key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=*.k8scamp.com/O=k8scamp" -addext "subjectAltName = DNS:k8scamp.com"
# create secret in istio namespace
kubectl create secret tls k8scamp-tls --cert=./tls.crt --key=./tls.key -n istio-system 



kubectl create -f httpserver-gwsec.yaml
kubectl create -f httpserver-vssec.yaml

# then test
curl  -kH "Host: k8scamp.com" https://$INGRESS_IP/httpserver
curl  -kH "Host: k8scamp.com" https://$INGRESS_IP/nginx
