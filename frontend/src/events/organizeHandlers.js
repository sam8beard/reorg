/*
 * Handlers and utils for organize page related actions
 */
import { showOrganize } from '../navigation.js';
import { fetchFiles } from '../api';
import { store } from '../state.js';

export async function onOrganizePageSubmit(e, root) {
	// Fetch files using upload ID for file preview on organize page
	const uploadId = store.upload.uploadID
	const response = await fetchFiles(uploadId);
	store.upload.files = response
	showOrganize();
}
