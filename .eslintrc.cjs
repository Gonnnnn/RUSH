module.exports = {
    root: true,
    env: { browser: true, es2020: true },
    extends: [
      'eslint:recommended',
      'plugin:@typescript-eslint/recommended',
      'plugin:react-hooks/recommended',
      'airbnb',
      'airbnb/hooks',
      'airbnb-typescript',
      'plugin:react/jsx-runtime', // https://github.com/jsx-eslint/eslint-plugin-react/blob/master/docs/rules/react-in-jsx-scope.md#when-not-to-use-it
      'prettier',
    ],
    ignorePatterns: ['dist', '.eslintrc.cjs', '**/generated/*'],
    parser: '@typescript-eslint/parser',
    parserOptions: {
      project: [
        'ui/tsconfig.json',
      ],
      tsconfigRootDir: __dirname,
    },
    plugins: ['react-refresh'],
    rules: {
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
      'react/function-component-definition': [
        2,
        {
          namedComponents: 'arrow-function',
          unnamedComponents: 'arrow-function',
        },
      ],
      'import/no-extraneous-dependencies': 'off',
      // We set this to 'off' to allow trivial variables, such as
      // styled components, to be placed at the bottom of a file.
      'no-use-before-define': 'off',
      '@typescript-eslint/no-use-before-define': 'off',
      'react/require-default-props': 'off',
      // These rules are for disabled people to read websites better, but we don't have resource yet to support it.
      'jsx-a11y/click-events-have-key-events': 'off',
      'jsx-a11y/no-noninteractive-element-interactions': 'off',
    },
  };
  