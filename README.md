[![Golang](https://img.shields.io/github/go-mod/go-version/nskryukov/chatgpt_tg_bot)](https://go.dev/blog/go1.18)
[![Tarantool](https://img.shields.io/badge/Tarantool-v2.10-blue)](https://www.tarantool.io/ru/doc/latest/release/2.10.0/)
[![Tarantool Go Connector](https://img.shields.io/badge/Tarantool%20Go%20Connector-v1.10-blue)](https://github.com/tarantool/go-tarantool)
[![Bot_API](https://img.shields.io/badge/Telegram%20Bot%20API-May%2C%202023-blue)](https://core.telegram.org/bots/api)
[![Docker](https://img.shields.io/badge/Docker-v23.0.5-blue)](https://www.docker.com/)

# Password manager Telegram bot
This bot is simple to use password manager, you can save, retrieve and delete passwords there. Bot uses Tarantool DB to store data. Also, it works in Docker containers 


**Try this bot**  
[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/tarantool_password_manager_bot)

## Available bot commands
1. ```/start``` - Bot sends greeting message
2. ```/set``` - Save service login and password. Bot sends info about message format to save password and user entering setting menu 
3. ```/get``` - Get service login and password. Bot sends info about message format to get password and user entering getting menu
4. ```/del``` - Delete service login and password. Bot sends info about message format to delete password and user entering deleting menu
5. ```/exit``` - Exit interactive menu called by commands ```/set``` ```/get``` ```/del```

## Quick start
**Required software**
1. Docker v23.0.5
2. Go v1.18

**Variables to specify in docker-compose.yml**
1. ```<TARANTOOL_USER_NAME>``` - db username
2. ```<TARANTOOL_USER_PASSWORD>``` - db user password
3. ```<CHAT_BOT_TOKEN>``` - your tg bot token
4. ```<PATH_TO_DB_LOG>``` - path directory on your server to save Tarantool snapshots and WAL logs
5. ```<PATH_TO_BOT_LOG>``` - path directory on your server to save bot logs

**Docker image with bot**
1. Go to directory with project
2. Build Go project ```go build -tags go_tarantool_ssl_disable .```
3. Specify in Dockerfile ```<PATH_TO_BOT_BINARY_FILE>``` - path to built on step 1 file
4. Build Docker image ```docker build -t password_manager .```

After preparations go to directory with bot and run ```docker-compose build``` and ```docker-compose up```