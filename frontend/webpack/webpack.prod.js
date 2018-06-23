const merge = require('webpack-merge');
const common = require('./webpack.common.js');
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');
const path = require('path')

const here = __dirname
const config = path.join(here, 'configurations', 'production.js')

module.exports = merge(common, {
    resolve: {
        alias: { config: config }
    },
    plugins: [
        new UglifyJSPlugin(),
    ]
});