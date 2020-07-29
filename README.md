# GO-utils

This project is a set of generic interfaces for GO Lang to interact with AWS using its native SDK. Golang 
has a powerful and robust framework that handles concurrency and parallelism well especially in distributed systems. 

Dependencies:
* [GoLang 1.14](https://golang.org/)
* [AWS GO SDK](https://docs.aws.amazon.com/sdk-for-go/api/aws/)


Intention
-
* The intention of this repo is to have a standardized toolset that can be used for basic/robust interaction with AWS via parameters and configs to abstract out patterns that are commonly used. The methods available are idioms that can be interpreted by other languages such as Python or Java via the same patterns i.e. [the bulder pattern](https://en.wikipedia.org/wiki/Builder_pattern).

Package Breakdown
-  

- [aws](./aws): utility methods that interact with AWS from a data engineers point of view
    1. s3
    2. DynamoDB
    3. Amazon Managed Apache Cassandra Service (MCS)
    4. Cloudwatch
    5. STS (Secured Token Service)
    6. SES (Managed Email Service)
    7. AWS Secrets Manager (Config)

- utilities: helper packages to compliment interactions within the sdk
    1. [database](./database): database interaction methods
    2. [reader](./reader): a consumer of some sort where it takes an input object and reads concurrently via [go channels](https://gobyexample.com/channels) 
    3. [scanner](./scanner): a producer of some sort where it takes an input and writes concurrently via [go routines](https://gobyexample.com/goroutines)
    4. [utils](./utils): utility methods i.e. [Logger](./utils/Logger.go) that compliments repo
    
Tips & Tricks
-    

* To use this as a standard library the import command is
    
        go get -u https://github.com/polyglotDataNerd/poly-GO-utils
        

        