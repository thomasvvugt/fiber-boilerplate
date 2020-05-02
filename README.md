# Fiber Boilerplate
A boilerplate for the Fiber web framework


## Configuration
All configuration for your application can be found in the `config.yml` file. Various options can be changed depending on your needs such as Database settings, Fiber settings and Fiber Middleware setting such as Logger, Public and Helmet.

Keep in mind if configurations are not set, they default to Fiber's default settings which can be found [here](https://docs.gofiber.io/).


## Compiling assets
This boilerplate uses [Laravel Mix](https://github.com/JeffreyWay/laravel-mix) as an elegant wrapper around [Webpack](https://github.com/webpack/webpack) (a bundler for javascript and friends).

In order to compile your assets, you must first add them in the `webpack.mix.js` file. Examples of the Laravel Mix API can be found [here](https://laravel-mix.com/docs/5.0/mixjs).

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
