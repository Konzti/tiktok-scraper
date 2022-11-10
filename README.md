## Tiktok Video Downloader

This is a Tiktok Downloader written in Go accompanied by a React client. 

This app uses the awesome Gin framework
https://github.com/gin-gonic/gin

## Important:

This is not a Tiktok API and this app does not use any 3rd party API.


## Usage

### --- Development ---
#### 1. Run Server
```go build main.go -o tiktok```
<br/>
```./tiktok```

#### 2. Run Client
```cd client && yarn && yarn dev```
<br/>
<br/>
### --- Build For Production (using docker-compose)  ---
#### 1. Create volume (to persist and delete videos)
 ```docker volume create tiktok```
#### To skip the first step you can edit docker-compose.yml like so:
```
volumes:
    tiktok:
      external: false
```

#### 2. Build and run app containers
 ```docker-compose up -d --build```


