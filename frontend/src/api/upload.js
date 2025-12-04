import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

/* Post file data to backend */
export async function uploadFileForm(formData) { 
	var uploadURL = DEV_API_BASE + "/upload";
	console.log(uploadURL);
	try { 
		const response = await axios.post(uploadURL, {
			formData: formData,
		});
		return response.data;

	} catch (err) { 
		if (err.response) { 
			console.log(err);
			return { error: err.response.data.error || 'Unknown error' };
		} else { 
			console.log(err);
			return { error: 'Network error' };
		} 
	} 

}

