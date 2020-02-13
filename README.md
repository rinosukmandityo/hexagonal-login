# hexagonal-login
Sample of hexagonal architecture to handle login logic and user CRUD

How to run
---
#### Set Environment Variable
This application support 2 kind of database Redis and MongoDB to prove that by using 1 repository logic it can connect into different database  
By default it will connect into our mongo DB database with default host & port `localhost:27017` and collection `local`  
To connect into different database we need to set database information in environment variable

`set mongo_url=mongodb://localhost:27017/local`
`set mongo_timeout=30`  
`set mongo_db=local`  
`set url_db=mongo`  

After setting the database information we only need to run the main.go file  
`go run main.go`  

Hexagonal Architecture Concept
---
The concept of Hexagonal Architecture is to make sure that we divide our software in such a way that each piece of the software maintains its separation of concerns so that our application is modular.  
We have App & Domain Logic in the middle and Ports & Adapters layer on the outside which connect the App & Domain Logic to outside things like a user interface, Repository, REST API, External API or Message Queue.  
You can refer to below images to make it clearly, I got [this image from medium post](https://medium.com/@msmechatronics/hexagonal-architecture-in-java-5b21ebea849d)
 :

![alt hexagonal-architecture](https://miro.medium.com/max/1689/1*dayDz6OTNc2qSS3QhppATA.png)

By using this arhitecture we can make sure that the business logic itself is independent of any kind of framework.  
So if we rely on framework to one of our Ports & Adapters, if the framework become depreciated or if want to use a different fremework we could just take our business logic and move it over to that other framework.  
We also can make sure that our App & Domain Logic are testable without any of the Port & Adapters,  so if we don't have a database or an API we can still test the business logic and make sure that it actually works properly

Here is the service that we are going to build  

So we have our service which is a user management and login and it will connect to serializer which will either serialize the data into json or message pack before serving it through REST API  
And then on the other side we have our repository which will either choose to use MongoDB or Redis based on how we start the application from command line.  
So basically our API will be able to accept JSON or message pack format and also our repository is able to use both MongoDB and Redis and it won't really affect our service