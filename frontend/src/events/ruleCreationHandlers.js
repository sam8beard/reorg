import { store } from '../state.js';

export async function onRuleSubmit(event, root) { 
	event.preventDefault();

	const form = event.currentTarget;

	// This allows us to get the form data as key:value pairs
	const formData = new FormData(form);
	console.log(Object.fromEntries(formData));
	for (var pair of formData.entries()) {
		console.log(pair[0] + ": " + pair[1]);
	}
}
