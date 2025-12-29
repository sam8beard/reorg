import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

/* Post file data to backend */
export async function uploadFileForm(formData, progressBar) { 
	const uploadURL = DEV_API_BASE + "/upload";
	try { 
		console.log(formData);
		const token = sessionStorage.getItem('authToken');
		const response = await axios.post(uploadURL, formData, {
			headers: {
				'Authorization': `Bearer ${token}`
			},
			onUploadProgress: (progEvent) => {
				const percent = Math.round((progEvent.loaded * 100) / progEvent.total);
				progressBar.style.width = percent + '%';
			}
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

