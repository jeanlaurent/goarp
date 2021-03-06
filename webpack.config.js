module.exports = {
    entry: [
      './src/index.js'
    ],
    module: {
        rules: [
          {
            test: /\.(js|jsx)$/,
            exclude: /node_modules/,
            use: ['babel-loader']
          }
        ]
      },
      resolve: {
        extensions: ['*', '.js', '.jsx']
      },
    output: {
      path: __dirname + '/app/',
      publicPath: '/',
      filename: 'goarp.js'
    },
    devServer: {
      contentBase: './app/'
    }
  };