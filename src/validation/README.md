# libs/validation

The validation directive, `x-validate`, must be set on a form tag

- Sets `novalidate` on the form if JavaScript is supported, falls back to built-in form validation if JavaScript is disabled

- Makes use of the `type` attribute, and `required`, `pattern`, `min`, etc

- Rules like `required` must behave the same as built-in form validation. For example, if `type=checkbox`, then it must be checked to pass the rule

- Custom element and form level rules using `x-valid-RULE`

- Form level rules are set on hidden inputs

- All validated elements, including the form, must have unique IDs

- Validation messages are displayed inline, not in tooltips. This is mobile friendly, and makes the behavior more predictable

- Optionally makes use of [Alpine.js Components](https://alpinejs.dev/components), e.g. date picker. Enabled by default, can be disabled with a constructor option. Component dependencies must be included separately

- Provides a consistent validation UX across supported browsers


**TODO**

- [ ] Make styling interoperable with [built-in validation pseudo classes](https://developer.mozilla.org/en-US/docs/Learn/Forms/Form_validation#using_built-in_form_validation), `:valid` and `:invalid`?

