const merge = require('webpack-merge');
const common = require('./webpack.common.js');
const path = require('path')

const here = __dirname
const config = path.join(here, 'configurations', 'development.js')

module.exports = merge(common, {
    resolve: {
        alias: { config: config }
    },
    devtool: 'inline-source-map',
});
    