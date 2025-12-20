/*
 * State management for application
 */

// Initial state 
const initialState = {
	user: null,
	isLoggedIn: false,
	currentPage: 'landing',
	upload: {
		uploadUUID: null,
		sessionID: null,
		files: []
	},
	
	// { targetUUID, targetName } 
	targets: [],

	// { ruleUUID, name, when: {}, then: {} }
	rules: [],

	// { ruleUUID, targetUUID }
	ruleBindings: [],

	// UI state
	activeTarget: null,
	activeRule: null,

	resultZipURL: null,
	loginError: null,
};

// Functions that react to state change
const subscribers = [];

// Subscribe listener to state changes
export function subscribe(listener) { 
	subscribers.push(listener);
}

// Storage for state management
export const store = new Proxy(initialState, {
	set(target, prop, value) { 
		// Get old value for debugging 
		const oldValue = target[prop];

		// Set new state value
		target[prop] = value;

		// Print state for debugging
		console.log(`State changed: ${prop}`);
		console.log('Old value: ', oldValue);
		console.log('New value: ', value);
		console.log("Updated state:\n", JSON.stringify(target, null, 2));
		// Notify all subscribers
		subscribers.forEach(fn => fn(prop, value)); 
		return true;
	}
});
