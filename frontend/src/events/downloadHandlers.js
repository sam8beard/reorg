import { downloadZip, saveOrgData } from '../api';
import { Spinner } from 'spin.js';
import { store } from '../state.js';

const spinnerOpts = {
	lines: 10, // The number of lines to draw
	length: 6, // The length of each line
	width: 3, // The line thickness
	radius: 6, // The radius of the inner circle
	scale: 1, // Scales overall size of the spinner
	corners: 1, // Corner roundness (0..1)
	speed: 1, // Rounds per second
	rotate: 0, // The rotation offset
	animation: 'spinner-line-fade-quick', // The CSS animation name for the lines
	direction: 1, // 1: clockwise, -1: counterclockwise
	color: '#ffffff', // CSS color or array of colors
	fadeColor: 'transparent', // CSS color or array of colors
	shadow: '0 0 1px transparent', // Box-shadow for the lines
	zIndex: 2000000000, // The z-index (defaults to 2e9)
	className: 'spinner', // The CSS class to assign to the spinner
	position: 'absolute'
};

export async function onDownloadClick(e, container, fileStructure) { 
	const button = e.currentTarget;
	container.innerHTML = '';
	const spinner = new Spinner(spinnerOpts).spin(container);
	button.disabled = true;
	button.innerText = 'Downloading your organized files...';
	let success = false;
	try {
		const blob = await downloadZip(fileStructure);
		// If server does not return blob representing zip archive
		//if (!(blob instanceof Blob)) {
		//	throw new Error('Invalid download response');
		//}
		console.log(blob);
		
		const href = window.URL.createObjectURL(blob);

		const anch = document.createElement('a');

		anch.href = href;
		anch.download = "organized_files.zip";


		document.body.appendChild(anch);
		anch.click();

		document.body.removeChild(anch);
		window.URL.revokeObjectURL(href);
		success = true;

	} catch (err) {
		spinner.stop();
		const cntr = document.querySelector('#download-status-container');
		cntr.innerHTML = `
			<p> Failed to download organized files. Please try again. </p>
		`;
		const btn = document.querySelector('#download-zip-btn');
		btn.innerText = 'Download your organized files';
		btn.disabled = false;

		// Indicate unsuccessful download
		success = false;
		console.log(err);
	}	
	
	// If zip was successfully downloaded
	if (success) {
		spinner.stop();
		const cntr = document.querySelector('#download-status-container');
		cntr.innerHTML = `
			<p> Organized files downloaded as organized_files.zip </p>
		`;
		const btn = document.querySelector('#download-zip-btn');
		btn.innerText = 'Download your organized files';
		btn.disabled = false;

	}
	
	// If registered user, send ruleset information to backend for persistent storage
	if (store.identity.type === 'user') {
		console.log("Registered user detected, saving org data");
		const orgData = buildOrgData()
		try {
			const response = await saveOrgData(orgData);
		} catch(err) {
			console.log(err);
		}
	}
}

/* Builds registered user organizational data for persistent storage */
function buildOrgData() {

	const ruleSet = {
		uploadID: store.upload.uploadID,
		files: store.upload.files,
		targets: store.targets,
	};

	const orgData = {
		rules: store.rules,
		ruleBindings: store.ruleBindings,
		ruleSet: ruleSet,
	};

	return orgData;

}
