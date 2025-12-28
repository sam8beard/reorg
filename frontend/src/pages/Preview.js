import { store } from '../state.js';
import { showOrganize } from '../navigation.js';
import { downloadZip } from '../api';
import { onDownloadClick } from '../events';
import { Spinner } from 'spin.js';
import { Tree, Folder, File } from 'https://cdn.jsdelivr.net/npm/@webreflection/file-tree/prod.js';

export default function Preview(root, user) { 
	root.innerHTML = `
		<div id='preview-root'>	
			<div id='preview-options'>
				<button id='add-folder-btn'>Make another folder</button>
				<button id='download-zip-btn'>Download your organized files</button>
			</div>
			<div id='download-status-container' style='height: 2em; margin: 1em auto; display: flex; align-items: center; justify-content: center;'></div>

			<div style='margin: 1em;'>
				<h2>How your files will be organized</h2>
			</div>
			<div id='tree-container'>
			</div>
		</div>
	`;

	const addFolderBtn = root.querySelector('#add-folder-btn');
	addFolderBtn.addEventListener('click', () => showOrganize());


	const downloadFilesBtn = root.querySelector('#download-zip-btn');
	const downloadStatusContainer = root.querySelector('#download-status-container');
	downloadFilesBtn.addEventListener('click', (e) => onDownloadClick(e, downloadStatusContainer, store.preview));

	const previewContainer = root.querySelector('#tree-container');
	displayPreview(previewContainer, store.preview);
}

function displayPreview(previewContainer, preview) {
	const tree = new Tree;
	const folders = Object.entries(preview.folders);
	const unsorted = Object.entries(preview.unmatched);

	// Add each folder
	folders.forEach(([key, value]) => {
		const newFolder = new Folder(value.targetName);
		tree.append(newFolder);
		// Add all files to folder
		for (let file of value.files) {
			const fileName = file.fileName;
			newFolder.append(new File([], fileName));	
		}
	});
	
	// Add unsorted folder if unsorted files exist
	if (haveUnsorted(unsorted)) {
		const unsortedFolder = new Folder("Unsorted Files");
		tree.append(unsortedFolder);
		unsorted.forEach(([key, value]) => {
			if (key === "files") {
				const files = value;
				if (files) {
					for (let file of files) {
						const fileName = file.fileName;
						unsortedFolder.append(new File([], fileName));
					}
				}
			}
		});
	}

	// Add file tree to DOM 
	previewContainer.appendChild(tree);
}

function haveUnsorted(unsortedObj) {
	for (const [key, value] of unsortedObj) {
		if (key === "files") {
			return value.length > 0; 
		}
	}
}
