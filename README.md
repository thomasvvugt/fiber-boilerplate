# Fiber Boilerplate
A boilerplate for the Fiber web framework

## Configuration
Different to previous versions of this boilerplate, configurations are in a single file called `.env`. You can copy the `.env.example` and change it to your needs.

This `.env` file represents system environment variables on your machine. This change was made with the ease-of-use with Docker in mind.

A full version of all available configurations is located in the `.env.full` file. Various options can be changed depending on your needs such as Database, Fiber and Middleware settings.

Keep in mind if configurations are not set, they default to Fiber's default settings which can be found [here](https://docs.gofiber.io/).

## Routing
Routing examples can be found within the `/routes` directory. Both web and API routes are split, but you can adjust this to your likings.

## Views
Views are located and be edited under the `/resources/views` directory. 

You are able change this behavior using the `.env` file, as well as the ability to modify the Views Engine and other templating configurations using this file or using environment variables.

## Controllers
Example controllers can be found within the `/app/controllers` directory. You can extend or edit these to your preferences.

## Database
We use GORM v2 as an ORM to provide useful features to your models. Please check out their documentation [here](https://gorm.io/index.html).

## Models
Models are located within the `/app/models` directory and are also based on the GORM v2 package.

## Compiling assets
This boilerplate uses [Laravel Mix](https://github.com/JeffreyWay/laravel-mix) as an elegant wrapper around [Webpack](https://github.com/webpack/webpack) (a bundler for javascript and friends).

In order to compile your assets, you must first add them in the `webpack.mix.js` file. Examples of the Laravel Mix API can be found [here](https://laravel-mix.com/docs/5.0/mixjs).

Then you must run either `npm install` or `yarn install` to install the packages required to compile your assets.

Next, run one of the following commands to compile your assets with either `npm` or `yarn`:

```bash
# Run all Mix tasks
npm run dev

# Run all Mix tasks and minify output
yarn run production
# Run all Mix tasks and watch for changes (useful when developing)
yarn run watch
# Run all Mix tasks with hot module replacement
yarn run hot
```

## Docker
You can run your own application using the Docker example image.
To build and run the Docker image, you can use the following commands.

Please note, I am using host.docker.internal to point to my Docker host machine. You are free to use Docker's internal networking to point to your desired database host.


```bash
docker build -t fiber-boilerplate .
docker run -it --rm --name fiber-boilerplate -e DB_HOST=host.docker.internal -e DB_USER=fiber -e DB_PASSWORD=secret -e DB_DATABASE=boilerplate -p 8080:8080 fiber-boilerplate
```


## Live Reloading (Air)
Example configuration files for [Air](https://github.com/cosmtrek/air) have also been included.
This allows you to live reload your Go application when you change a model, view or controller which is very useful when developing your application.

To run Air, use the following commands. Also, check out [Air its documentation](https://github.com/cosmtrek/air) about running the `air` command.
```bash
# Windows
air -c .air.windows.conf
# Linux
air -c .air.linux.conf
```
