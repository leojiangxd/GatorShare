const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

export const getCsrfToken = () => {
  const cookies = document.cookie.split("; ");
  const csrfCookie = cookies.find((row) => row.startsWith("csrf_token="));
  return csrfCookie ? decodeURIComponent(csrfCookie.split("=")[1]) : null;
};

export async function getUsername() {
	const csrfToken = getCsrfToken();
	const response = await axios.get(`${apiBaseUrl}/api/v1/current-user`, {
		headers: {
			"X-CSRF-Token": csrfToken || "",
		},
		withCredentials: true,
	});
	return response.data.username;
}

export const formatTime = (dateString) => {
  const date = new Date(dateString);
  const now = new Date();
  const diffInMilliseconds = now - date;
  
  // Less than a minute
  if (diffInMilliseconds < 60000) {
    const seconds = Math.floor(diffInMilliseconds / 1000);
    return `${seconds} second${seconds !== 1 ? 's' : ''} ago`;
  }
  
  // Less than an hour
  if (diffInMilliseconds < 3600000) {
    const minutes = Math.floor(diffInMilliseconds / 60000);
    return `${minutes} minute${minutes !== 1 ? 's' : ''} ago`;
  }

  // Less than a day
  if (diffInMilliseconds < 86400000) {
    const hours = Math.floor(diffInMilliseconds / 3600000);
    return `${hours} hour${hours !== 1 ? 's' : ''} ago`;
  }
  
  // Less than a week (7 days)
  if (diffInMilliseconds < 604800000) {
    const days = Math.floor(diffInMilliseconds / 86400000);
    return `${days} day${days !== 1 ? 's' : ''} ago`;
  }

  // Older than a week, format as MM/DD/YY
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const year = date.getFullYear().toString().slice(-2);
  return `${month}/${day}/${year}`;
};
