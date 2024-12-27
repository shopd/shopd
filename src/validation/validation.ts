import { Alpine } from "../../deps.ts"
import { sprintf } from "../sprintf/sprintf.esm.js"
import { Events } from "../events/events.ts"
import { alpineInit } from "../events/alpine.ts"
import { validationValidate } from "./events.ts"
import { submit, change, keyup } from "../events/standard.ts"

// See "TypeScript Style Guide" for naming identifiers
// https://google.github.io/styleguide/tsguide.html#identifiers

// TypeScript contribution guide says "use undefined. Do not use null", but also
// "This is NOT a prescriptive guideline for the TypeScript community. These 
// guidelines are meant for contributors to the TypeScript project's codebase"
// https://github.com/Microsoft/TypeScript/wiki/Coding-guidelines#null-and-undefined
// "undefined means a variable has been declared but not assigned a value...
// null can be assigned to a variable as a representation of no value".
// Therefore pass or assign null if value is not applicable
// https://stackoverflow.com/a/5076962/639133

// ValidationRule is a function that accepts any value,
// it returns a ValidationResult
export interface ValidationRule {
	// deno-lint-ignore no-explicit-any
	(value: any): ValidationResult
}

// ValidationRuleAsync is a function that accepts any value,
// it returns a Promise that resolved to a ValidationResult
export interface ValidationRuleAsync {
	// deno-lint-ignore no-explicit-any
	(value: any): Promise<ValidationResult>
}

// ValidationResult for use with ValidationRule
export interface ValidationResult {
	rule: RuleName
	valid: boolean
	msg: string
	// deno-lint-ignore no-explicit-any
	value: any
}

export type RuleName = string

export type RuleMap = Map<RuleName, ValidationRule>

export type RuleMapAsync = Map<RuleName, ValidationRuleAsync>

// CustomRuleMap is used to register custom rules for a field,
// and report undefined custom rule names
type CustomRuleMap = Map<RuleName, boolean>

interface Field {
	// timeout is used to debounce events on a field
	timeout: null | number
	// compareValue is used to check if a field is "dirty".
	// Multiple events (submit, change, keyup) can cause a validation check.
	// The check only proceeds if the comparison values do not match.
	// This value might not match the actual value of the element,
	// consider comments for "textarea"
	compareValue: null | string
	// result of validation on latest change
	result?: ValidationResult
	// invalidMsg to display if not valid
	invalidMsg?: HTMLElement
	// validMsg to display if valid
	validMsg?: HTMLElement
	// customRuleNames associated with this field
	customRuleMap?: CustomRuleMap
	// TODO Triggers used as a verb, rather use a noun?
	// triggers validation on other elements
	triggers?: HTMLElement[]
}

type ElementID = string

type FieldMap = Map<ElementID, Field>

interface Form {
	// timeout is used to debounce events on a form
	timeout: null | number
	// fields for this form
	fields: FieldMap
}

type FormMap = Map<ElementID, Form>

// DebounceCB is the debounce callback function
export interface DebounceCB {
	(): void
}

// ValidFieldCB is called if a field is valid
export interface ValidFieldCB {
	(target: HTMLElement): void
}

// .............................................................................

export const RuleEmail = "email"
export const MsgEmail = "Email is invalid"

export const RuleNotApplicable = "n/a"
export const MsgNotApplicable = ""

export const RulePattern = "pattern"
export const MsgPattern = "Must match pattern"

export const RuleRequired = "required"
export const MsgRequired = "Required"

export function getDefaultRules(): RuleMap {
	const rules = new Map<RuleName, ValidationRule>()
	rules.set(RuleEmail, (value: string): ValidationResult => {
		// Required is checked separately
		if (value.trim() == "") {
			return {
				rule: RuleEmail,
				valid: true,
				msg: "",
				value: value
			}
		}
		// "Top-level domain has two to six letters... 
		// If you want to avoid sending too many undeliverable emails, 
		// while still not blocking any real email addresses"
		// https://www.oreilly.com/library/view/regular-expressions-cookbook/9780596802837/ch04.html#validation-email-solution-tld
		// Read more about JavaScript regex here
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/RegExp/test
		// Note the flag in the second arg to make the match case-insensitive
		const regex = new RegExp(/^[\w!#$%&'*+/=?`{|}~^-]+(?:\.[\w!#$%&'*+/=?`{|}~^-]+)*@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}$/, "i")
		return {
			rule: RuleEmail,
			valid: regex.test(value),
			msg: MsgEmail,
			value: value
		}
	})
	return rules
}

// .............................................................................

// ValidationErr must be returned (by functions or methods) on error...
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error#custom_error_types
// ...instead of null or undefined
// https://stackoverflow.com/a/48197438/639133
// Callers can then use "instanceof ValidationErr" to check for errors
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/instanceof
// "explicitly check for errors where they occur [instead of] 
// throwing exceptions and sometimes catching them"
// https://go.dev/blog/error-handling-and-go
export class ValidationErr extends Error {
	target: HTMLElement | null = null

	// deno-lint-ignore no-explicit-any
	constructor(target: HTMLElement | null, ...params: any) {
		super(...params)
		if (target) {
			this.target = target
		}
	}

	// error can be used for logging, 
	// i.e. console.error(...validationError.error())
	// deno-lint-ignore no-explicit-any
	error(): any[] {
		if (this.target) {
			// Target goes after message,
			// it displays better that way in the Dev Console
			return [this.message, this.target]
		}
		return [this.message]
	}

	string(): string {
		if (this.target) {
			return sprintf("#%s %s", this.target.id, this.message)
		}
		return this.message
	}
}

// .............................................................................

// LanguageTranslation for validation messages
export interface LanguageTranslation {
	(result: ValidationResult): string
}

export type TranslationMap = Map<RuleName, LanguageTranslation>

// .............................................................................

export interface ValidationOptions {
	InvalidClass?: string
	CustomRules?: RuleMap
	CustomRulesAsync?: RuleMapAsync
}

export class Validation {
	#initialised = false

	#invalidClass: string

	// "Unlike TypeScript private, 
	// JavaScript private fields remain private after compilation"
	// https://www.typescriptlang.org/docs/handbook/2/classes.html#private
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes/Private_class_fields
	#rules: RuleMap

	#rulesAsync: RuleMapAsync

	// forms using the validation directive
	#forms: FormMap

	// translations for validation messages
	#translations: TranslationMap

	#debounceWaitMS: number

	#debug: boolean

	constructor(options?: ValidationOptions) {
		// Bulma classes by default
		// https://bulma.io/documentation/form/general/
		this.#invalidClass = "is-danger"

		this.#rules = getDefaultRules()

		this.#rulesAsync = new Map<RuleName, ValidationRuleAsync>()

		this.#forms = new Map<ElementID, Form>()

		// TODO Language translation
		this.#translations = new Map<RuleName, LanguageTranslation>()

		// TODO Option to adjust debounceWaitMS
		this.#debounceWaitMS = 900

		// TODO Option to toggle debug
		this.#debug = false

		// Options override defaults
		if (options) {
			if (options.InvalidClass) {
				this.#invalidClass = options.InvalidClass
			}
			if (options.CustomRules) {
				options.CustomRules.forEach((value, key) => {
					this.#rules.set(key, value)
				})
			}
			if (options.CustomRulesAsync) {
				options.CustomRulesAsync.forEach((value, key) => {
					this.#rulesAsync.set(key, value)
				})
			}
		}
	}

	// deno-lint-ignore no-explicit-any
	static isErr(err: any): boolean {
		if (err instanceof ValidationErr) {
			return true
		}
		return false
	}

	#getFormID(target: HTMLElement): string | ValidationErr {
		const form = target.closest("form")
		if (form == null) {
			return new ValidationErr(target, "element not inside a form")
		}
		return form.id
	}

	#getFields(formID: ElementID): FieldMap | ValidationErr {
		const form = this.#forms.get(formID)
		if (form == undefined) {
			return new ValidationErr(
				null, sprintf("form %s is not registered", formID))
		}
		return form.fields
	}

	#getField(formID: ElementID, fieldID: ElementID): Field | ValidationErr {
		const form = this.#forms.get(formID)
		if (form == undefined) {
			return new ValidationErr(
				null, sprintf("form %s is not registered", formID))
		}

		const field = form.fields.get(fieldID)
		if (field == undefined) {
			return new ValidationErr(
				null, sprintf("field %s is not registered", fieldID))
		}

		return field
	}

	getFieldElementIDs(formID: ElementID): ElementID[] | ValidationErr {
		const form = document.getElementById(formID)
		if (form == null) {
			return new ValidationErr(
				null, sprintf("form %s not found", formID))
		}
		const selector = "input, select, textarea"
		const elements = form.querySelectorAll(selector)
		const ids: ElementID[] = []
		elements.forEach((element) => {
			const id = (element as HTMLElement).id
			if (id.trim() != "") {
				// Ignore fields with empty id
				ids.push(id)
			}
		})
		return ids
	}

	getFieldElement(formID: ElementID, fieldID: ElementID): HTMLElement | ValidationErr {
		// https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelector
		const selector = sprintf("#%s #%s", formID, fieldID)
		const element = document.querySelector(selector)
		if (element == null) {
			return new ValidationErr(
				null, sprintf("selector \"%s\" not found", selector))
		}
		return element as HTMLElement
	}

	// checkRule checks if a value is valid for a rule
	// deno-lint-ignore no-explicit-any
	checkRule(name: RuleName, value: any): ValidationResult | ValidationErr {
		const rule = this.#rules.get(name)
		if (rule == undefined) {
			return new ValidationErr(
				null, sprintf("rule %s is undefined", name))
		}
		return rule(value)
	}

	// visible returns true if the element is visible
	#visible(e: HTMLElement): boolean {
		// TODO Use checkVisibility method if available?
		// https://stackoverflow.com/a/72717388/639133

		// Fall back to same check as jQuery source code
		// https://github.com/jquery/jquery/blob/main/src/css/hiddenVisibleSelectors.js
		return !!(e.offsetWidth || e.offsetHeight || e.getClientRects().length)
	}

	// hiddenInput return true if element is an input with type hidden
	#hiddenInput(e: HTMLElement): boolean {
		const tagName = e.tagName.toLowerCase()
		if (tagName == "input") {
			const type = (e as HTMLInputElement).getAttribute("type")
			if (type == "hidden") {
				return true
			}
		}
		return false
	}

	// Consider the built-in checkValidity and reportValidity functions
	// https://developer.mozilla.org/en-US/docs/Web/HTML/Constraint_validation

	// checkElement checks if the element is valid, 
	// and updates the validation flag on the local scope
	checkElement(target: HTMLElement, cb: null | ValidFieldCB) {
		if (cb == null) {
			// Callback does nothing
			cb = () => { }
		}

		if (!this.#visible(target)) {
			if (this.#hiddenInput(target)) {
				// Inputs with type hidden are always checked.
				// These can be used for "form level" rules
			} else {
				// Invisible fields are always valid.
				// Assuming messages associated with this field are also invisible,
				// therefore #reportElement is not called from here
				cb(target)
				return
			}
		}

		// Valid by default
		let result: ValidationResult = {
			rule: RuleNotApplicable,
			valid: true,
			msg: "",
			value: undefined
		}

		// Get IDs
		const formID = this.#getFormID(target)
		if (formID instanceof ValidationErr) {
			console.error(...formID.error())
			return
		}
		const fieldID = target.id
		if (fieldID.trim() == "") {
			console.error(sprintf("missing id in form %s for field", formID), target)
			return
		}

		// Get field
		const field = this.#getField(formID, fieldID)
		if (field instanceof ValidationErr) {
			console.error(...field.error())
			return
		}

		// Checks similar to built-in browser behavior
		const tagName = target.tagName.toLowerCase()
		// deno-lint-ignore no-explicit-any
		let value: any
		switch (tagName) {
			// deno-lint-ignore no-case-declarations
			case "input":
				const targetType = target.getAttribute("type")?.toLowerCase()
				if (targetType == "checkbox") {
					value = String((target as HTMLInputElement).checked)
				} else if (targetType == "hidden") {
					// Hidden inputs can't be changed by the user,
					// assume the value is dirty if this code is called.
					// Typical use-case is validating a combination of other elements
					value = (target as HTMLInputElement).value
					result = this.#checkInput(target as HTMLInputElement)
					break
				} else {
					value = (target as HTMLInputElement).value
				}
				if (field.compareValue == value) {
					if (field.result?.valid) {
						// Value hasn't changed, report on previous validation result
						cb(target)
					}
					return
				} else {
					field.compareValue = value
				}
				result = this.#checkInput(target as HTMLInputElement)
				break

			case "select":
				value = (target as HTMLSelectElement).value
				if (field.compareValue == value) {
					if (field.result?.valid) {
						// Value hasn't changed, report on previous validation result
						cb(target)
					}
					return
				} else {
					field.compareValue = value
				}
				result = this.#checkSelect(target as HTMLSelectElement)
				break

			case "textarea":
				// TODO Compare hashCode to check for changes
				// https://stackoverflow.com/a/7616484/639133
				value = (target as HTMLSelectElement).value
				if (field.compareValue == value) {
					if (field.result?.valid) {
						// Value hasn't changed, report on previous validation result
						cb(target)
					}
					return
				} else {
					field.compareValue = value
				}
				result = this.#checkTextarea(target as HTMLTextAreaElement)
				break

			default:
				console.error(
					sprintf("unknown tag %s for target", tagName), target)
		}
		if (!result.valid) {
			this.#reportElement(formID, fieldID, result)
			return
		}

		// Additional rules are checked after emulating built-in behavior
		if (field.customRuleMap) {
			// Avoid throwing exceptions in libs, but support it for user-defined code
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/try...catch

			// Reset
			for (const ruleName of field.customRuleMap.keys()) {
				field.customRuleMap.set(ruleName, false)
			}

			// Custom rules
			try {
				for (const ruleName of field.customRuleMap.keys()) {
					const rule = this.#rules.get(ruleName)
					if (rule != undefined) {
						// Custom rule is defined
						field.customRuleMap.set(ruleName, true)
						result = rule(value)
						if (!result.valid) {
							this.#reportElement(formID, fieldID, result)
							return
						}
					}
				}
			} catch (err) {
				console.error(err)
			}

			// Defer error logging and element reporting for async custom rules
			const defer = (formID: ElementID, fieldID: ElementID, result: ValidationResult) => {
				const field = this.#getField(formID, fieldID)
				if (field instanceof ValidationErr) {
					console.error(...field.error())
					return
				}

				if (field.customRuleMap) {
					// Log errors for undefined custom rules registered on this field
					for (const ruleName of field.customRuleMap.keys()) {
						const defined = field.customRuleMap.get(ruleName)
						if (!defined) {
							console.error(sprintf(
								"rule %s on field %s is undefined", ruleName, fieldID))
						}
					}
				}

				this.#reportElement(formID, fieldID, result)
				if (result.valid) {
					if (cb != null) {
						cb(target)
					}
				}
			}

			// Async custom rules
			const promises: Promise<ValidationResult>[] = []
			for (const ruleName of field.customRuleMap.keys()) {
				const rule = this.#rulesAsync.get(ruleName)
				if (rule != undefined) {
					// Custom rule is defined
					if (this.#debug) { console.info("ruleName", ruleName) }
					field.customRuleMap.set(ruleName, true)
					promises.push(rule(value))
				}
			}
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/allSettled
			Promise.allSettled(promises).then(
				(results) => {
					for (const r of results) {
						switch (r.status) {
							case "fulfilled":
								result = r.value
								if (!result.valid) {
									// Field is invalid, stop iterating results
									defer(formID as ElementID, fieldID as ElementID, result)
									return
								}
								break
							case "rejected":
								if (r.reason instanceof Error) {
									console.error(r.reason.message)
								} else {
									console.error("promise rejected", r.reason)
								}
								// Field is valid by default
								defer(formID as ElementID, fieldID as ElementID, result)
								return
						}
					}
					// Field is valid by default
					defer(formID as ElementID, fieldID as ElementID, result)
				})

		} else {
			// No custom rules, or triggers, and field is valid by default
			this.#reportElement(formID, fieldID, result)
			cb(target)
		}

		// Trigger validation on related elements
		if (field.triggers) {
			field.triggers.forEach((el) => {
				Events.dispatchOnElement(Validation.name, el, validationValidate)
			})
		}
	}

	// TODO Constructor option to toggle "is-hidden" and "is-invisible"
	#show(msg: HTMLElement) {
		// msg.style.removeProperty("visibility")
		// msg.style.removeProperty("display")
		// Using bulma
		msg.classList.remove("is-hidden")
	}

	#hide(msg: HTMLElement) {
		// msg.style.setProperty("visibility", "hidden")
		// msg.style.setProperty("display", "none")
		msg.classList.add("is-hidden")
	}

	#reportElement(formID: ElementID, fieldID: ElementID, result: ValidationResult) {
		const field = this.#getField(formID, fieldID)
		if (field instanceof ValidationErr) {
			console.error(...field.error())
			return
		}

		// Update current result
		field.result = result

		// Update visual validation indicators
		const el = this.getFieldElement(formID, fieldID)
		if (el instanceof ValidationErr) {
			console.error(...el.error())
			return
		}
		let msg = result.msg
		const translate = this.#translations.get(result.rule)
		if (translate != undefined) {
			msg = translate(result)
		}
		if (result.valid) {
			// https://developer.mozilla.org/en-US/docs/Web/API/Element/classList
			el.classList.remove(this.#invalidClass)
			// https://developer.mozilla.org/en-US/docs/Web/API/CSSStyleDeclaration/setProperty
			if (field.validMsg != undefined) {
				field.validMsg.innerText = msg
				this.#show(field.validMsg)
			}
			if (field.invalidMsg != undefined) {
				this.#hide(field.invalidMsg)
			}

		} else {
			el.classList.add(this.#invalidClass)
			if (field.validMsg != undefined) {
				this.#hide(field.validMsg)
			}
			if (field.invalidMsg != undefined) {
				field.invalidMsg.innerText = msg
				this.#show(field.invalidMsg)
			}
		}
	}

	#checkInput(target: HTMLInputElement): ValidationResult {
		// Valid by default
		let result: ValidationResult = {
			rule: RuleNotApplicable,
			valid: true,
			msg: MsgNotApplicable,
			value: undefined
		}

		const targetType = target.getAttribute("type")?.toLowerCase()

		// Required is checked first
		// https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/required
		const required = target.getAttribute("required")
		if (required != null) {
			const invalid = {
				rule: RuleRequired,
				valid: false,
				msg: MsgRequired,
				value: undefined
			}
			if (targetType == "checkbox") {
				if (!target.checked) {
					return invalid
				}
			} else {
				if (target.value.trim() == "") {
					return invalid
				}
			}
		}

		// Check regex if pattern attribute is set
		// https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/pattern
		// "the 'u' flag is specified... so that the pattern is treated as a 
		// sequence of Unicode code points, instead of as ASCII"
		// "If the specified pattern is invalid... this attribute is ignored"
		// TODO Option to make pattern matching case-sensitive?
		const pattern = target.getAttribute("pattern")?.toLowerCase()
		if (pattern != null && pattern.trim() != "") {
			try {
				const regex = new RegExp(sprintf("^%s$", pattern), "u")
				if (!regex.test(target.value.toLowerCase())) {
					return {
						rule: RulePattern,
						valid: false,
						msg: MsgPattern,
						value: undefined
					}
				}
			} catch (err) {
				console.error(err)
				return result
			}
		}

		// Checks specific to input type
		let typeResult: ValidationResult | ValidationErr
		switch (targetType) {
			case "checkbox":
				break
			case "email":
				typeResult = this.checkRule("email", target.value)
				if (typeResult instanceof ValidationErr) {
					// Just log the error don't return,
					// element are valid by default
					console.error(...typeResult.error())
					break
				}
				result = typeResult
				break
			case "file":
				// Can be validated with the required attr
				break
			case "hidden":
				break
			case "number":
				// TODO Create a rule for this type? Maybe not necessary,
				// Google Chrome only allows values like: "123", "1.23" and "1,23"
				break
			case "radio":
				break
			case "text":
				break
			default:
				console.error(
					sprintf("unknown type %s for target", targetType), target)
		}

		return result
	}

	#checkSelect(target: HTMLSelectElement): ValidationResult {
		// Valid by default
		const result: ValidationResult = {
			rule: RuleNotApplicable,
			valid: true,
			msg: MsgNotApplicable,
			value: undefined
		}

		// Required is checked first
		const required = target.getAttribute("required")
		if (required != null) {
			if (target.value.trim() == "") {
				return {
					rule: RuleRequired,
					valid: false,
					msg: MsgRequired,
					value: undefined
				}
			}
		}

		// Check regex if pattern attribute is set
		const pattern = target.getAttribute("pattern")?.toLowerCase()
		if (pattern != null && pattern.trim() != "") {
			try {
				const regex = new RegExp(sprintf("^%s$", pattern), "u")
				if (!regex.test(target.value.toLowerCase())) {
					return {
						rule: RulePattern,
						valid: false,
						msg: MsgPattern,
						value: undefined
					}
				}
			} catch (err) {
				console.error(err)
				return result
			}
		}

		return result
	}

	#checkTextarea(target: HTMLTextAreaElement): ValidationResult {
		// Valid by default
		const result: ValidationResult = {
			rule: RuleNotApplicable,
			valid: true,
			msg: MsgNotApplicable,
			value: undefined
		}

		// Required is checked first
		const required = target.getAttribute("required")
		if (required != null) {
			if (target.value.trim() == "") {
				return {
					rule: RuleRequired,
					valid: false,
					msg: MsgRequired,
					value: undefined
				}
			}
		}

		return result
	}

	// ...........................................................................

	// registerForm, if the form was registered before, 
	// delete it and create as new 
	#registerForm(formID: ElementID) {
		// fieldMap items are registered on change.
		// Consider that the elements inside the form could be 
		// created or removed dynamically with JavaScript (i.e. Alpine.js)
		const fieldMap = new Map<ElementID, Field>()
		const form = <Form>{ fields: fieldMap }

		// Register form
		this.#forms.set(formID, form)

		// Register fields
		const ids = this.getFieldElementIDs(formID)
		if (ids instanceof ValidationErr) {
			console.error(...ids.error())
			return
		}
		for (const fieldID of ids) {
			this.#registerField(formID, fieldID)
		}
	}

	// registerField or return existing
	#registerField(formID: ElementID, fieldID: ElementID): Field | ValidationErr {
		const form = this.#forms.get(formID)
		if (form == undefined) {
			// Directives declared on form fields (x-valid, x-invalid, etc) 
			// must be called after the form directive x-validate.
			// Does alpine call them in the order they're listed on the page?
			return new ValidationErr(
				null, sprintf("form %s is not registered", formID))
		}

		let field = form.fields.get(fieldID)
		if (field == undefined) {
			// This field has not been registered
			field = <Field>{
				timeout: null,
				compareValue: null,
				result: <ValidationResult>{
					rule: "",
					valid: true,
					msg: "",
					value: "",
				}
			}
			form.fields.set(fieldID, field)
		}

		return field
	}

	// debounce events by formID and fieldID.
	// Pass empty string for fieldID is not applicable
	// This is why it's required https://alpinejs.dev/directives/on#debounce
	// Debounce logic inspired by node_modules/alpinejs/src/utils/debounce.js
	debounce(formID: string, fieldID: string, cb: DebounceCB, wait: number):
		void | ValidationErr {

		// Get form
		const form = this.#forms.get(formID)
		if (form == undefined) {
			return new ValidationErr(
				null, sprintf("form %s is not registered", formID))
		}

		// Get current timeout
		let timeout: null | number
		const field = form.fields.get(fieldID)
		if (fieldID.trim() == "") {
			timeout = form.timeout
		} else {
			// Get field
			if (field == undefined) {
				return new ValidationErr(
					null, sprintf("field %s is not registered", fieldID))
			}
			timeout = field.timeout
		}
		if (timeout) {
			// Clear timeout
			clearTimeout(timeout)
		}

		// Register new timeout
		const later = () => {
			cb.apply(this, [])
		}
		if (fieldID.trim() == "") {
			form.timeout = setTimeout(later, wait)
		} else {
			if (field != undefined) {
				field.timeout = setTimeout(later, wait)
			} else {
				// This code should be unreachable, but return an error anyway
				return new ValidationErr(
					null, sprintf("field %s is undefined", fieldID))
			}
		}

		return
	}

	// register Alpine.js directives etc
	init() {
		if (this.#initialised) {
			console.error("already initialised")
			return
		}

		Events.addListener(Validation.name, alpineInit, () => {
			// .......................................................................
			// x-validate
			// Validation directive for forms
			Alpine.directive('validate', (el: Node) => {
				if (el.nodeName.toLowerCase() != "form") {
					console.error("element is not a form ", el)
					return
				}

				const form = (el as HTMLFormElement)

				if (form.id.trim() == "") {
					console.error("form must have an id ", el)
					return
				}

				// Disabling built-in form validation is required, see README
				form.setAttribute("novalidate", "")

				this.#registerForm(form.id)

				// .....................................................................
				// Event listeners

				// Submit event
				// https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement/submit_event
				Events.addListenerToElement(Validation.name, 
					form, submit, (event: Event) => {

					// https://developer.mozilla.org/en-US/docs/Web/API/Event/cancelable
					if (event.cancelable) {
						event.preventDefault()
					}

					if (this.#debug) {
						console.info("submit", (event.target as HTMLElement).id)
					}

					// Get all registered fields
					const fields = this.#getFields(form.id)
					if (fields instanceof ValidationErr) {
						console.error(...fields.error())
						return
					}

					// If all the fields are valid, then re-trigger the submit event
					// Some fields may have async rules, a callback is required for that.
					// The callback triggers the event if it's called for every field
					let validCounter = 0
					const fieldCount = fields.size
					const validFieldCB = () => {
						validCounter++
						if (validCounter == fieldCount) {
							if (this.#debug) {
								console.info("valid", (event.target as HTMLElement).id)
							}
							// TODO Constructor option to toggle what happens here
							// form.submit()
							// form.dispatchEvent(new Event("valid"))
							// https://developer.mozilla.org/en-US/docs/Web/Events/Creating_and_triggering_events
							form.dispatchEvent(new CustomEvent("valid"))
						}
					}

					// Check all fields
					for (const [fieldID] of fields) {
						const target = this.getFieldElement(form.id, fieldID)
						if (target instanceof ValidationErr) {
							console.error(...target.error())
							return
						}
						this.checkElement(target, validFieldCB)
					}
				})

				// Change event
				// https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/change_event
				// https://stackoverflow.com/a/51846602/639133
				Events.addListenerToElement(Validation.name, 
					form, change, (event: Event) => {

					if (event.target != null) {
						if (this.#debug) {
							console.info("change", (event.target as HTMLElement).id)
						}
						this.checkElement(event.target as HTMLElement, null)
					}
				})

				// Keyup event
				// https://developer.mozilla.org/en-US/docs/Web/API/Element/keyup_event
				Events.addListenerToElement(Validation.name, 
					form, keyup, (event: Event) => {

					const kbe = (event as KeyboardEvent)
					// See keyup_event link above: "Firefox bug 354358..."
					if (kbe.isComposing || (kbe.keyCode && kbe.keyCode === 229)) {
						return
					}
					// "key property [take] into consideration the state of modifier keys 
					// such as Shift as well as the keyboard locale and layout"
					// https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/key
					const key = kbe.key
					switch (key) {
						// Keys to ignore
						case "Alt":
						case "Control":
						case "Enter":
						case "Escape":
						case "Meta":
						case "Shift":
						case "Tab":
							return
					}
					if (kbe.target != null) {
						const el = (kbe.target as HTMLElement)
						if (el.tagName.toLowerCase() == "button") {
							// Prevent missing id error when key is pressed on a button
							// e.g. pressing space-bar on the submit button
							return
						}

						// Check element after a delay
						const formID = this.#getFormID(el)
						if (formID instanceof ValidationErr) {
							console.error(...formID.error())
							return
						}
						const fieldID = (el).id
						if (fieldID.trim() == "") {
							console.error(sprintf(
								"missing id in form %s for field", formID), kbe.target)
							return
						}
						const err = this.debounce(formID, fieldID, () => {
							if (this.#debug) { console.info("keyup", fieldID) }
							this.checkElement(el, null)
						}, this.#debounceWaitMS)
						if (err instanceof ValidationErr) {
							console.error(...err.error())
							return
						}
					}
				})

				// Validate event can be triggered programmatically
				Events.addListenerToElement(Validation.name, 
					form, validationValidate, (event: Event) => {

					if (event.target != null) {
						if (this.#debug) {
							console.info("validate", (event.target as HTMLElement).id)
						}
						this.checkElement(event.target as HTMLElement, null)
					}
				})
			})

			// .......................................................................
			const invalidMsg = "invalid"
			const validMsg = "valid"
			const validationMsg = (mode: string, el: HTMLElement, fieldID: ElementID) => {
				// Get or register field 
				const formID = this.#getFormID(el as HTMLElement)
				if (formID instanceof ValidationErr) {
					console.error(...formID.error())
					return
				}
				const element = this.getFieldElement(formID, fieldID)
				if (element instanceof ValidationErr) {
					console.error(...element.error())
					return
				}
				const field = this.#registerField(formID, fieldID)
				if (field instanceof ValidationErr) {
					console.error(...field.error())
					return
				}

				// Set validation message
				switch (mode) {
					case invalidMsg:
						field.invalidMsg = el
						break
					case validMsg:
						field.validMsg = el
						break
					default:
						console.error(sprintf("unknown mode %s ", mode))
						return
				}
			}

			// x-invalid
			// Message directive for invalid form fields
			// TODO Make object destructuring syntax work with deno lint?
			// https://alpinejs.dev/advanced/extending#custom-directives
			// https://stackoverflow.com/a/37661289/639133
			// deno-lint-ignore no-explicit-any
			Alpine.directive(invalidMsg, (el: Node, o: any) => {
				validationMsg(invalidMsg, el as HTMLElement, o.expression)
			})

			// x-valid
			// Message directive for valid form fields
			// deno-lint-ignore no-explicit-any
			Alpine.directive(validMsg, (el: Node, o: any) => {
				validationMsg(validMsg, el as HTMLElement, o.expression)
			})

			// .......................................................................
			// x-validate-rules
			// Directive for custom rules to check, space delimited
			// deno-lint-ignore no-explicit-any
			Alpine.directive("validate-rules", (el: Node, o: any) => {
				const ruleNames = o.expression.split(" ")
				const target = el as HTMLElement
				if (!["input", "select", "textarea"].includes(
					target.tagName.toLowerCase())) {
					console.error(sprintf("invalid tag %s", target.tagName))
					return
				}
				if (target.id.trim() == "") {
					console.error("element with empty id", target)
					return
				}

				// Get or register field 
				const formID = this.#getFormID(target)
				if (formID instanceof ValidationErr) {
					console.error(...formID.error())
					return
				}
				const field = this.#registerField(formID, target.id)
				if (field instanceof ValidationErr) {
					console.error(...field.error())
					return
				}

				// Set rule
				field.customRuleMap = new Map<RuleName, boolean>
				for (const ruleName of ruleNames) {
					// boolean is set to true the first time the rule is checked
					field.customRuleMap.set(ruleName, false)
				}
			})

			// .......................................................................
			// x-validation-trigger
			// Can be used to trigger validation on an element,
			// when other elements are changed.
			// deno-lint-ignore no-explicit-any
			Alpine.directive("validation-trigger", (el: Node, o: any) => {
				const triggers = document.querySelectorAll(o.expression)
				const target = el as HTMLElement
				const formID = this.#getFormID(target)
				if (formID instanceof ValidationErr) {
					console.error(...formID.error())
					return
				}
				// Setup fields that trigger this target
				triggers.forEach((node) => {
					const el = node as HTMLElement
					const field = this.#getField(formID, el.id)
					if (field instanceof ValidationErr) {
						console.error(...field.error())
						return
					}
					if (field.triggers) {
						field.triggers.push(target)
					} else {
						field.triggers = [target]
					}
				})
			})
		})

		this.#initialised = true
	}
}
