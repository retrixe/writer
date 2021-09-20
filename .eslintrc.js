module.exports = {
  env: {
    es6: true,
    browser: true
  },
  extends: ['plugin:react/recommended', 'standard', 'standard-react', 'standard-jsx'],
  plugins: ['react-hooks'],
  ignorePatterns: ['.eslintrc.js', 'dist', '.yarn/*', '.pnp.*'],
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
    ecmaFeatures: { jsx: true }
  },
  rules: {
    // React Hooks rules.
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn'
  }
}
