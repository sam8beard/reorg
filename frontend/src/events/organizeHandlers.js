/*
 * Handlers and utils for organize related actions
 */
import { showOrganize } from '../navigation.js';
import { store } from '../state.js';
/*
 * Adds event for guest button click on landing page.
 *
 * Shows organize page
 */
export function attachOrganizePageHandler(button, root) {
	button.addEventListener('click', () => showOrganize());
}

/*
 * Handle uploaded files via drop
 */
export async function dropHandler(e) { 
}
