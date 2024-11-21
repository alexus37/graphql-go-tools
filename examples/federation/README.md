# Federation Demo

## Getting started

1. Install go modules

```shell
go mod download
```

2. Run start script

```
chmod +x start.sh
./start.sh
```

## To debug

1. Install go modules

```shell
go mod download
```

2. Run start script

```
chmod +x start.sh
./start_debug.sh
```

3. Launch the example federation server in vscode

Run the `Debug srv-gateway` in vscode.

4. Clone the playground plus plus repo and run the dev server

```shell
git clone https://github.com/alexus37/playground
npm start
```

5. Open chrome without CORS

```shell
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --disable-web-security --user-data-dir="/tmp"
```

6. Go to localhost:3000 and run the following query

```graphql
{
  me {
    username
    id
    reviews {
      body
    }
  }
  topProducts(first: 2) {
    upc
    reviews {
      body
    }
  }
}
```
