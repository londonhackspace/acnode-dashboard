const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack')

let gitHash = require('child_process')
    .execSync("git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty")
    .toString();

module.exports = {
    entry: './src/dash.ts',
    module: {
        rules: [
            // rule to run js and jsx through babel
            {
                test: /\.(js|jsx|ts|tsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        plugins: ['@babel/plugin-proposal-class-properties'],
                    },
                },
            },
            {
                test: /\.module\.css$/,
                use: [
                    'style-loader',
                    "@teamsupercell/typings-for-css-modules-loader",
                    {
                        loader: "css-loader",
                        options: { modules: true }
                    },
                ],
            },
            {
                test: /(?<!\.module)\.css$/,
                use: [
                    'style-loader',
                    "@teamsupercell/typings-for-css-modules-loader",
                    'css-loader',
                ]
            },
        ],
    },
    plugins: [
        // Make our index.html file
        new HtmlWebpackPlugin({
            title: 'ACNode Dashboard',
            template: "src/index.html",
        }),
        new webpack.DefinePlugin({
                gitHash: JSON.stringify(gitHash),
            }),
    ],
    resolve: { extensions: ["*", ".js", ".jsx", ".ts", ".tsx"] },
    output: {
        publicPath: '/',
        filename: 'static/[name].[contenthash].js',
        path: path.resolve(__dirname, 'dist'),
        clean: true,
    },
    optimization: {
        usedExports: true,
    },
};