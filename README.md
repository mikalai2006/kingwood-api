# Template for rest API

![example event parameter](https://github.com/mikalai2006/kingwood-api/actions/workflows/deploy.yml/badge.svg?branch=main&event=push)

For start project, need to create file .env in parent directory

```
HOST_DATA=***** // directory for save docker volume
APP_ENV=*****  // local or host

CLIENT_HOST=***** // Client host (Ex: http://localhost:3000)

HOST=***** // API host (Ex: localhost)
PORT=*****

// OAuth vk
VK_CLIENT_ID=*****
VK_CLIENT_SECRET=*****
// OAuth google
GOOGLE_CLIENT_ID=*****
GOOGLE_CLIENT_SECRET=*****
NAME_COOKIE_REFRESH=***** // name cookie for resresh token

SALT=***** // custom string
SIGNING_KEY=***** // custom string
NAME_COOKIE_REFRESH=***** // custom string

MONGODB_HOST=*****
MONGODB_PORT=*****
MONGODB_USER=*****
MONGODB_PASSWORD=*****
```

Option type action:
1 - remove node
2 - patch node

11 - remove tagopt
12 - add tagopt
