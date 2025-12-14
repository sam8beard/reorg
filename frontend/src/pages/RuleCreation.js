import { store } from '../state.js';
export default function RuleCreation(root, user) { 
	root.innerHTML = `
		<div id='rule-creation-root'>
			<!-- How to get most recently added target dir??? -->
			<h1> Add rules for ${store.targets[0].targetName} </h1>
		</div>
	`;
}
