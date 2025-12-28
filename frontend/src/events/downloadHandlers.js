import { downloadZip } from '../api';
import { Spinner } from 'spin.js';

export async function onDownloadClick(e, container, fileStructure) { 
	const spinner = new Spinner().spin();
	e.currentTarget.disabled = true;
	try {
		const blob = await downloadZip(fileStructure);
		container.appendChild(spinner.el);

		console.log(blob);
		
		const href = window.URL.createObjectURL(blob);

		const anch = document.createElement('a');

		anch.href = href;
		anch.download = "organized_files.zip";


		document.body.appendChild(anch);
		anch.click();

		document.body.removeChild(anch);
		window.URL.revokeObjectURL(href);
		spinner.stop();

	} catch (err) {
		console.log(err);
	}
}
