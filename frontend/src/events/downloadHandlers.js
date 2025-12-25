import { downloadZip } from '../api';

export async function onDownloadClick(fileStructure) { 
	try {
		console.log("firing 2");
		const blob = await downloadZip(fileStructure);
		console.log(blob);
		
		const href = window.URL.createObjectURL(blob);

		const anch = document.createElement('a');

		anch.href = href;
		anch.download = "organized_files.zip";


		document.body.appendChild(anch);
		anch.click();

		document.body.removeChild(anch);
		window.URL.revokeObjectURL(href);

	} catch (err) {
		console.log(err);
	}
}
