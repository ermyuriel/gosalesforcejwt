openssl genrsa -out salesforce.key 2048
openssl req -x509 -sha256 -nodes -days 3650  -key salesforce.key -out salesforce.crt