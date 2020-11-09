# hoopback
Hoopback is a webhook to webhook middleware, intended to bridge two webhooks that are not initally compatible. For example, Discord webhooks have a special form that likely doesn't match the form of a webhook from a service like Heroku. Therefore, Hoopback can convert the Heroku webhook payload to one that is Discord compatible and then send the webhook for you. More information at https://hoopback.schwa.tech.

# Getting Started
0. Create Discord OAuth App and a MongoDB instance
1. Clone the repository
```sh
$ git clone https://github.com/Kalissaac/hoopback && cd hoopback
```
2. Populate .env
```
.env
------------------------------------------
CLIENT_ID=<Discord OAuth client ID>
CLIENT_SECRET=<Discord OAuth client secret>
MONGODB_URI=<MongoDB connection URI>
PORT=<optional, port to run web server on>
```
3. Get dependencies (`npm` is used for TailwindCSS)
```sh
$ go get
$ npm install
```
4. Build project
```sh
$ npm run build
$ go build web.go
```
5. Run project
```sh
$ ./build/hoopback
```
