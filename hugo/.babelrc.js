module.exports = {
  presets: [
    [
      '@babel/preset-env',
      {
        modules: false,
        forceAllTransforms: true,
        useBuiltIns: 'usage',
        corejs: 3,
        browserslist: 'default, not IE 11',
      },
    ],
  ],
  plugins: [
    '@babel/plugin-syntax-dynamic-import',
    '@babel/plugin-proposal-class-properties',
    '@babel/plugin-proposal-object-rest-spread',
    ['@babel/plugin-transform-react-jsx', { pragma: 'h' }],
  ],
};
