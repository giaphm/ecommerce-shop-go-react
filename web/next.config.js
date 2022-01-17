/** @type {import('next').NextConfig} */
module.exports = {
  distDir: '../.next',
  reactStrictMode: true,
  webpackDevMiddleware: config => {
    config.watchOptions = {
      poll: 800,
      aggregateTimeout: 300,
    }
    return config
  },
}
