import { store } from '../state.js';
export default function RuleCreation(root, user) {
	// Get most recently pushed target onto target list
	const currTarget = store.targets[store.targets.length - 1]
	root.innerHTML = `
		<div id='rule-creation-root'>
			<!-- Is this the optimal way to get the most recently added target dir??? -->
			<h1> Add rules for ${currTarget.targetName} </h1>
		</div>
	`;
}
