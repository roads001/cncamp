# install ingress-nginx controller 
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace

# create tls key
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=k8scamp.com/O=k8scamp" -addext "subjectAltName = DNS:k8scamp.com"
# create secret
kubectl create secret tls my-tls --cert=./tls.crt --key=./tls.key

# check the website
curl -H "Host: k8scamp.com" https://10.100.150.160 -k -v
