import { store } from '../state.js';
import { showOrganize } from '../navigation.js';
export default function Preview(root, user) { 
	const previewJSON = JSON.stringify(store.preview, null, " ");
	root.innerHTML = `
		<div id='preview-root'>	
			<h1>Preview</h1>
			<div>
				<h3>How your files will be organized</h3>
			</div>
			<div>
				<p>${previewJSON}</p>
			</div>
		</div>	
		<div id='preview-options'>
			<button id='add-folder-btn'>Make another folder</button>
		</div>
	`;

	const addFolderBtn = root.querySelector('#add-folder-btn');
	addFolderBtn.addEventListener('click', () => showOrganize());
}
