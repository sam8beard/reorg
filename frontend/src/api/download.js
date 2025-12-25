import axios from 'axios'; 
const DEV_API_BASE = "http://localhost:5173/api";

/* Downloads zip file of organized file structure */
export async function downloadZip(organizeResult) {
	const downloadURL = DEV_API_BASE + "/download/zip";
	try {
		const response = await axios.post(downloadURL, organizeResult);
		return response.data;
	} catch(err) { 
		if (err.response) {
			console.log(err); 
			return { error: err.response.data.error || 'Unknown error' };
		} else { 
			console.log(err);
			return { error: 'Network error' };
		}
	}
}
