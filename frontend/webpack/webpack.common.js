const CopyWebpackPlugin = require('copy-webpack-plugin');
const ExtractTextPlugin = require('extract-text-webpack-plugin');

const path = require('path')
const here = __dirname
const root = path.resolve(path.join(here, `..`))
const output = path.join(root, 'dist')

module.exports = {
  context: root,
  resolve: {
    modules: ['node_modules', 'src'],
    extensions: ['.js']
  },
  output: {
    path: output,
    filename: 'bundle.js',
    publicPath: '/',
  },
  entry: 'index.js',
  module: {
    rules: [{
      test: require.resolve('jquery'),
      use: [{
        loader: 'expose-loader',
        query: 'jQuery',
      },{
        loader: 'expose-loader',
        query: '$',
      }]
    },{
      test: /\.js$/,
      use: {
        loader: 'babel-loader',
        options: {
          presets: ['env']
        }
      }
    },{
      test: /\.(png|jpg|jpeg|svg|woff|woff2|ttf|eot)$/i,
      use: {
        loader: 'url-loader',
        options: {
          name: '[name].[ext]',
          outputPath: 'assets',
          limit: 16 * 1024
        }
      }
    },{
      test: /\.css$/,
      use: ExtractTextPlugin.extract({
        use: {
          loader: 'css-loader',
          options: {
              url: true,
              minimize: true,
          }
        }
      })
    }]
  },
  plugins: [
    new CopyWebpackPlugin([
      { from: 'src/index.html' },
      { from: 'assets', to: 'assets' }
    ]),
    new ExtractTextPlugin('bundle.css')
  ]
};