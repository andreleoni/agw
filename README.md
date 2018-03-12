
# Route AWS APIGateway Requests with ease
 AGW transform AWS lambda event message to standard http.Request which can make it easy to work with existing http routers and chaining libraries. With AWS's native support for golang wrapper, shim like nodejs or python is no longer needed.

In short, the usage is
```go
//your standard http handler
func testhandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("test", "test header")
    w.(*agw.LPResponse).WriteBody(map[string]string{
        "test":    "test body",
        "funcArn": agw.LambdaContext.InvokedFunctionArn, //can access context as global variable
        "event":   string(agw.RawMessage),               //can access RawMessage as global variable
    }, false)
}

func main() {
    //use any exsiting router supporting the standard http.Handler
    //like 	"github.com/gorilla/mux"
    mux := mux.NewRouter()
    mux.HandleFunc("/test1/hello", testhandler)
    //lambda is from official sdk "github.com/aws/aws-lambda-go/lambda"
    lambda.Start(agw.Handler(mux))
}
```

### The Full Picture
To use it in the real project we might need some more setups
 1. AWS APIGateway **must** be configured with lambda **proxy** mode, typically like {/proxy+} here's the [doc](http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-as-simple-proxy-for-lambda.html) from aws
 2. Using any http router such as lightning fast [Bone](https://github.com/go-zoo/bone) or popular and feature rich [Gorilla Mux](https://github.com/gorilla/mux) and even with chaining libraries like [Alice](https://github.com/justinas/alice) to write your middlewares


### Complex example
You can deploy this code to aws lambda and link it to apigateway to see how it works in test console of aws apigateway.
```go
func handler1(w http.ResponseWriter, r *http.Request) {
    p1 := bone.GetValue(r, "var")
    bd := string(r.Context().Value(agw.ContextKeyBody).([]byte))
    w.(*agw.LPResponse).WriteBody(map[string]string{
        "agent": r.UserAgent(),
        "var":   p1,
        "bd":    bd,
    }, false)
}

func main() {
    mux := bone.New()
    cors := alice.New(agw.Logging, agw.EnableCORS, agw.ParseBodyBytes)
    mux.Post("/test1/:var", cors.ThenFunc(handler1))

    lambda.Start(func() agw.GatewayHandler {
        return func(ctx context.Context, event json.RawMessage) (interface{}, error) {
            //might be useful to store ctx and event as global variable here
            agp := agw.NewAPIGateParser(event)
            lctx, _ := lambdacontext.FromContext(ctx)
            //deal with the different ctx and event,such as connecting to different db endpoint
            //or setting up global variables
            log.Printf("init here with method %s, req ctx: %+v", agp.Method(), lctx)
            return agw.Process(agp, mux), nil
        }
    }())
} 
```



### Notes

 - The ResponseWriter.Write([]byte) (int, error) is not going to work as normal http response due to the way how lambda and aws apigateway work
 - You need to type assert  ResponseWriter as (*agw.LPResponse) and use WriteBody(out, false|true) to set the return body 
 - Information about aws's base64 support is [here](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-set-up-simple-proxy.html#api-gateway-simple-proxy-for-lambda-output-format)
```go
 func MyHandler(w http.ResponseWriter, r *http.Request) {
    //your logic ....
    w.(*agw.LPResponse).WriteBody(out, false)//false|true indicates whether the body is encoded with base64 or not
}
```


### License
BSD licensed. See the LICENSE file for details.
