/*
 * State management for application
 */

// Initial state 
const initialState = {
	user: null,
	isLoggedIn: false,
	currentPage: 'landing',
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
		target[prop] = value;
		// Notify all subscribers
		subscribers.forEach(fn => fn(prop, value)); 
		return true;
	}
});
