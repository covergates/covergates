module.exports = {
  devServer: {
    disableHostCheck: true,
    proxy: process.env.VUE_APP_PROXY
  },
  publicPath: process.env.NODE_ENV === 'production' ? '{{.}}' : process.env.BASE_URL,
  transpileDependencies: [
    'vuetify'
  ],
  configureWebpack: {
    devtool: 'source-map'
  }
};
