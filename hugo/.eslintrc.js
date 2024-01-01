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
    'plugin:react/recommended',
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
  plugins: ['@typescript-eslint', 'react'],
  rules: {
    'one-var': ['error', 'never'],
    'prefer-template': 'error',
    'react/prop-types': 'off',
    'no-use-before-define': 'error',
  },
  settings: {
    react: {
      pragma: 'h',
      fragment: 'Fragment',
      version: '17',
    },
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
