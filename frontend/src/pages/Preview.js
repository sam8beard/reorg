import { store } from '../state.js';
import { showOrganize } from '../navigation.js';
import { Tree, Folder, File } from 'https://cdn.jsdelivr.net/npm/@webreflection/file-tree/prod.js';

export default function Preview(root, user) { 
	const previewJSON = JSON.stringify(store.preview, null, " ");
	console.log(store.preview);
	root.innerHTML = `
		<div id='preview-root'>	
			<h1>Preview</h1>
			<div>
				<h3>How your files will be organized</h3>
			</div>
			<div id='tree-container'>
			</div>
		</div>	
		<div id='preview-options'>
			<button id='add-folder-btn'>Make another folder</button>
			<button id='download-zip-btn'>Download your organized files</button>
		</div>
	`;

	const addFolderBtn = root.querySelector('#add-folder-btn');
	addFolderBtn.addEventListener('click', () => showOrganize());

	const downloadFilesBtn = root.querySelector('#download-zip-btn');
	// downloadFilesBtn.addEventListener('click', () => down

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

	// Add unsorted folder
	const unsortedFolder = new Folder("Unsorted Files");
	tree.append(unsortedFolder);
	unsorted.forEach(([key, value]) => {
		if (key === "files") {
			const files = value;
			for (let file of files) {
				const fileName = file.fileName;		
				unsortedFolder.append(new File([], fileName));
			}
		}
	});

	previewContainer.appendChild(tree);
			

}
