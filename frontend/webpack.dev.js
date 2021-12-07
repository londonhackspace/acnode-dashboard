const { merge } = require('webpack-merge');
const path = require('path');
const common = require('./webpack.common.js');

module.exports = merge(common, {
    mode: 'development',
    devtool: 'inline-source-map',
    devServer: {
        client: {
            webSocketURL: {
                pathname: "/ws_debug"
            }
        },
        devMiddleware: {
            publicPath: "http://localhost:3000/",
        },
        static: {
            directory: path.join(__dirname, "static"),
        },
        port: 3000,
        hot: true,
        proxy: {
            '/api': 'http://localhost:8080',
            '/swagger': 'http://localhost:8080',
            '/static/swagger': 'http://localhost:8080',
            '/static/api.yaml': 'http://localhost:8080',
            '/ws': {
                target: 'ws://localhost:8080',
                ws: true
            },
        },
    },
});