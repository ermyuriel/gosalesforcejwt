# gosalesforcejwt

Simple JWT Salesforce authentication in native Go.

## Installation

Generate certificate and key using openssl (or use existing keys and certificate) and set up app in Salesforce. E.g:

```bash
openssl genrsa -out salesforce.key 2048
openssl req -x509 -sha256 -nodes -days 3650  -key salesforce.key -out salesforce.crt
```
Set up environment variables or provide them at runtime:

```bash
SALESFORCE_CLIENT_ID=afdfasgdsagdsdgsdgsdg
SALESFORCE_AUDIENCE=https://test.salesforce.com
SALESFORCE_USER=eucj@eucj.mx
SALESFORCE_INSTANCE=https://login.salesforce.com
SALESFORCE_SANDBOX_INSTANCE=https://test.salesforce.com
SALESFORCE_KEY_PATH=./salesforce.key
```


## Usage



```go

keyBytes, _ := ioutil.ReadFile(os.Getenv("SALESFORCE_KEY_PATH"))

request, _ := BuildRequest(os.Getenv("SALESFORCE_CLIENT_ID"), os.Getenv("SALESFORCE_USER"), os.Getenv("SALESFORCE_AUDIENCE"))

signature, _ := SignRequest(keyBytes, request)

token, _, _ := LogIn(request, signature, os.Getenv("SALESFORCE_SANDBOX_INSTANCE"))

log.Println(token.AccessToken)

```

## Contributing
Just tell me what I did wrong.

## License
[MIT](https://choosealicense.com/licenses/mit/)