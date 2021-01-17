module.exports = {
  env: {
    browser: true,
    es2021: true,
    amd: true,
  },
  globals: {
    process: true,
  },
  extends: [
    'eslint:recommended',
    'prettier',
    'prettier/@typescript-eslint',
    'prettier/prettier',
    'prettier/react',
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 12,
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint'],
  rules: {
    'one-var': ['error', 'never'],
    'prefer-template': 'error',
  },
  overrides: [
    {
      files: './*.js',
      env: {
        node: true,
      },
    },
    {
      files: ['*.ts', '*.tsx'],
      extends: ['plugin:@typescript-eslint/recommended'],
      rules: {
        '@typescript-eslint/explicit-module-boundary-types': 'off',
      },
    },
  ],
};
