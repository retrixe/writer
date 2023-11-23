module.exports = {
  root: true,
  env: {
    es2024: true,
    browser: true
  },
  extends: [
    'plugin:react/recommended',
    'standard-with-typescript',
    'standard-react',
    'standard-jsx',
    'plugin:prettier/recommended',
  ],
  ignorePatterns: ['.eslintrc.js', 'dist'],
  overrides: [{ files: ['*.ts', '*.tsx'] }],
  parser: '@typescript-eslint/parser',
  parserOptions: { project: './tsconfig.json' },
  rules: {
    'react/react-in-jsx-scope': 'off',
    'react/no-unknown-property': ['error', { ignore: ['css'] }],
    // Make TypeScript ESLint less strict.
    '@typescript-eslint/no-confusing-void-expression': 'off',
  }
}
