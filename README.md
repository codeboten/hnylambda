# hnylambda
Simple wrapper for AWS Lambda's using Honeycomb

## Usage
```golang
beeline.Init(beeline.Config{
    WriteKey: os.Getenv("HONEYCOMB_KEY"),
    Dataset:  os.Getenv("HONEYCOMB_DATASET"),
})
lambda.Start(hnylambda.Middleware(Handler))
```

Full example code available in example.go