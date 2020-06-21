module.exports = {
  devServer: {
    disableHostCheck: true
  },
  publicPath: process.env.NODE_ENV === 'production' ? '{{.}}' : process.env.BASE_URL,
  transpileDependencies: [
    'vuetify'
  ],
  configureWebpack: {
    devtool: 'source-map'
  }
};
