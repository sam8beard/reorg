/*
 * Handlers and utils for organize page related actions
 */
import { showOrganize } from '../navigation.js';
import { store } from '../state.js';

export function onOrganizePageSubmit(e, root) {
	showOrganize(store.user);
} 

