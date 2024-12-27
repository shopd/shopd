// deno-lint-ignore-file

// TODO Create PR on https://github.com/bigskysoftware/htmx
// Additional types not defined in the official htmx.d.ts

export interface HtmxResponseEvent extends Event {
	detail: {
		xhr: XMLHttpRequest
		elt: HTMLElement
		target: HTMLElement
		requestConfig: any
		successful: boolean
		shouldSwap: boolean
	}
}
