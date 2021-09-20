// const fs = require('fs')
const path = require('path')
/* const { bold, greenBright, yellowBright } = require('colorette')

const jsgo = bold(greenBright('js.go'))
const generatedFromMain = bold(yellowBright('[generated from main.js]'))
class InlineGolangJsPlugin {
  apply (compiler) {
    compiler.hooks.done.tap('InlineGolangHtmlPlugin', stats => {
      const js = fs.readFileSync(path.resolve(__dirname, 'dist', 'main.js'), { encoding: 'utf8' })
      const goFile = `package main\n\nconst js = \`\n${js.replace(/`/g, '`+"`"+`')}\n\`\n`
      const buffer = Buffer.from(goFile)
      fs.writeFileSync('js.go', buffer)

      console.log(`asset ${jsgo} ${Math.ceil(buffer.byteLength / 1024)} KiB ${generatedFromMain}`)
    })
  }
} */

const isDev = env => (
  (env && env.NODE_ENV === 'development') ||
  process.env.NODE_ENV === 'development'
)

module.exports = env => ({
  entry: './src/index.js',
  mode: isDev(env) ? 'development' : 'production',
  devtool: isDev(env) ? 'inline-source-map' : undefined,
  output: {
    filename: '[name].js', // filename: '[name].[contenthash].js',
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
  // plugins: [new InlineGolangJsPlugin()]
})
