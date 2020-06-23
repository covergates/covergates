const path = require('path');
module.exports = {
  // set your styleguidist configuration here
  title: 'Default Style Guide',
  // serverPort: 8081,
  components: 'src/components/**/[A-Z]*.vue',
  styleguideDir: 'styleguide/dist',
  ignore: ['**/__tests__/**', '**/__examples__/**'],
  require: [
    path.join(__dirname, 'styleguide/global.requires.ts'),
    path.join(__dirname, 'styleguide/router-mock.ts')
  ],
  validExtends: fullFilePath => !/(?=node_modules)(?!node_modules\/vuetify)/.test(fullFilePath),
  renderRootJsx: path.join(__dirname, 'styleguide/styleguide.root.ts'),
  exampleMode: 'expand'
};
