## Prepare

- install node-sass
- install yarn
- yarn add bootstrap
- yarn add jquery
- yarn add marked
- yarn add hightlight.js
- yarn install

## Env variables

- PORT
- ConfigDir 
- ADMIN_USER
- ADMIN_PASSWORD

## Configfile 

1. $HOME/.fiberapp/tagconfig.txt

    ```
    tag1,tag2,tag3
    ```

1. $HOME/.fiberapp/dbconfig.txt

    ```
    DbHost="127.0.0.1"
    DbPort="27017"
    DbName="blog"
    DbUser="root"
    DbPassword="password"
    ```

## Usage

- `make assets` to bundle to app.css and app.js
- `make dev` to start the app
- `go to http://ip:port/login` to manage pages 

## LICENSE

MIT