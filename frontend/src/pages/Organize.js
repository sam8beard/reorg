/*
 * Organize page
 */
import { store } from '../state.js';

export default function Organize(root, userData) { 
	root.innerHTML = `
		<div id='organize-root'>
			<div id='organize-options'>
				<button id='ai-organize-btn' type='click'>
					Organize with AI
				</button>
				<button id='rule-organize-btn' type='click'>
					Organize with Rules
				</button>
			</div>
			<div id='organize-file-preview'>
				<!-- Fetch all files with the correct upload id and display here -->	
			</div>
		</div>
	`;

	
}
