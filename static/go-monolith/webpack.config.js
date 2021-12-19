const path = require('path');
const webpack = require('webpack');
const CssMinimizerPlugin = require("css-minimizer-webpack-plugin");

module.exports = {
  mode: "production", // "production" | "development" | "none",
  module: {
  rules: [{
    test: require.resolve('jquery'),
    use: [{
        loader: 'expose-loader',
        options: {exposes: ['$', 'jQuery']}
    }]
  }],  
  },
  entry: [
       './node_modules/jquery/dist/jquery.min.js', './node_modules/jquery-ui-bundle/jquery-ui.min.js', './node_modules/eonasdan-bootstrap-datetimepicker/build/js/bootstrap-datetimepicker.min.js', './js/sprintf.js', './assets/js/tether.min.js', './assets/bootstrap/3.3.7/js/bootstrap.min.js',
       './assets/js/wow.js', './assets/spinner/src/jRoll.js', './assets/js/floatHead.min.js', './assets/js/staticdata.js',
       './js/notify.min.js', './assets/chosen/docsupport/prism.js', './assets/cropper/cropper.min.js'
   ],
  // Chosen mode tells webpack to use its built-in optimizations accordingly.
  // defaults to ./src
  // Here the application starts executing
  // and webpack starts bundling
  output: {
    // options related to how webpack emits results
    path:path.resolve(__dirname, "dist"), // string (default)
    // the target directory for all output files
    // must be an absolute path (use the Node.js path module)
    filename: "[name].js", // string (default)
    // the filename template for entry chunks
    publicPath: "/assets/", // string
    // the url to the output directory resolved relative to the HTML page
    library: { // There is also an old syntax for this available (click to show)
      type: "umd", // universal module definition
      // the type of the exported library
      name: "MyLibrary", // string | string[]
      // the name of the exported library

      /* Advanced output.library configuration (click to show) */
    },
    uniqueName: "my-application", // (defaults to package.json "name")
    // unique name for this build to avoid conflicts with other builds in the same HTML
    // name of the configuration, shown in output
    /* Advanced output configuration (click to show) */
    /* Expert output configuration 1 (on own risk) */
    /* Expert output configuration 2 (on own risk) */
  },
  plugins: [
      new webpack.ProvidePlugin({
          $: "jquery",
          jquery: "jquery",
          "window.jQuery": "jquery",
          jQuery:"jquery"
      })  
  ]
}
