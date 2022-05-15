const path = require('path')

const isDev = env => (
  (env && env.NODE_ENV === 'development') ||
  process.env.NODE_ENV === 'development'
)

module.exports = env => ({
  entry: './src/index.js',
  mode: isDev(env) ? 'development' : 'production',
  devtool: isDev(env) ? 'inline-source-map' : undefined,
  output: {
    filename: 'index.js', // filename: '[name].[contenthash].js',
    path: path.resolve(__dirname, 'dist')
  },
  module: {
    rules: [{
      test: /\.(c|m)?js$/,
      exclude: /(node_modules|bower_components)/,
      use: {
        loader: 'babel-loader',
        options: {
          presets: [
            ['@babel/preset-react', { runtime: 'automatic', importSource: '@emotion/react' }],
            ['@babel/preset-env', { corejs: 3, targets: { ie: 9 }, useBuiltIns: 'usage' }]
          ]
        }
      }
    }]
  }
})
