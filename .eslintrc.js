// https://eslint.org/docs/user-guide/configuring

module.exports = {
  extends: 'eslint:recommended',
  root: true,
  parser: 'babel-eslint',
  parserOptions: {
    sourceType: 'module'
  },
  env: {
    browser: true,
    es6: true
  },
  // required to lint *.vue files
  plugins: ['html'],
  // add your custom rules here
  rules: {
    'no-case-declarations': 0,
    // allow debugger during development
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off'
  }
};
