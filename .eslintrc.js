module.exports = {
  env: {
    es6: true,
    browser: true
  },
  extends: [
    'plugin:react/recommended',
    'plugin:react-hooks/recommended',
    'standard-with-typescript',
    'standard-react',
    'standard-jsx'
  ],
  plugins: ['react', 'react-hooks', '@typescript-eslint'],
  overrides: [{
    files: ['*.ts', '*.tsx'],
    parser: '@typescript-eslint/parser',
    parserOptions: { project: './tsconfig.json' }
  }],
  ignorePatterns: ['.eslintrc.js', 'dist', '.yarn/*', '.pnp.*'],
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
    ecmaFeatures: { jsx: true }
  },
  rules: {
    'react/react-in-jsx-scope': 'off'
  }
}
