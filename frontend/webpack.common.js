const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

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
        ],
    },
    plugins: [
        // Make our index.html file
        new HtmlWebpackPlugin({
            title: 'ACNode Dashboard',
            template: "src/index.html",
        }),
    ],
    resolve: { extensions: ["*", ".js", ".jsx", ".ts", ".tsx"] },
    output: {
        filename: 'static/[name].[contenthash].js',
        path: path.resolve(__dirname, 'dist'),
        clean: true,
    },
    optimization: {
        usedExports: true,
    },
};