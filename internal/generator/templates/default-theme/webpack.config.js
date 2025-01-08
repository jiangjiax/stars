const { WebpackManifestPlugin } = require('webpack-manifest-plugin');
const path = require('path');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  entry: {
    main: './static/js/main.js',
    nft: './static/js/nft.js'
  },
  output: {
    filename: '[name].[contenthash].bundle.js',
    path: path.resolve(__dirname, './static/dist'),
    publicPath: ''
  },
  devtool: process.env.NODE_ENV === 'development' ? 'eval-source-map' : false,
  optimization: {
    minimize: process.env.NODE_ENV === 'production',
    minimizer: [new TerserPlugin()]
  },
  performance: {
    hints: process.env.NODE_ENV === 'production' ? 'warning' : false
  },
  stats: {
    colors: true,
    modules: false,
    children: false,
    chunks: false,
    chunkModules: false
  },
  mode: process.env.NODE_ENV || 'development',
  plugins: [
    new WebpackManifestPlugin({
      fileName: path.resolve(__dirname, './static/dist/manifest.json'),
      writeToFileEmit: true,
      generate: (seed, files) => {
        const manifestFiles = files.reduce((manifest, file) => {
          const key = file.name.replace('.js', '') + '.bundle.js';
          manifest[key] = file.path;
          return manifest;
        }, seed);
        console.log('Generated manifest:', manifestFiles);
        return manifestFiles;
      }
    })
  ]
}; 